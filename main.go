package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	qdrantutils "github.com/AstraBert/search-web-app/qdrant-utils"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("POST /search", handleSearch)

	fmt.Println("Starting server on port 8000")
	http.ListenAndServe(":8000", mux)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

type SearchArgs struct {
	Query string `json:"query"`
	Limit uint64 `json:"limit"`
}

type Results struct {
	SearchResults []qdrantutils.SearchResult `json:"results"`
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	var searchArgs SearchArgs
	err := json.NewDecoder(r.Body).Decode(&searchArgs)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}
	searchResults, err := qdrantutils.SearchText(searchArgs.Query, &searchArgs.Limit)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(Results{SearchResults: searchResults})
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
