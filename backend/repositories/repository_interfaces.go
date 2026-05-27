package repositories

import (
	"database/sql"
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

type UserRepository interface {
	Save(email, hashedPassword string) (int64, error)
	FindByEmail(email string) (int64, string, string, error)
}

type userRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

type RagRepository interface {
	AddData(userID int64, data *models.RagData) error
	SearchSimilar(userID int64, embedding []float32, limit int) ([]models.RagData, error)
	GetAllData(userID int64) ([]models.RagData, error)
}

type ragRepositoryImpl struct {
	db *sql.DB
}

func NewRagRepository(db *sql.DB) RagRepository {
	return &ragRepositoryImpl{db: db}
}
