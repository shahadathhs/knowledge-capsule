package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"knowledge-capsule/app/middleware"
	"knowledge-capsule/app/models"
	"knowledge-capsule/app/store"
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
	msgStore  = store.MessageStore{FileStore: store.FileStore[models.Message]{FilePath: "data/messages.json"}}
)

// InitChat initializes the chat storage and CORS for WebSocket.
func InitChat(allowedOrigins []string) error {
	corsOrigins = allowedOrigins
	return msgStore.FileStore.Init()
}

type ChatMessage struct {
	ReceiverID string             `json:"receiver_id"`
	Content    string             `json:"content"`
	Type       models.MessageType `json:"type"`
	FileURL    string             `json:"file_url,omitempty"`
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
		var msg ChatMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			slog.Error("WebSocket read error", "error", err)
			break
		}

		// Save message
		savedMsg, err := msgStore.SaveMessage(userID, msg.ReceiverID, msg.Content, msg.Type, msg.FileURL)
		if err != nil {
			slog.Error("Chat save error", "error", err)
			continue
		}

		// Send to receiver if online
		clientsMu.Lock()
		receiverConn, ok := clients[msg.ReceiverID]
		clientsMu.Unlock()

		if ok {
			if err := receiverConn.WriteJSON(savedMsg); err != nil {
				slog.Error("WebSocket write error", "error", err)
			}
		}
	}
}

// GetChatHistoryHandler godoc
// @Summary Get chat history
// @Description Get paginated chat history between current user and another user
// @Tags chat
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param user_id query string true "Other user ID"
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 20, max 100)"
// @Success 200 {object} models.PaginatedResponse "Paginated messages: data, page, limit, total"
// @Failure 400 {object} map[string]interface{}
// @Router /api/chat/history [get]
func GetChatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserContextKey).(string)
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	otherUserID := r.URL.Query().Get("user_id")
	if otherUserID == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("missing user_id parameter"))
		return
	}

	messages, err := msgStore.GetMessagesBetweenUsers(userID, otherUserID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("failed to fetch messages"))
		return
	}

	page, limit := utils.ParsePagination(r)
	paged, total := utils.SlicePage(messages, page, limit)
	utils.JSONPaginatedResponse(w, http.StatusOK, "Messages fetched successfully", paged, page, limit, total)
}
