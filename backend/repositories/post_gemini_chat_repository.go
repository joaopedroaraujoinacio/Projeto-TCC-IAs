package repositories

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go-project/models"
)

func (r *chatRepository) SendToGemini(request *models.ChatRequest) (<-chan models.StreamChunk, error) {
	key := r.geminiKey
	if key == "" {
		return nil, fmt.Errorf("Gemini API key not configured")
	}

	model := request.Model
	if model == "" {
		model = "gemini-2.0-flash"
	}

	// Gemini uses a different message format — no "system" role, alternating user/model
	contents := []models.GeminiContent{}
	for _, msg := range request.History {
		role := msg.Role
		if role == "assistant" {
			role = "model" // Gemini calls it "model" not "assistant"
		}
		contents = append(contents, models.GeminiContent{
			Role:  role,
			Parts: []models.GeminiPart{{Text: msg.Content}},
		})
	}
	contents = append(contents, models.GeminiContent{
		Role:  "user",
		Parts: []models.GeminiPart{{Text: request.Message}},
	})

	fmt.Printf("[Gemini] model=%s messages=%d\n", model, len(contents))

	payload := models.GeminiRequest{Contents: contents}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal gemini request: %w", err)
	}

	url := fmt.Sprintf("%s/models/%s:streamGenerateContent?alt=sse&key=%s", r.geminiURL, model, key)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to reach gemini: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("gemini returned status %d: %s", resp.StatusCode, string(body))
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
			if data == "" {
				continue
			}

			var chunk models.GeminiStreamChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				fmt.Printf("[Gemini] parse error: %v\n", err)
				continue
			}

			if len(chunk.Candidates) == 0 {
				continue
			}

			candidate := chunk.Candidates[0]
			isDone := candidate.FinishReason == "STOP"

			text := ""
			if len(candidate.Content.Parts) > 0 {
				text = candidate.Content.Parts[0].Text
			}

			streamChan <- models.StreamChunk{
				Text: text,
				Done: isDone,
			}

			if isDone {
				return
			}
		}

		if err := scanner.Err(); err != nil {
			streamChan <- models.StreamChunk{Error: fmt.Errorf("scanner error: %w", err)}
		}
	}()

	return streamChan, nil
}
