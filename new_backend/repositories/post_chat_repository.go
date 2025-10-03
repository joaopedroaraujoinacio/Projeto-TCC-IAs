package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-project/models"
	"io"
	"net/http"
	"strings"
	"time"
)


type ChatRepository interface {
	SendToLLM(request *models.ChatRequest) (*models.ChatResponse, error)
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
	client	*http.Client
}

func (r *chatRepository) SendToLLM(request *models.ChatRequest) (*models.ChatResponse, error) {
	model := request.Model
	if model == "" {
		model = "llama3.2:3b"
	}

	ollamaReq := models.OllamaRequest{
		Model: model,
		Prompt: request.Message,
	}

	jsonData, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, fmt.Errorf("failed to format request: %w", err)
	}

	resp, err := r.client.Post(r.ollamaURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to send the formated request to ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get response from ollama: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// fmt.Printf(string(body))



	// var ollamaResp models.OllamaResponse
	// if err := json.Unmarshal(body, &ollamaResp); err != nil {
	// 	return nil, fmt.Errorf("failed to format ollama response: %w", err)
	// }

// 	return &models.ChatResponse{
// 		Response: ollamaResp.Response,
// 		Model: model,
// 	}, nil
// }

	lines := strings.Split(string(body), "\n")
	var fullResponse strings.Builder
	// var model string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		var ollamaResp models.OllamaResponse
		if err := json.Unmarshal([]byte(line), &ollamaResp); err != nil {
			continue // Skip invalid JSON lines
		}
		
		fullResponse.WriteString(ollamaResp.Response)
		// model = ollamaResp.Model
	}

	return &models.ChatResponse{
		Response: fullResponse.String(),
		Model: model,
	}, nil
}
