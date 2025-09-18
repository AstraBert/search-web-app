package qdrantutils

import (
	"os"
	"testing"
)

func TestSearchText(t *testing.T) {
	_, okOpenAi := os.LookupEnv("OPENAI_API_KEY")
	_, okQdApi := os.LookupEnv("QDRANT_API_KEY")
	_, okQdUrl := os.LookupEnv("QDRANT_ENDPOINT")
	if !okOpenAi || !okQdApi || !okQdUrl {
		t.Skip("Skipping test because the environment variables are not configured correctly.")
	}
	testCases := []struct {
		text  string
		limit uint64
	}{
		{text: "Sommerzeit", limit: 5},
		{text: "Nachbar", limit: 3},
	}
	for _, tc := range testCases {
		results, err := SearchText(tc.text, &tc.limit)
		if err != nil {
			t.Errorf("Expecting no error, got %s", err.Error())
		}
		if len(results) != int(tc.limit) {
			t.Errorf("Expecting results to be of length %d, but they are of length %d", int(tc.limit), len(results))
		}
	}
}
