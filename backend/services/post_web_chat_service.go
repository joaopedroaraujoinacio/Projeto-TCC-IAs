package services

import (
	"fmt"
	"go-project/models"
	"go-project/repositories"
	"go-project/utils"
	"strings"
)


type WebSearchService interface {
	WebSearchChatService(request *models.WebSearchRequest) (*models.WebSearchResponse, <-chan string, <-chan error)
}

type webSearchService struct {
	searchRepo utils.WebSearchRepository
	chatRepo   repositories.ChatRepository 
}

func NewWebSearchService(searchRepo utils.WebSearchRepository, chatRepo repositories.ChatRepository) WebSearchService {
	return &webSearchService{
		searchRepo: searchRepo,
		chatRepo:   chatRepo,
	}
}

func (s *webSearchService) WebSearchChatService(request *models.WebSearchRequest) (*models.WebSearchResponse, <-chan string, <- chan error) {
	if strings.TrimSpace(request.Query) == "" {
		errorChan := make(chan error, 1)
		messageChan := make(chan string)
		errorChan <- fmt.Errorf("search query cannot be empty")
		close(errorChan)
		close(messageChan)
		return nil, messageChan, errorChan
	}

	maxSources := request.MaxSources
	if maxSources <= 0 {
		maxSources = 3
	}

	sources, err := s.searchRepo.SearchDuckDuckGo(request.Query, maxSources)
	if err != nil {
		errorChan := make(chan error, 1)
		messageChan := make(chan string)
		errorChan <- fmt.Errorf("failed to search DuckDuckGo: %w", err)
		close(errorChan)
		close(messageChan)
		return nil, messageChan, errorChan
	}

	if len(sources) == 0 {
		response := &models.WebSearchResponse{
			Query:     request.Query,
			AISummary: "No results found for your query.",
			Sources:   []models.SearchSource{},
			Count:     0,
		}
		errorChan := make(chan error, 1)
		messageChan := make(chan string, 1)
		messageChan <- "No results found for your query"
		close(errorChan)
		close(messageChan)
		return response, messageChan, errorChan
	}

	for i := range sources {
		if sources[i].URL != "" && len(sources[i].Content) < 200 {
			if content, err := s.searchRepo.ScrapeContent(sources[i].URL); err == nil {
				sources[i].Content = content
			}
		}
	}

	messageChan, errorChan := s.generateAISummary(request.Query, sources)

	cleanSources := make([]models.SearchSource, len(sources))
	for i := range sources {
		cleanSources[i] = sources[i]
		cleanSources[i].Content = ""
	}

	response := &models.WebSearchResponse{
		Query:     request.Query,
		AISummary: "",
		Sources:   sources,
		Count:     len(sources),
	}
	return response, messageChan, errorChan
}

func (s *webSearchService) generateAISummary(query string, sources []models.SearchSource) (<-chan string, <-chan error) {
	messageChan := make(chan string, 10)
	errorChan := make(chan error, 1)


	context := fmt.Sprintf("Based on the following search results, provide a comprehensive answer to the query: '%s'\n\n", query)
	
	for i, source := range sources {
		context += fmt.Sprintf("Source %d (%s):\n%s\n\n", i+1, source.Title, source.Content)
	}
	
	context += fmt.Sprintf("Please provide a well-structured, informative answer to: '%s'", query)

	chatRequest := &models.ChatRequest{
		Message: context,
		// Model:   "llama3.2:3b",
	}

	go func() {
		defer close(messageChan)
		defer close(errorChan)
	
		streamChan, err := s.chatRepo.SendToLLM(chatRequest)
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

