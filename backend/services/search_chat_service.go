package services


import (

	"fmt"
	"strings"
	"golang_crud/models"
	"golang_crud/repositories"



)


type SmartSearchService interface {
	ProcessSmartSearch(request *models.SmartSearchRequest) (*models.SmartSearchResponse, error)
}

type smartSearchService struct {
	searchRepo repositories.SmartSearchRepository 
	chatRepo   repositories.ChatRepository 
}

func NewSmartSearchService(searchRepo repositories.SmartSearchRepository, chatRepo repositories.ChatRepository) SmartSearchService {
	return &smartSearchService{
		searchRepo: searchRepo,
		chatRepo:   chatRepo,
	}
}

func (s *smartSearchService) ProcessSmartSearch(request *models.SmartSearchRequest) (*models.SmartSearchResponse, error) {
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
		return &models.SmartSearchResponse{
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

	return &models.SmartSearchResponse{
		Query:     request.Query,
		AISummary: aiSummary,
		Sources:   sources,
		Count:     len(sources),
	}, nil
}

func (s *smartSearchService) generateAISummary(query string, sources []models.SearchSource) (string, error) {
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

