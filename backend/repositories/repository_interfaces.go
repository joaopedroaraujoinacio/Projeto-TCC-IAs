package repositories

import (
	"net/http"
	"time"

	"go-project/models"
)

type ChatRepository interface {
	SendToLLM(request *models.ChatRequest) (<-chan models.StreamChunk, error)
	SendToOpenAI(request *models.ChatRequest) (<-chan models.StreamChunk, error)
	SendToGemini(request *models.ChatRequest) (<-chan models.StreamChunk, error)
}

func NewChatRepository(ollamaURL, ollamaDefaultModel, openAIKey, geminiKey string) ChatRepository {
	return &chatRepository{
		ollamaURL: ollamaURL,
		ollamaDefaultModel: ollamaDefaultModel,
		openAIURL: "https://api.openai.com/v1/chat/completions",
		openAIKey: openAIKey,
		geminiURL: "https://generativelanguage.googleapis.com/v1beta",
		geminiKey: geminiKey,
		client: &http.Client{
			Timeout: 500 * time.Second,
		},
	}
}

type chatRepository struct {
	ollamaURL string
	openAIURL string
	openAIKey string
	geminiURL string
	geminiKey string
	ollamaDefaultModel string
	client    *http.Client
}
