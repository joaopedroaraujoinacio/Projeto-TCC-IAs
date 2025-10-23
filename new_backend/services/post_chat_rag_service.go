package services

import (
	"log"
	"fmt"
	"strings"
	"database/sql"
	"go-project/utils"
	"go-project/models"
	"go-project/repositories"
)

type RagChatService interface {
	RagChatService(db *sql.DB, query string, limit int, request *models.ChatRequest) (<-chan string, <-chan error)
}

type ragChatService struct {
    chatRepo repositories.ChatRepository
}

func NewRagChatService(chatRepo repositories.ChatRepository) RagChatService {
    return &ragChatService{
        chatRepo: chatRepo,
    }
}

func (s *ragChatService) RagChatService(db *sql.DB, query string, limit int, request *models.ChatRequest) (<-chan string, <-chan error) {
	messageChan := make(chan string, 10)
	errorChan := make(chan error, 1)

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

	log.Printf("Generated embedding for query: '%s' (length: %d)", query, len(queryEmbedding))

	if len(queryEmbedding) > 0 {
		log.Printf("first 5 embedding values: %v", queryEmbedding[:5])
	}

	results, err := repositories.RagDataSimilaritySearch(db, request.Message, queryEmbedding, limit, 0.5, 0.5)
	if err != nil {
		log.Printf("failed to fetch search query from duckdb: %v", err)
	}

	log.Printf("repository returned %d results", len(results) )

	context := utils.FormatContext(results)

	request.Message = fmt.Sprintf(
		"Context: %s\n\nQuestion: %s. Please, if the context is not found or unrelated, make it clear to the user.", 
		context, request.Message,
	)

	go func() {
		defer close(messageChan)
		defer close(errorChan)

		streamChan, err := s.chatRepo.SendToLLM(request)
		if err != nil {
			errorChan <- fmt.Errorf("failed to start streaming: %w", err)
			return
		}

		buffer := ""
		
		for chunk := range streamChan {
			if chunk.Error != nil {
				errorChan <- chunk.Error
				return
			}

			buffer += chunk.Text

			for {
				spaceIdx := strings.Index(buffer, " ")
				if spaceIdx == -1 {
					break
				}
				word := buffer[:spaceIdx+1] // envia palavra completa
				messageChan <- word
				buffer = buffer[spaceIdx+1:]
			}

			if chunk.Done {
				if buffer != "" {
					messageChan <- buffer
				}
				return
			}
		}
	}()

	return messageChan, errorChan
}

