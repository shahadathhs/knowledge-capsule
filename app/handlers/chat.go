package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"knowledge-capsule/app/middleware"
	"knowledge-capsule/app/models"
	"knowledge-capsule/pkg/utils"

	"github.com/gorilla/websocket"
)

var (
	corsOrigins []string
	upgrader    = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			if len(corsOrigins) == 0 {
				return false
			}
			return middleware.IsOriginAllowed(r.Header.Get("Origin"), corsOrigins)
		},
	}
	clients   = make(map[string]*websocket.Conn) // UserID -> Conn
	clientsMu sync.Mutex
)

// InitChat initializes CORS for WebSocket. MessageStore is set via InitStores.
func InitChat(allowedOrigins []string) {
	corsOrigins = allowedOrigins
}

// WSMessage is the WebSocket message envelope.
type WSMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// WSSendPayload is the payload for type "send".
type WSSendPayload struct {
	ReceiverID string             `json:"receiver_id"`
	Content    string             `json:"content"`
	Type       models.MessageType `json:"type"`
	FileURL    string             `json:"file_url,omitempty"`
}

// WSGetHistoryPayload is the payload for type "get_history".
type WSGetHistoryPayload struct {
	UserID string `json:"user_id"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

func sendWSResponse(conn *websocket.Conn, msgType string, payload interface{}) {
	conn.WriteJSON(map[string]interface{}{"type": msgType, "payload": payload})
}

func ChatWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserContextKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("WebSocket upgrade error", "error", err)
		return
	}
	defer conn.Close()

	clientsMu.Lock()
	clients[userID] = conn
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(clients, userID)
		clientsMu.Unlock()
	}()

	for {
		var raw map[string]json.RawMessage
		err := conn.ReadJSON(&raw)
		if err != nil {
			slog.Error("WebSocket read error", "error", err)
			break
		}

		msgType := ""
		if t, ok := raw["type"]; ok {
			json.Unmarshal(t, &msgType)
		}

		// Legacy format: { receiver_id, content, type } without envelope - treat as send
		if msgType == "" && raw["receiver_id"] != nil {
			msgType = "send"
			raw["payload"], _ = json.Marshal(raw)
		}

		payloadBytes := raw["payload"]
		if payloadBytes == nil {
			payloadBytes = []byte("{}")
		}

		switch msgType {
		case "send":
			var sendPayload WSSendPayload
			if err := json.Unmarshal(payloadBytes, &sendPayload); err != nil {
				sendWSResponse(conn, "error", map[string]string{"message": "invalid send payload"})
				continue
			}
			savedMsg, err := MessageStore.SaveMessage(userID, sendPayload.ReceiverID, sendPayload.Content, sendPayload.Type, sendPayload.FileURL)
			if err != nil {
				slog.Error("Chat save error", "error", err)
				sendWSResponse(conn, "error", map[string]string{"message": "failed to save message"})
				continue
			}
			sendWSResponse(conn, "message", savedMsg)
			clientsMu.Lock()
			receiverConn, ok := clients[sendPayload.ReceiverID]
			clientsMu.Unlock()
			if ok {
				receiverConn.WriteJSON(map[string]interface{}{"type": "message", "payload": savedMsg})
			}

		case "get_history":
			var histPayload WSGetHistoryPayload
			if err := json.Unmarshal(payloadBytes, &histPayload); err != nil {
				sendWSResponse(conn, "error", map[string]string{"message": "invalid get_history payload"})
				continue
			}
			if histPayload.UserID == "" {
				sendWSResponse(conn, "error", map[string]string{"message": "user_id required"})
				continue
			}
			page, limit := histPayload.Page, histPayload.Limit
			if page < 1 {
				page = 1
			}
			if limit < 1 || limit > 100 {
				limit = 20
			}
			messages, err := MessageStore.GetMessagesBetweenUsers(userID, histPayload.UserID)
			if err != nil {
				sendWSResponse(conn, "error", map[string]string{"message": "failed to fetch history"})
				continue
			}
			paged, total := utils.SlicePage(messages, page, limit)
			sendWSResponse(conn, "history", map[string]interface{}{
				"data":  paged,
				"page":  page,
				"limit": limit,
				"total": total,
			})

		default:
			sendWSResponse(conn, "error", map[string]string{"message": "unknown message type: " + msgType})
		}
	}
}
