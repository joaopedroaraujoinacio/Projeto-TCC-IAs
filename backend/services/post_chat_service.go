package services

import (
	"fmt"
	"go-project/models"
	"strings"
)


func (s *chatService) StreamChat(request *models.ChatRequest) (<-chan string, <-chan error) {
	messageChan := make(chan string, 10)
	errorChan := make(chan error, 1)

	if strings.TrimSpace(request.Message) == "" {
		errorChan <- fmt.Errorf("message cannot be empty")
		close(messageChan)
		close(errorChan)
		return messageChan, errorChan
	}

	go func() {
		defer close(messageChan)
		defer close(errorChan)

		streamChan, err := s.chatRepo.SendToLLM(request)
		if err != nil {
			errorChan <- fmt.Errorf("failed to start streaming: %w", err)
			return
		}

		buffer := ""

		for chunk := range streamChan {
			if chunk.Error != nil {
				errorChan <- chunk.Error
				return
			}

			buffer += chunk.Text

			for {
				spaceIdx := strings.Index(buffer, " ")
				if spaceIdx == -1 {
					break
				}
				word := buffer[:spaceIdx+1] 
				messageChan <- word
				buffer = buffer[spaceIdx+1:]
			}

			if chunk.Done {
				if buffer != "" {
					messageChan <- buffer
				}
				return
			}
		}
	}()

	return messageChan, errorChan
}
