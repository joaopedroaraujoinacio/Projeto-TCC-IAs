package utils

import (
	"encoding/json"
	"fmt"
	"go-project/models"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type WebSearchRepository interface {
	SearchDuckDuckGo(query string, maxSources int) ([]models.SearchSource, error)
	ScrapeContent(url string) (string, error)
}

type webSearchRepository struct {
	client     *http.Client
	searxngURL string
}

func NewWebSearchRepository() WebSearchRepository {
	searxngURL := os.Getenv("SEARXNG_URL")
	if searxngURL == "" {
		searxngURL = "http://localhost:8080"
	}

	fmt.Printf("Using SearXNG at: %s\n", searxngURL)

	return &webSearchRepository{
		client:     &http.Client{Timeout: 30 * time.Second},
		searxngURL: searxngURL,
	}
}

type SearXNGResponse struct {
	Query   string          `json:"query"`
	Results []SearXNGResult `json:"results"`
}

type SearXNGResult struct {
	URL     string `json:"url"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (r *webSearchRepository) SearchDuckDuckGo(query string, maxSources int) ([]models.SearchSource, error) {
	fmt.Printf("DEBUG: Received query: '%s', maxSources: %d\n", query, maxSources)

	if maxSources <= 0 {
		maxSources = 5
	}

	if strings.TrimSpace(query) == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	params := url.Values{}
	params.Add("q", query)
	params.Add("format", "json")

	searchURL := fmt.Sprintf("%s/search?%s", r.searxngURL, params.Encode())
	fmt.Printf("DEBUG: Search URL: %s\n", searchURL)

	resp, err := r.client.Get(searchURL)
	if err != nil {
		fmt.Printf("DEBUG: Request failed with error: %v\n", err)
		return nil, fmt.Errorf("failed to search: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("DEBUG: Got response with status: %d\n", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// CRITICAL DEBUG LINE
	fmt.Printf("DEBUG: Response body length: %d, content: %s\n", len(body), string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search returned status %d: %s", resp.StatusCode, string(body))
	}

	var searxngResp SearXNGResponse
	if err := json.Unmarshal(body, &searxngResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w. Body: %s", err, string(body))
	}

	fmt.Printf("DEBUG: Found %d results\n", len(searxngResp.Results))

	if len(searxngResp.Results) == 0 {
		return nil, fmt.Errorf("no results found for query: %s", query)
	}

	sources := make([]models.SearchSource, 0, maxSources)
	seen := make(map[string]bool)

	for _, result := range searxngResp.Results {
		if len(sources) >= maxSources {
			break
		}

		if seen[result.URL] {
			continue
		}
		seen[result.URL] = true

		if !strings.HasPrefix(result.URL, "http://") && !strings.HasPrefix(result.URL, "https://") {
			continue
		}

		fmt.Printf("DEBUG: Found result: %s\n", result.Title)

		sources = append(sources, models.SearchSource{
			Title:   result.Title,
			URL:     result.URL,
			Summary: result.Content,
			Content: result.Content,
		})
	}

	if len(sources) == 0 {
		return nil, fmt.Errorf("no valid results found")
	}

	fmt.Printf("DEBUG: Returning %d sources\n", len(sources))
	return sources, nil
}


func (r *webSearchRepository) ScrapeContent(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Bot/1.0)")

	resp, err := r.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to scrape content: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("scraping returned status: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("unsupported content type: %s", contentType)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	var content strings.Builder
	r.extractText(doc, &content)

	text := content.String()
	text = r.cleanHTML(text)

	if len(text) > 5000 {
		text = text[:5000] + "..."
	}

	if len(text) < 50 {
		return "", fmt.Errorf("insufficient content extracted")
	}

	return text, nil
}

func (r *webSearchRepository) extractText(n *html.Node, buf *strings.Builder) {
	if n.Type == html.TextNode {
		text := strings.TrimSpace(n.Data)
		if text != "" {
			buf.WriteString(text)
			buf.WriteString(" ")
		}
	}

	if n.Type == html.ElementNode {
		if n.Data == "script" || n.Data == "style" || n.Data == "noscript" {
			return
		}

		if n.Data == "p" || n.Data == "div" || n.Data == "br" {
			buf.WriteString("\n")
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		r.extractText(c, buf)
	}
}

func (r *webSearchRepository) cleanHTML(html string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	text := re.ReplaceAllString(html, " ")

	re = regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")

	return strings.TrimSpace(text)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
