package handlers

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"knowledge-capsule-api/app/middleware"
	"knowledge-capsule-api/app/models"
	"knowledge-capsule-api/app/store"
	"knowledge-capsule-api/pkg/utils"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for now
		},
	}
	clients   = make(map[string]*websocket.Conn) // UserID -> Conn
	clientsMu sync.Mutex
	msgStore  = store.MessageStore{FileStore: store.FileStore[models.Message]{FilePath: "data/messages.json"}}
)

// InitChatStore initializes the chat storage.
func InitChatStore() error {
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
		log.Println("Upgrade error:", err)
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
			log.Println("Read error:", err)
			break
		}

		// Save message
		savedMsg, err := msgStore.SaveMessage(userID, msg.ReceiverID, msg.Content, msg.Type, msg.FileURL)
		if err != nil {
			log.Println("Save error:", err)
			continue
		}

		// Send to receiver if online
		clientsMu.Lock()
		receiverConn, ok := clients[msg.ReceiverID]
		clientsMu.Unlock()

		if ok {
			if err := receiverConn.WriteJSON(savedMsg); err != nil {
				log.Println("Write error:", err)
			}
		}
	}
}

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

	utils.JSONResponse(w, http.StatusOK, true, "Messages fetched successfully", messages)
}
