package repositories

import (
	"fmt"
	"time"
	"bufio"
	"bytes"
	"net/http"
	"encoding/json"
	"go-project/models"
)

type ChatRepository interface {
    SendToLLM(request *models.ChatRequest) (<-chan models.StreamChunk, error)
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
    client    *http.Client
}

func (r *chatRepository) SendToLLM(request *models.ChatRequest) (<-chan models.StreamChunk, error) {
    model := request.Model
    if model == "" {
        model = "llama3.2:3b"
    }

    // Build messages array with history + current message
    messages := []map[string]string{}
    
    // Add all history
    for _, msg := range request.History {
        messages = append(messages, map[string]string{
            "role":    msg.Role,
            "content": msg.Content,
        })
    }
    
    // Add current message
    messages = append(messages, map[string]string{
        "role":    "user",
        "content": request.Message,
    })

    fmt.Printf("Building request with %d messages (history: %d, current: 1)\n", 
        len(messages), len(request.History))

    ollamaReq := models.OllamaRequest{
        Model:    model,
        Messages: messages,
        Stream:   true,
    }

    jsonData, err := json.Marshal(ollamaReq)
    if err != nil {
        return nil, fmt.Errorf("failed to format request: %w", err)
    }

    fmt.Printf("Sending to Ollama: %s\n", string(jsonData)) // DEBUG

    resp, err := r.client.Post(r.ollamaURL+"/api/chat", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to send request to ollama: %w", err)
    }

    if resp.StatusCode != http.StatusOK {
        resp.Body.Close()
        return nil, fmt.Errorf("failed to get response from ollama: %d", resp.StatusCode)
    }

    streamChan := make(chan models.StreamChunk, 10)

    go func() {
        defer close(streamChan)
        defer resp.Body.Close()
        scanner := bufio.NewScanner(resp.Body)

        lineCount := 0
        for scanner.Scan() {
            line := scanner.Text()
            lineCount++
            
            fmt.Printf("Ollama line %d: %s\n", lineCount, line) // DEBUG
            
            if line == "" {
                continue
            }

            var ollamaResp models.OllamaResponse
            if err := json.Unmarshal([]byte(line), &ollamaResp); err != nil {
                fmt.Printf("Parse error on line %d: %v\n", lineCount, err) // DEBUG
                streamChan <- models.StreamChunk{
                    Error: fmt.Errorf("failed to parse response: %w", err),
                }
                continue
            }

            fmt.Printf("Parsed chunk: content='%s', done=%v\n", 
                ollamaResp.Message.Content, ollamaResp.Done) // DEBUG

            streamChan <- models.StreamChunk{
                Text:  ollamaResp.Message.Content,
                Done:  ollamaResp.Done,
            }

            if ollamaResp.Done {
                break
            }
        }

        if err := scanner.Err(); err != nil {
            fmt.Printf("Scanner error: %v\n", err) // DEBUG
            streamChan <- models.StreamChunk{
                Error: fmt.Errorf("scanner error: %w", err),
            }
        }
        
        fmt.Printf("Stream ended. Total lines: %d\n", lineCount) // DEBUG
    }()

    return streamChan, nil
}

// type ChatRepository interface {
// 	SendToLLM(request *models.ChatRequest) (<-chan models.StreamChunk, error)
// }
//
// func NewChatRepository(ollamaURL string) ChatRepository {
// 	return &chatRepository{
// 		ollamaURL: ollamaURL,
// 		client: &http.Client{
// 			Timeout: 500 * time.Second,
// 		},
// 	}
// }
//
// type chatRepository struct {
// 	ollamaURL string
// 	client	*http.Client
// }
//
// func (r *chatRepository) SendToLLM(request *models.ChatRequest) (<-chan models.StreamChunk, error) {
// 	model := request.Model
// 	if model == "" {
// 		model = "llama3.2:3b"
// 	}
//
// 	ollamaReq := models.OllamaRequest{
// 		Model:  model,
// 		Prompt: request.Message,
// 		Stream: true, 
// 	}
//
// 	jsonData, err := json.Marshal(ollamaReq)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to format request: %w", err)
// 	}
//
// 	resp, err := r.client.Post(r.ollamaURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to send request to ollama: %w", err)
// 	}
//
// 	if resp.StatusCode != http.StatusOK {
// 		resp.Body.Close()
// 		return nil, fmt.Errorf("failed to get response from ollama: %d", resp.StatusCode)
// 	}
//
// 	streamChan := make(chan models.StreamChunk, 10)
//
// 	go func() {
// 		defer close(streamChan)      
// 		defer resp.Body.Close()     
//
// 		scanner := bufio.NewScanner(resp.Body)
//
// 		for scanner.Scan() {
// 			line := scanner.Text()
// 			if line == "" {
// 				continue 
// 			}
//
// 			var ollamaResp models.OllamaResponse
// 			if err := json.Unmarshal([]byte(line), &ollamaResp); err != nil {
// 				streamChan <- models.StreamChunk{
// 					Error: fmt.Errorf("failed to parse response: %w", err),
// 				}
// 				continue
// 			}
//
// 			streamChan <- models.StreamChunk{
// 				Text:  ollamaResp.Response, 
// 				Done:  ollamaResp.Done,
// 			}
//
// 			if ollamaResp.Done {
// 				break
// 			}
// 		}
//
// 		if err := scanner.Err(); err != nil {
// 			streamChan <- models.StreamChunk{
// 				Error: fmt.Errorf("scanner error: %w", err),
// 			}
// 		}
// 	}()
//
// 	return streamChan, nil
// }
//
