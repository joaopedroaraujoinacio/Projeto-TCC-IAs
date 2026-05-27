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

type UserHandler struct {
	svc *services.UserService
}

func NewUserHandler(svc *services.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

type RagHandler struct {
	svc *services.RagService
}

func NewRagHandler(svc *services.RagService) *RagHandler {
	return &RagHandler{svc: svc}
}

