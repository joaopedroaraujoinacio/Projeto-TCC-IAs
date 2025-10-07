package services

import (
	"fmt"
	"log"
	"strings"
	"database/sql"
	"go-project/utils"
	"go-project/models"
	"go-project/repositories"
)

type RagChatService interface {
    RagChatService(db *sql.DB, query string, limit int, request *models.ChatRequest) (*models.ChatResponse, error)
}

type ragChatService struct {
    chatRepo repositories.ChatRepository
}

func NewRagChatService(chatRepo repositories.ChatRepository) RagChatService {
    return &ragChatService{
        chatRepo: chatRepo,
    }
}

func (s *ragChatService) RagChatService(db *sql.DB, query string, limit int, request *models.ChatRequest) (*models.ChatResponse, error) {
	if strings.TrimSpace(request.Message) == "" {
		return nil, fmt.Errorf("message cannot be empty")
	}

	queryEmbedding, err := utils.GenerateEmbedding(request.Message)
	if err != nil {
		log.Printf("failed to generate query embedding: %v", err)
		return nil, fmt.Errorf("failed to generate query embedding :%w", err)
	}

	log.Printf("Generated embedding for query: '%s' (lenght: %d)", query, len(queryEmbedding))

	if len(queryEmbedding) > 0 {
		log.Printf("first 5 embedding values: %v", queryEmbedding[:5])
	}

	results, err := repositories.RagDataSimilaritySearch(db, queryEmbedding, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch search query from duckdb: %w", err)
	}

	log.Printf("repository returned %d results", len(results) )

	context := formatContext(results)

	request.Message = fmt.Sprintf("Context: %s\n\nQuestion: %s. Please, if the context is not found or unrelaed, make it clear to the user.", context, request.Message)

	response, err := s.chatRepo.SendToLLM(request)
	if err != nil {
		return  nil, fmt.Errorf("failed to get response from LLM: %w", err)
	}

	return response, nil
}

func formatContext(ragData []models.RagData) string {
    if len(ragData) == 0 {
        return "No relevant context found."
    }
    
    var contextBuilder strings.Builder
    contextBuilder.WriteString("Relevant information:\n")
    
    for i, data := range ragData {
        contextBuilder.WriteString(fmt.Sprintf("%d. %s\n", i+1, data.Content))
    }
    
    return contextBuilder.String()
}

