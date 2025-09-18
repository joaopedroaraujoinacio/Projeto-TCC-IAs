package repositories

import (
	"io"
	"fmt"
	"time"
	"regexp"
	"strings"
	"net/url"
	"net/http"
	"encoding/json"
	"golang_crud/models"
)

type SmartSearchRepository interface {
	SearchDuckDuckGo(query string, maxSources int) ([]models.SearchSource, error)
	ScrapeContent(url string) (string, error)
}

type smartSearchRepository struct {
	client *http.Client
}

func NewSmartSearchRepository() SmartSearchRepository {
	fmt.Println("Using direct DuckDuckGo connection")
	return &smartSearchRepository{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (r *smartSearchRepository) SearchDuckDuckGo(query string, maxSources int) ([]models.SearchSource, error) {
	fmt.Printf("DEBUG: Received query: '%s', maxSources: %d\n", query, maxSources)
	
	if maxSources <= 0 {
		maxSources = 3
	}

	if strings.TrimSpace(query) == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	encodedQuery := url.QueryEscape(query)
	fmt.Printf("DEBUG: Encoded query: '%s'\n", encodedQuery)
	
	// Try DuckDuckGo Instant Answer API first
	searchURL := fmt.Sprintf("https://api.duckduckgo.com/?q=%s&format=json&no_html=1&skip_disambig=1", encodedQuery)
	fmt.Printf("DEBUG: Search URL: %s\n", searchURL)
	
	var resp *http.Response
	var err error
	
	// Retry logic for 202 responses
	for i := range []int{0, 1, 2} {  
	// for i := 0; i < 3; i++ {
		resp, err = r.client.Get(searchURL)
		if err != nil {
			return nil, fmt.Errorf("failed to search DuckDuckGo: %w", err)
		}
		
		if resp.StatusCode == http.StatusOK {
			break
		} else if resp.StatusCode == 202 {
			resp.Body.Close()
			fmt.Printf("DEBUG: Got 202, retrying in %d seconds...\n", i+1)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		} else {
			resp.Body.Close()
			return nil, fmt.Errorf("DuckDuckGo returned status: %d", resp.StatusCode)
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var ddgResp models.DuckDuckGoResponse
	if err := json.Unmarshal(body, &ddgResp); err != nil {
		return nil, fmt.Errorf("failed to parse DuckDuckGo response: %w", err)
	}

	sources	:= r.buildSources(ddgResp, query, encodedQuery, maxSources)

	return sources, nil
}

func (r *smartSearchRepository) buildSources(
	ddgResp models.DuckDuckGoResponse,
	query, encodedQuery string, 
	maxSources int,
) []models.SearchSource {

	var sources []models.SearchSource
	count := 0

	// Add abstract/instant answer if available
	if ddgResp.AbstractText != "" && count < maxSources {
		fmt.Printf("DEBUG: Found abstract: %s\n", ddgResp.AbstractText[:min(50, len(ddgResp.AbstractText))])
		sources = append(sources, models.SearchSource{
			Title:   r.extractTitle(ddgResp.AbstractText),
			URL:     ddgResp.AbstractURL,
			Summary: ddgResp.AbstractText,
			Content: ddgResp.AbstractText,
		})
		count++
	}

	// Add related topics
	for _, topic := range ddgResp.RelatedTopics {
		if count >= maxSources {
			break
		}
		if topic.FirstURL != "" && topic.Text != "" {
			fmt.Printf("DEBUG: Found related topic: %s\n", topic.Text[:min(50, len(topic.Text))])
			sources = append(sources, models.SearchSource{
				Title:   r.extractTitle(topic.Text),
				URL:     topic.FirstURL,
				Summary: topic.Text,
				Content: topic.Text,
			})
			count++
		}
	}

	if len(sources) == 0 {
		fmt.Println("DEBUG: No instant answers found, creating fallback source")
		sources = append(sources, models.SearchSource{
			Title:   fmt.Sprintf("Information about: %s", query),
			URL:     fmt.Sprintf("https://duckduckgo.com/?q=%s", encodedQuery),
			Summary: fmt.Sprintf("General information about: %s", query),
			Content: fmt.Sprintf("Please provide information about: %s", query),
		})
	}

	fmt.Printf("DEBUG: Returning %d sources\n", len(sources))
	return sources
}

func (r *smartSearchRepository) ScrapeContent(url string) (string, error) {
	resp, err := r.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to scrape content: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("scraping returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read scraped content: %w", err)
	}

	// Simple HTML cleaning (basic implementation)
	content := string(body)
	content = r.cleanHTML(content)
	
	// Limit content size for LLM
	if len(content) > 2000 {
		content = content[:2000] + "..."
	}

	return content, nil
}

func (r *smartSearchRepository) extractTitle(text string) string {
	parts := strings.Split(text, " - ")
	if len(parts) > 0 {
		title := strings.TrimSpace(parts[0])
		if len(title) > 80 {
			return title[:80] + "..."
		}
		return title
	}
	
	if len(text) > 50 {
		return text[:50] + "..."
	}
	return text
}

func (r *smartSearchRepository) cleanHTML(html string) string {
	// Remove HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	text := re.ReplaceAllString(html, " ")
	
	// Remove extra whitespace
	re = regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")
	
	return strings.TrimSpace(text)
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
