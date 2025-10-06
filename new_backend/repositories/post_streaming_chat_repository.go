package repositories

import (
	"fmt"
	"bufio"
	"bytes"
	"net/http"
	"encoding/json"
	"go-project/models"
)

type StreamChunk struct {
	Text  string
	Model string
	Done  bool
	Error error
}

func (r *chatRepository) SendToLLMStreaming(request *models.ChatRequest) (<-chan StreamChunk, error) {
	model := request.Model
	if model == "" {
		model = "llama3.2:3b"
	}

	ollamaReq := models.OllamaRequest{
		Model:  model,
		Prompt: request.Message,
		Stream: true, 
	}

	jsonData, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, fmt.Errorf("failed to format request: %w", err)
	}

	resp, err := r.client.Post(r.ollamaURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to send request to ollama: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to get response from ollama: %d", resp.StatusCode)
	}

	streamChan := make(chan StreamChunk, 10)

	go func() {
		defer close(streamChan)      
		defer resp.Body.Close()     

		scanner := bufio.NewScanner(resp.Body)
		
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue 
			}

			var ollamaResp models.OllamaResponse
			if err := json.Unmarshal([]byte(line), &ollamaResp); err != nil {
				streamChan <- StreamChunk{
					Error: fmt.Errorf("failed to parse response: %w", err),
				}
				continue
			}

			streamChan <- StreamChunk{
				Text:  ollamaResp.Response, 
				Done:  ollamaResp.Done,
			}

			if ollamaResp.Done {
				break
			}
		}

		if err := scanner.Err(); err != nil {
			streamChan <- StreamChunk{
				Error: fmt.Errorf("scanner error: %w", err),
			}
		}
	}()

	return streamChan, nil
}

