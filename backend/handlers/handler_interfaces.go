package handlers

import (
	"database/sql"
	"go-project/services"
)


type ChatHandler struct {
	chatService services.ChatService
	db             *sql.DB
}

func NewChatHandler(chatService services.ChatService, db *sql.DB) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
		db:             db,
	}
}

