package services

import (
	"fmt"
	"strings"
	"go-project/models"
	"go-project/repositories"
)

type ChatStreamService interface {
	StreamChat(request *models.ChatRequest) (<-chan string, <-chan error)
}

type chatStreamService struct {
	chatRepo repositories.ChatRepository
}

func NewChatStreamService(chatRepo repositories.ChatRepository) ChatStreamService {
	return &chatStreamService{
		chatRepo: chatRepo,
	}
}

func (s *chatStreamService) StreamChat(request *models.ChatRequest) (<-chan string, <-chan error) {
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
		streamChan, err := s.chatRepo.SendToLLMStreaming(request)
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
