package enrich

import (
	"context"
)

// Result is the structured enrichment output for a single email.
//
// MVP: everything is a string to keep CSV output simple and stable.
type Result struct {
	LinkedInURL string
	Company     string
	Title       string
	Description string
	Confidence  string
}

// Enricher enriches a single email address.
type Enricher interface {
	Enrich(ctx context.Context, email string) (Result, error)
}
