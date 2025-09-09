package repositories

import (
	"io"
	"fmt"
	"time"
	"bytes"
	"net/http"
	"encoding/json"
	"golang_crud/models"
)


type ChatRepository interface {
	SendToLLM(request *models.ChatRequest) (*models.ChatResponse, error)
}

type chatRepository struct {
	ollamaURL string
	client		*http.Client
}

func NewChatRepository(ollamaURL string) ChatRepository {
	return &chatRepository{
		ollamaURL: ollamaURL,
		client: &http.Client{
			Timeout: 500 * time.Second,
		},
	}
}

func (r *chatRepository) SendToLLM(request *models.ChatRequest) (*models.ChatResponse, error) {
	model := request.Model
	if model == "" {
		model = "llama3.2:3b"
	}

	ollamaReq := models.OllamaRequest{
		Model:  model,
		Prompt: request.Message,
		Stream: false,
	}

	jsonData, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := r.client.Post(r.ollamaURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var ollamaResp models.OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &models.ChatResponse{
		Response: ollamaResp.Response,
		Model:    model,
	}, nil
}

