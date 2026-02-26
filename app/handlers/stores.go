package handlers

import (
	"knowledge-capsule/app/store"

	"gorm.io/gorm"
)

var (
	UserStore    store.UserStore
	CapsuleStore store.CapsuleStore
	TopicStore   store.TopicStore
	MessageStore store.MessageStore
)

// InitStores initializes all stores with the database connection.
func InitStores(db *gorm.DB) {
	UserStore = store.NewUserStore(db)
	CapsuleStore = store.NewCapsuleStore(db)
	TopicStore = store.NewTopicStore(db)
	MessageStore = store.NewMessageStore(db)
}
