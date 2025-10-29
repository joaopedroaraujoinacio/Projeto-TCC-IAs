package services

import (
	"database/sql"
	"go-project/models"
	"go-project/repositories"
	"go-project/utils"
)


type ChatService interface {
	StreamChat(request *models.ChatRequest) (<-chan string, <-chan error)
	RagChat(db *sql.DB, query string, limit int, request *models.ChatRequest) (<-chan string, <-chan error)
	WebSearchChat(request *models.WebSearchRequest) (*models.WebSearchResponse, <-chan string, <-chan error)
}

type chatService struct {
	chatRepo repositories.ChatRepository
	searchRepo utils.WebSearchRepository
}

func NewChatService(chatRepo repositories.ChatRepository, searchRepo utils.WebSearchRepository) ChatService {
	return &chatService{
		chatRepo: chatRepo,
		searchRepo: searchRepo,
	}
}

