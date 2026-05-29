package services

import (
	"fmt"
	"log"
	"strings"
	"go-project/models"
	"go-project/utils"
)

func (s *chatService) RagChat(userID int64, request *models.ChatRequest) (<-chan string, <-chan error) {
	messageChan := make(chan string, 10)
	errorChan   := make(chan error, 1)

	if strings.TrimSpace(request.Message) == "" {
		errorChan <- fmt.Errorf("message cannot be empty")
		close(messageChan)
		close(errorChan)
		return messageChan, errorChan
	}

	queryEmbedding, err := utils.GenerateEmbedding(request.Message)
	if err != nil {
		log.Printf("failed to generate query embedding: %v", err)
	}
	log.Printf("generated embedding for query: '%s' (length: %d)", request.Message, len(queryEmbedding))

	results, err := s.ragRepo.SearchSimilar(userID, queryEmbedding, 5)
	if err != nil {
		log.Printf("failed to search rag data: %v", err)
	}
	log.Printf("repository returned %d results for user %d", len(results), userID)

	ragContext := utils.FormatContext(results)
	request.Message = fmt.Sprintf(
		"Context: %s\n\nQuestion: %s. Please, if the context is not found or unrelated, make it clear to the user.",
		ragContext, request.Message,
	)

	go func() {
		defer close(messageChan)
		defer close(errorChan)

		streamChan, err := s.chatRepo.SendToLLM(request)
		if err != nil {
			errorChan <- fmt.Errorf("failed to start streaming: %w", err)
			return
		}

		for chunk := range streamChan {
			if chunk.Error != nil {
				errorChan <- chunk.Error
				return
			}
			if chunk.Text != "" {
				messageChan <- chunk.Text
			}
			if chunk.Done {
				return
			}
		}
	}()

	return messageChan, errorChan
}
