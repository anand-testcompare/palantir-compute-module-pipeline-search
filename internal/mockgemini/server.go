package mockgemini

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Server implements a minimal Gemini-API-like surface used by the Google GenAI SDK.
//
// It is intentionally deterministic and public-safe for hermetic tests.
type Server struct{}

func New() *Server {
	return &Server{}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handle)
	return mux
}

type generateContentRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	// Expected path form: /v1beta/models/<model>:generateContent
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !strings.HasPrefix(r.URL.Path, "/v1beta/models/") || !strings.HasSuffix(r.URL.Path, ":generateContent") {
		http.NotFound(w, r)
		return
	}
	if strings.TrimSpace(r.Header.Get("X-Goog-Api-Key")) == "" {
		http.Error(w, "missing api key", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read body", http.StatusBadRequest)
		return
	}

	email := extractEmail(body)
	domain := ""
	if at := strings.LastIndex(email, "@"); at >= 0 && at+1 < len(email) {
		domain = email[at+1:]
	}

	respText, _ := json.Marshal(map[string]string{
		"linkedin_url": "",
		"company":      domain,
		"title":        "",
		"description":  "",
		"confidence":   "mock",
	})

	response := map[string]any{
		"candidates": []any{
			map[string]any{
				"content": map[string]any{
					"parts": []any{
						map[string]any{"text": string(respText)},
					},
				},
				"groundingMetadata": map[string]any{
					"groundingChunks": []any{
						map[string]any{"web": map[string]any{
							"uri":   fmt.Sprintf("https://source.invalid/%s", domain),
							"title": "mock source",
						}},
					},
					"webSearchQueries": []any{
						fmt.Sprintf("company %s linkedin", domain),
					},
				},
				"urlContextMetadata": map[string]any{
					"urlMetadata": []any{
						map[string]any{
							"retrievedUrl":       fmt.Sprintf("https://context.invalid/%s", domain),
							"urlRetrievalStatus": "URL_RETRIEVAL_STATUS_SUCCESS",
						},
					},
				},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func extractEmail(body []byte) string {
	var req generateContentRequest
	_ = json.Unmarshal(body, &req)
	for _, c := range req.Contents {
		for _, p := range c.Parts {
			if strings.TrimSpace(p.Text) == "" {
				continue
			}
			// The production prompt includes "Email: <email>" on its own line.
			if idx := strings.Index(p.Text, "Email:"); idx >= 0 {
				rest := strings.TrimSpace(p.Text[idx+len("Email:"):])
				line, _, _ := strings.Cut(rest, "\n")
				return strings.TrimSpace(line)
			}
			// Fallback: try to find a token containing '@'.
			for _, tok := range strings.Fields(p.Text) {
				if strings.Contains(tok, "@") {
					return strings.Trim(strings.TrimSpace(tok), ",;()<>\"'")
				}
			}
		}
	}
	return ""
}
