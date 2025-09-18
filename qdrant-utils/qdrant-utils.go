package qdrantutils

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/AstraBert/search-web-app/embeddings"
	"github.com/qdrant/go-client/qdrant"
)

type SearchResult struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func SearchText(text string, limit *uint64) ([]SearchResult, error) {
	qdrantClient, errQd := qdrant.NewClient(&qdrant.Config{
		Host:   os.Getenv("QDRANT_ENDPOINT"),
		Port:   6334,
		APIKey: os.Getenv("QDRANT_API_KEY"),
		UseTLS: true,
		Cloud:  true,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
		SkipCompatibilityCheck: true,
	})

	if errQd != nil {
		return nil, errQd
	}

	vec, errVec := embeddings.EmbedText(text)

	if errVec != nil {
		return nil, errVec
	}

	vectorSlice := make([]float32, len(vec))
	for k, v := range vec {
		vectorSlice[k] = float32(v)
	}

	searchResult, errSearch := qdrantClient.Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: "germanSnippets",
		Query:          qdrant.NewQuery(vectorSlice...),
		WithPayload:    qdrant.NewWithPayload(true),
		Limit:          limit,
	})

	if errSearch != nil {
		return nil, errSearch
	}

	results := []SearchResult{}

	for _, point := range searchResult {
		title, okTit := point.Payload["title"]
		content, okExp := point.Payload["explanation"]
		if !okTit || !okExp {
			fmt.Println("Continuing because either title or explanation are not in payload")
			continue
		}
		results = append(results, SearchResult{Title: title.GetStringValue(), Content: content.GetStringValue()})
	}

	return results, nil
}
