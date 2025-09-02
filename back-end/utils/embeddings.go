package utils

import (
	"log"
	"fmt"
	"bytes"
	"net/http"
	"encoding/json"
)


type EmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type EmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

func GenerateEmbedding(text string) ([]float32, error) {
	reqBody := EmbeddingRequest{
		Model:  "nomic-embed-text",
		Prompt: text,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	log.Printf("Sending to Ollama: %s", string(jsonData))

	// resp, err := http.Post("http://localhost:11434/api/embeddings", "application/json", bytes.NewBuffer(jsonData))
	// resp, err := http.Post("http://host.docker.internal:11434/api/embeddings", "application/josn", bytes.NewBuffer(jsonData))
	resp, err := http.Post("http://ollama:11434/api/embeddings", "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rawResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&rawResp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	log.Printf("Raw Ollama response: %+v", rawResp)
	
	if embedding, ok := rawResp["embedding"].([]interface{}); ok {
		result := make([]float32, len(embedding))
		for i, v := range embedding {
			if f, ok := v.(float64); ok {
				result[i] = float32(f)
			}
		}
		return result, nil
	}
	return nil, fmt.Errorf("no embedding found in response")
}

func GenerateCodeEmbedding(text string) ([]float32, error) {
	reqBody := EmbeddingRequest{
		Model:  "mxbai-embed-large",
		Prompt: text,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	log.Printf("Sending to Ollama: %s", string(jsonData))

	resp, err := http.Post("http://ollama:11434/api/embeddings", "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rawResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&rawResp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	log.Printf("Raw Ollama response: %+v", rawResp)
	
	if embedding, ok := rawResp["embedding"].([]interface{}); ok {
		result := make([]float32, len(embedding))
		for i, v := range embedding {
			if f, ok := v.(float64); ok {
				result[i] = float32(f)
			}
		}
		return result, nil
	}
	return nil, fmt.Errorf("no embedding found in response")
}

