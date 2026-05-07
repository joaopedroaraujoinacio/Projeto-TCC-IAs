package providers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go-project/models"
	"net/http"
	"strings"
)


func NewOpenAIProvider(apiKey string) *OpenAIProvider {
	return &OpenAIProvider{
		Client: &http.Client{},
		ApiKey: apiKey,
		ApiURL: "https://api.openai.com/v1/chat/completions",
	}
}


type OpenAIProvider struct {
	Client *http.Client
	ApiKey string
	ApiURL string
}

type OpenAIRequest struct {
	Model    string              `json:"model"`
	Messages []map[string]string `json:"messages"`
	Stream   bool                `json:"stream"`
}

type OpenAIStreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
	Usage *struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func (p *OpenAIProvider) SendToOpenAI(request *models.ChatRequest) (<-chan models.StreamChunk, error) {
	model := request.Model
	if model == "" {
		// model = "gpt-4o-mini"
	}

	messages := []map[string]string{}
	for _, msg := range request.History {
		messages = append(messages, map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		})
	}
	messages = append(messages, map[string]string{
		"role":    "user",
		"content": request.Message,
	})

	fmt.Printf("[OpenAI] Sending request with model=%s, messages=%d\n", model, len(messages))

	payload := OpenAIRequest{
		Model:    model,
		Messages: messages,
		Stream:   true,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal openai request: %w", err)
	}

	req, err := http.NewRequest("POST", p.ApiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create openai request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.ApiKey)

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to openai: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("openai returned status: %d", resp.StatusCode)
	}

	streamChan := make(chan models.StreamChunk, 10)

	go func() {
		defer close(streamChan)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()

			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")

			if data == "[DONE]" {
				streamChan <- models.StreamChunk{Done: true}
				return
			}

			var chunk OpenAIStreamChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				fmt.Printf("[OpenAI] Parse error: %v\n", err)
				continue
			}

			if len(chunk.Choices) == 0 {
				continue
			}

			choice := chunk.Choices[0]
			isDone := choice.FinishReason != nil && *choice.FinishReason == "stop"

			streamChan <- models.StreamChunk{
				Text: choice.Delta.Content,
				Done: isDone,
			}

			if isDone {
				return
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("[OpenAI] Scanner error: %v\n", err)
			streamChan <- models.StreamChunk{
				Error: fmt.Errorf("scanner error: %w", err),
			}
		}
	}()

	return streamChan, nil
}

