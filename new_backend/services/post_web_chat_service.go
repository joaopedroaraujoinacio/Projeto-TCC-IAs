package services

import (
	"fmt"
	"go-project/models"
	"go-project/repositories"
	"go-project/utils"
	"strings"
)


type WebSearchService interface {
	WebSearchChatService(request *models.WebSearchRequest) (*models.WebSearchResponse, error)
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

func (s *webSearchService) WebSearchChatService(request *models.WebSearchRequest) (*models.WebSearchResponse, error) {
	if strings.TrimSpace(request.Query) == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	maxSources := request.MaxSources
	if maxSources <= 0 {
		maxSources = 3
	}

	sources, err := s.searchRepo.SearchDuckDuckGo(request.Query, maxSources)
	if err != nil {
		return nil, fmt.Errorf("failed to search DuckDuckGo: %w", err)
	}

	if len(sources) == 0 {
		return &models.WebSearchResponse{
			Query:     request.Query,
			AISummary: "No results found for your query.",
			Sources:   []models.SearchSource{},
			Count:     0,
		}, nil
	}

	for i := range sources {
		if sources[i].URL != "" && len(sources[i].Content) < 200 {
			if content, err := s.searchRepo.ScrapeContent(sources[i].URL); err == nil {
				sources[i].Content = content
			}
		}
	}

	aiSummary, err := s.generateAISummary(request.Query, sources)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AI summary: %w", err)
	}

	for i := range sources {
		sources[i].Content = ""
	}

	return &models.WebSearchResponse{
		Query:     request.Query,
		AISummary: aiSummary,
		Sources:   sources,
		Count:     len(sources),
	}, nil
}

func (s *webSearchService) generateAISummary(query string, sources []models.SearchSource) (string, error) {
	context := fmt.Sprintf("Based on the following search results, provide a comprehensive answer to the query: '%s'\n\n", query)
	
	for i, source := range sources {
		context += fmt.Sprintf("Source %d (%s):\n%s\n\n", i+1, source.Title, source.Content)
	}
	
	context += fmt.Sprintf("Please provide a well-structured, informative answer to: '%s'", query)

	chatRequest := &models.ChatRequest{
		Message: context,
		Model:   "llama3.2:3b",
	}

	response, err := s.chatRepo.SendToLLM(chatRequest)
	if err != nil {
		return "", fmt.Errorf("failed to get LLM summary: %w", err)
	}

	return response.Response, nil
}

