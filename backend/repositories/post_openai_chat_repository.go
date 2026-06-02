package repositories

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go-project/models"
	"net/http"
	"strings"
)

func (r *chatRepository) SendToOpenAI(request *models.ChatRequest) (<-chan models.StreamChunk, error) {
    if r.openAIKey == "" {
        return nil, fmt.Errorf("OpenAI API key not configured")
    }

    model := request.Model
    if model == "" {
        model = "gpt-4.1-mini"
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

    payload := models.OpenAIRequest{
        Model:         model,
        Messages:      messages,
        Stream:        true,
        StreamOptions: &models.StreamOptions{IncludeUsage: true},
    }

    jsonData, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal openai request: %w", err)
    }

    req, err := http.NewRequest("POST", r.openAIURL, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to create openai request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+r.openAIKey)

    resp, err := r.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to reach openai: %w", err)
    }

    if resp.StatusCode != http.StatusOK {
        resp.Body.Close()
        return nil, fmt.Errorf("openai returned status: %d", resp.StatusCode)
    }

    streamChan := make(chan models.StreamChunk, 10)

    go func() {
        defer close(streamChan)
        defer resp.Body.Close()

        var completionTokens, promptTokens int

        scanner := bufio.NewScanner(resp.Body)
        for scanner.Scan() {
            line := scanner.Text()

            if !strings.HasPrefix(line, "data: ") {
                continue
            }

            data := strings.TrimPrefix(line, "data: ")

            if data == "[DONE]" {
                streamChan <- models.StreamChunk{
                    Done:         true,
                    TokenCount:   completionTokens,
                    PromptTokens: promptTokens,
                }
                return
            }

            var chunk models.OpenAIStreamChunk
            if err := json.Unmarshal([]byte(data), &chunk); err != nil {
                fmt.Printf("[OpenAI] parse error: %v\n", err)
                continue
            }

            if chunk.Usage != nil {
                completionTokens = chunk.Usage.CompletionTokens
                promptTokens = chunk.Usage.PromptTokens
            }

            if len(chunk.Choices) == 0 {
                continue
            }

            if chunk.Choices[0].Delta.Content != "" {
                streamChan <- models.StreamChunk{
                    Text: chunk.Choices[0].Delta.Content,
                }
            }
        }

        if err := scanner.Err(); err != nil {
            streamChan <- models.StreamChunk{Error: fmt.Errorf("scanner error: %w", err)}
        }
    }()

    return streamChan, nil
}
