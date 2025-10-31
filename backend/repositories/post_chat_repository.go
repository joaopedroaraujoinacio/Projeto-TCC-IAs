package repositories

import (
	"fmt"
	"bufio"
	"bytes"
	"net/http"
	"encoding/json"
	"go-project/models"
)


func (r *chatRepository) SendToLLM(request *models.ChatRequest) (<-chan models.StreamChunk, error) {
    model := request.Model
    if model == "" {
        model = "llama3.2:3b"
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

    fmt.Printf("Sending to Ollama: %s\n", string(jsonData)) 

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
            
            fmt.Printf("Ollama line %d: %s\n", lineCount, line) 
            
            if line == "" {
                continue
            }

            var ollamaResp models.OllamaResponse
            if err := json.Unmarshal([]byte(line), &ollamaResp); err != nil {
                fmt.Printf("Parse error on line %d: %v\n", lineCount, err) 
                streamChan <- models.StreamChunk{
                    Error: fmt.Errorf("failed to parse response: %w", err),
                }
                continue
            }

            fmt.Printf("Parsed chunk: content='%s', done=%v\n", 
                ollamaResp.Message.Content, ollamaResp.Done) 

            streamChan <- models.StreamChunk{
                Text:  ollamaResp.Message.Content,
                Done:  ollamaResp.Done,
            }

            if ollamaResp.Done {
                break
            }
        }

        if err := scanner.Err(); err != nil {
            fmt.Printf("Scanner error: %v\n", err) 
            streamChan <- models.StreamChunk{
                Error: fmt.Errorf("scanner error: %w", err),
            }
        }
        
        fmt.Printf("Stream ended. Total lines: %d\n", lineCount) 
    }()

    return streamChan, nil
}

