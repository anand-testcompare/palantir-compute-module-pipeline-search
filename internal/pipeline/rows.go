package pipeline

import (
	"context"
	"strings"

	"github.com/palantir/palantir-compute-module-pipeline-search/internal/enrich"
)

// Row is the stable output schema contract for the MVP.
type Row struct {
	Email       string
	LinkedInURL string
	Company     string
	Title       string
	Description string
	Confidence  string
	Status      string
	Error       string
}

// Header returns the stable CSV header for Row.
func Header() []string {
	return []string{
		"email",
		"linkedin_url",
		"company",
		"title",
		"description",
		"confidence",
		"status",
		"error",
	}
}

// EnrichEmails runs the enricher over all emails and returns stable output rows.
//
// Errors from enrichment are recorded per-row and do not fail the full run.
func EnrichEmails(ctx context.Context, emails []string, enricher enrich.Enricher) ([]Row, error) {
	rows := make([]Row, 0, len(emails))
	for _, raw := range emails {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		email := strings.TrimSpace(raw)
		if email == "" {
			rows = append(rows, Row{
				Email:  "",
				Status: "error",
				Error:  "empty email",
			})
			continue
		}

		res, err := enricher.Enrich(ctx, email)
		if err != nil {
			rows = append(rows, Row{
				Email:  email,
				Status: "error",
				Error:  err.Error(),
			})
			continue
		}

		rows = append(rows, Row{
			Email:       email,
			LinkedInURL: res.LinkedInURL,
			Company:     res.Company,
			Title:       res.Title,
			Description: res.Description,
			Confidence:  res.Confidence,
			Status:      "ok",
			Error:       "",
		})
	}
	return rows, nil
}
