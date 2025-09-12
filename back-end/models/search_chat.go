package models


type SmartSearchRequest struct {
	Query      string `json:"query" binding:"required"`
	MaxSources int    `json:"max_sources,omitempty"`

}


type SmartSearchResponse struct {
	Query     string         `json:"query"`
	AISummary string         `json:"ai_summary"`
	Sources   []SearchSource `json:"sources"`
	Count     int           `json:"count"`
}

type SearchSource struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Summary string `json:"summary"`
	Content string `json:"-"` 
}

type DuckDuckGoResponse struct  {
	AbstractText   string                `json:"AbstractText"`
	AbstractSource string                `json:"AbstractSource"`
	AbstractURL    string                `json:"AbstractURL"`
	RelatedTopics  []DuckDuckGoRelated   `json:"RelatedTopics"`
	Results        []DuckDuckGoInstant   `json:"Results"`
}

type DuckDuckGoRelated struct {
	FirstURL string `json:"FirstURL"`
	Text     string `json:"Text"`
}

type DuckDuckGoInstant struct {
	FirstURL string `json:"FirstURL"`
	Text     string `json:"Text"`
}

type LLMSummaryRequest struct {
	Query   string         `json:"query"`
	Sources []SearchSource `json:"sources"`
}

