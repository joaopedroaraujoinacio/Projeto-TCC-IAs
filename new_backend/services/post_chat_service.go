package services

import (
	"fmt"
	"strings"
	"go-project/models"
	"go-project/repositories"
)


type ChatService interface {
	StreamChat(request *models.ChatRequest) (<-chan string, <-chan error)
}

type chatService struct {
	chatRepo repositories.ChatRepository
}

func NewChatService(chatRepo repositories.ChatRepository) ChatService {
	return &chatService{
		chatRepo: chatRepo,
	}
}

func (s *chatService) StreamChat(request *models.ChatRequest) (<-chan string, <-chan error) {
	messageChan := make(chan string, 10)
	errorChan := make(chan error, 1)

	// Validate message
	if strings.TrimSpace(request.Message) == "" {
		errorChan <- fmt.Errorf("message cannot be empty")
		close(messageChan)
		close(errorChan)
		return messageChan, errorChan
	}

	// Start streaming
	go func() {
		defer close(messageChan)
		defer close(errorChan)

		// Call repository's SendToLLMStreaming
		streamChan, err := s.chatRepo.SendToLLM(request)
		if err != nil {
			errorChan <- fmt.Errorf("failed to start streaming: %w", err)
			return
		}

		// Read from repository's channel
		for chunk := range streamChan {
			// Check for errors in chunk
			if chunk.Error != nil {
				errorChan <- chunk.Error
				return
			}

			// Send text to handler
			if chunk.Text != "" {
				messageChan <- chunk.Text
			}

			// Stop if done
			if chunk.Done {
				return
			}
		}
	}()

	return messageChan, errorChan
}
