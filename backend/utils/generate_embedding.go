package utils

import (
	"log"
	"fmt"
	"bytes"
	"net/http"
	"encoding/json"
	"go-project/models"
)


func GenerateEmbedding(content string) ([]float32, error) {
	requestBody := models.EmbeddingRequest{
		Model: "nomic-embed-text",
		Prompt: content,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to format request: %w", err)
	}

	log.Printf("sending request to ollama: %s", string(jsonData))

	ollamaResponse, err := http.Post("http://ollama:11434/api/embeddings", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("ollama failed to get the request: %w", err)
	}
	defer ollamaResponse.Body.Close()

	var rawResponse map[string]any
	err = json.NewDecoder(ollamaResponse.Body).Decode(&rawResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to format ollama response: %w", err)
	}

	log.Printf("raw ollama response: %s", rawResponse)

	embeddingRaw, ok := rawResponse["embedding"];
	if !ok {
		return nil, fmt.Errorf("no embedding found in the response: %w", err)
	}

	embeddingSlice, ok := embeddingRaw.([]any);
	if !ok {
		return nil, fmt.Errorf("embedding is not a slice")
	}

	embeddingResult := make ([]float32, len(embeddingSlice))
	for i, v := range embeddingSlice {
		f, ok := v.(float64)
		if !ok {
			return  nil, fmt.Errorf("faild to format embedded response: %w", err)
		}
		embeddingResult[i] = float32(f)
	}
	return embeddingResult, nil
}

