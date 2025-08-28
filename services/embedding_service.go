package services

import (
	"log"
	"fmt"
	"bytes"
	"encoding/json"
	"net/http"
)


// type EmbeddingRequest struct {
//   Input string `json:"input"`
// }
//
// type EmbeddingResponse struct {
// 	Data[]struct {
// 			Embedding []float32 `json:"embedding"`
// 	} `json:"data"`
// }
//
// func GenerateEmbedding(text string) ([]float32, error) {
// 	reqBody := EmbeddingRequest{Input: text}
// 	jsonData, _ := json.Marshal(reqBody)
//
//
// 	resp, err := http.Post("http://localhost:11434/api/embeddings", "application/json", bytes.NewBuffer(jsonData))
//
// 	if err != nil {
// 			return nil, err
// 	}
// 	defer resp.Body.Close()
//
// 	var embResp EmbeddingResponse
// 	json.NewDecoder(resp.Body).Decode(&embResp)
//
// 	log.Printf("Ollama response: %+v", embResp)
// 	log.Printf("Data length: %d", len(embResp.Data))
//
// 	if len(embResp.Data[0].Embedding) == 0 {
// 		return nil, fmt.Errorf("empty embedding returned from ollama")
// 	}
//
// 	return embResp.Data[0].Embedding, nil
// }

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

	// Log what we're sending to debug
	log.Printf("Sending to Ollama: %s", string(jsonData))

	// resp, err := http.Post("http://localhost:11434/api/embeddings", "application/json", bytes.NewBuffer(jsonData))
	resp, err := http.Post("http://host.docker.internal:11434/api/embeddings", "application/josn", bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the raw response for debugging
	var rawResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&rawResp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	log.Printf("Raw Ollama response: %+v", rawResp)
	
	// Try to extract embedding from the raw response
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
