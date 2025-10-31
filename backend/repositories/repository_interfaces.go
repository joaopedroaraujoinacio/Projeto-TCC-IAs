package repositories

import (
	"time"
	"net/http"
	"go-project/models"
)


type ChatRepository interface {
    SendToLLM(request *models.ChatRequest) (<-chan models.StreamChunk, error)
}

func NewChatRepository(ollamaURL string) ChatRepository {
    return &chatRepository{
        ollamaURL: ollamaURL,
        client: &http.Client{
            Timeout: 500 * time.Second,
        },
    }
}

type chatRepository struct {
    ollamaURL string
    client    *http.Client
}
