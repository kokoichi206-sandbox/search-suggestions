package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/kokoichi206-sandbox/search-suggestions/api/model"
	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

const (
	// opensearch に関する情報。
	user = "admin"
	pass = "ad.PASS#1"
	addr = "https://opensearch:9201"

	index          = "searched-history"
	searchTemplate = "history-search-template"
)

func (h handler) searchSuggestions(ctx context.Context, index, template, query string) (*model.HistorySearchResult, error) {
	body, err := json.Marshal(model.HistorySearchRequest{
		ID: template,
		Params: model.Params{
			Query: query,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal: %w", err)
	}

	req := opensearchapi.SearchTemplateRequest{
		Index: []string{index},
		Body:  bytes.NewReader(body),
	}

	res, err := req.Do(ctx, h.opensearchClient)
	if err != nil {
		return nil, fmt.Errorf("failed to search document: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		return nil, fmt.Errorf("failed to execute opensearch request: %s", b)
	}

	var v model.HistorySearchResult

	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		return nil, fmt.Errorf("failed to decode: %w", err)
	}

	return &v, nil
}

type handler struct {
	opensearchClient *opensearch.Client
}

// NewHandler creates a new handler
//
// 簡単のため const に定義した値を使って初期化している。
func NewHandler() (*handler, error) {
	opensearchClient, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			// 自己証明書を許可する。
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses: []string{addr},
		Username:  user,
		Password:  pass,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create opensearch client: %w", err)
	}

	return &handler{
		opensearchClient: opensearchClient,
	}, nil
}

func (h handler) autoComplete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		if query == "" {
			http.Error(w, "query is required", http.StatusBadRequest)
			return
		}

		suggestions, err := h.searchSuggestions(
			r.Context(),
			index,
			searchTemplate,
			query,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res := model.SuggestResponse{
			Suggestions: []string{},
		}

		for _, sug := range suggestions.Aggregations.Keywords.Buckets {
			res.Suggestions = append(res.Suggestions, sug.Key)
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}

func main() {
	h, err := NewHandler()
	if err != nil {
		log.Fatal(err)
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /auto-complete", h.autoComplete())

	srv := &http.Server{
		Addr:              ":8085",
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
