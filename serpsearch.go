package goserp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rdhillbb/util"
	"io"
	"os"
	"errors"
	"net/http"
	"time"
)

// SearchResult represents the unified structure for results
type SearchResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

// PeopleAlsoAskItem represents a Q&A item
type PeopleAlsoAskItem struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Source   string `json:"source"`
	URL      string `json:"url"`
}

// QueryResult represents results for a single query
type QueryResult struct {
	Query           string              `json:"query"`
	AnswerBox       *SearchResult       `json:"answerBox,omitempty"`
	OrganicResults  []SearchResult      `json:"organicResults"`
	PeopleAlsoAsk   []PeopleAlsoAskItem `json:"peopleAlsoAsk"`
	RelatedSearches []string            `json:"relatedSearches"`
}

// SERPResponse represents the complete response structure
type SERPResponse struct {
	StructureType string        `json:"structureType"`
	Timestamp     string        `json:"timestamp"`
	Credits       int           `json:"credits"`
	Queries       []QueryResult `json:"queries"`
}

// Name:SerpSearch
// Description: SerpSearh should be called for basic Internet search
// Required Parmater query:string
func SerpSearch(query string) (string, error) {
	queries := []string{query} // Changed from array to slice
	results, err := serpSearch(queries)
	if err != nil {
		return "", err // Fixed error variable name
	}
	jsonData, err := json.MarshalIndent(results, "", "  ")
	return string(jsonData), err // Fixed inconsistent error return
}

// Name:SerpExtensiveSearch
// Description: SerpExtensiveSearh should be called for deep extensive search is requested.
// Required Parmater query:string
func SerpExtensiveSearch(query string) (string, error) {
	queries, err := util.ReWriteQR(query)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	results, err1 := serpSearch(queries)
	if err1 != nil {
		return "", err
	}
	jsonData, err := json.MarshalIndent(results, "", "  ")
	return string(jsonData), err1
}

// SerpSearch performs search on a list of searched provided.
// The input parameters is string array.
func serpSearch(queries []string) (*SERPResponse, error) {
	    value := os.Getenv("SERP_API_KEY")
    if value == "" {
	     return nil, errors.New("SERP_API_KEY environment variable is not set")
    }
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	response := &SERPResponse{
		StructureType: "serp",
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		Credits:       0,
		Queries:       make([]QueryResult, 0, len(queries)),
	}

	for _, query := range queries {
		// Create maps to track unique URLs for this query
		seenOrganicURLs := make(map[string]SearchResult)
		seenPAAURLs := make(map[string]PeopleAlsoAskItem)

		reqBody, err := json.Marshal(map[string]interface{}{
			"q":   query,
			"num": 40,
		})
		if err != nil {
			return nil, fmt.Errorf("error marshaling request: %w", err)
		}

		req, err := http.NewRequest("POST", "https://google.serper.dev/search", bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, fmt.Errorf("error creating request: %w", err)
		}

		req.Header.Set("X-API-KEY", "5ec69c6ba87e1fc14dc4382c060d3d190e1cbcff")
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response: %w", err)
		}

		var rawResponse map[string]interface{}
		if err := json.Unmarshal(body, &rawResponse); err != nil {
			return nil, fmt.Errorf("error parsing response: %w", err)
		}

		// Extract credits if available
		if credits, ok := rawResponse["credits"].(float64); ok {
			response.Credits += int(credits)
		}

		queryResult := QueryResult{
			Query:           query,
			OrganicResults:  []SearchResult{},
			PeopleAlsoAsk:   []PeopleAlsoAskItem{},
			RelatedSearches: []string{},
		}

		// Process answer box
		if answerBox, ok := rawResponse["answerBox"].(map[string]interface{}); ok {
			result := SearchResult{
				Title:   answerBox["title"].(string),
				URL:     answerBox["link"].(string),
				Snippet: answerBox["snippet"].(string),
			}
			queryResult.AnswerBox = &result
			// Also add answer box to organic results map to prevent duplication
			seenOrganicURLs[result.URL] = result
		}

		// Process organic results
		if organic, ok := rawResponse["organic"].([]interface{}); ok {
			for _, result := range organic {
				res := result.(map[string]interface{})
				url := res["link"].(string)
				searchResult := SearchResult{
					Title:   res["title"].(string),
					URL:     url,
					Snippet: res["snippet"].(string),
				}
				seenOrganicURLs[url] = searchResult
			}
		}

		// Process People Also Ask
		if paa, ok := rawResponse["peopleAlsoAsk"].([]interface{}); ok {
			for _, item := range paa {
				qa := item.(map[string]interface{})
				url := qa["link"].(string)
				paaItem := PeopleAlsoAskItem{
					Question: qa["question"].(string),
					Answer:   qa["snippet"].(string),
					Source:   qa["title"].(string),
					URL:      url,
				}
				seenPAAURLs[url] = paaItem
			}
		}

		// Convert maps to slices for the query result
		for _, result := range seenOrganicURLs {
			queryResult.OrganicResults = append(queryResult.OrganicResults, result)
		}

		for _, item := range seenPAAURLs {
			queryResult.PeopleAlsoAsk = append(queryResult.PeopleAlsoAsk, item)
		}

		// Process Related Searches (these are already unique by nature)
		if related, ok := rawResponse["relatedSearches"].([]interface{}); ok {
			seenRelated := make(map[string]bool)
			for _, r := range related {
				if search, ok := r.(map[string]interface{}); ok {
					query := search["query"].(string)
					if !seenRelated[query] {
						queryResult.RelatedSearches = append(queryResult.RelatedSearches, query)
						seenRelated[query] = true
					}
				}
			}
		}

		response.Queries = append(response.Queries, queryResult)
	}

	return response, nil
}
