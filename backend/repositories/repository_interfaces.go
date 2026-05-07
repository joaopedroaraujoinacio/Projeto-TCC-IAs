package repositories

import (
	"time"
	"net/http"
	"go-project/models"
)


type ChatRepository interface {
    SendToLLM(request *models.ChatRequest) (<-chan models.StreamChunk, error)
    SendToOpenAI(request *models.ChatRequest) (<-chan models.StreamChunk, error)
}

func NewChatRepository(ollamaURL string, openAIKey string) ChatRepository {
    return &chatRepository{
        ollamaURL: ollamaURL,
        openAIURL: "https://api.openai.com/v1/chat/completions",
        openAIKey: openAIKey,
        client: &http.Client{
            Timeout: 500 * time.Second,
        },
    }
}

type chatRepository struct {
    ollamaURL string
		openAIKey string
		openAIURL string
    client    *http.Client
}

