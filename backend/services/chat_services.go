package services

import (
	"fmt"
	"golang_crud/models"
	"golang_crud/repositories"
	"strings"
)


type ChatService interface {
	ProcessChat(request *models.ChatRequest) (*models.ChatResponse, error)
}

type chatService struct {
	chatRepo repositories.ChatRepository
}

func NewChatService(chatRepo repositories.ChatRepository) ChatService {
	return &chatService{
		chatRepo: chatRepo,
	}
}

func (s *chatService) ProcessChat(request *models.ChatRequest) (*models.ChatResponse, error) {
	if strings.TrimSpace(request.Message) == "" {
		return nil, fmt.Errorf("message cannot be empty")
	}

	response, err := s.chatRepo.SendToLLM(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get response from LLM: %w", err)
	}

	return response, nil
}

