package enrich

import (
	"context"
	"errors"
	"strings"
)

// Stub is a deterministic, no-network enricher used for hermetic local runs and tests.
type Stub struct{}

func (Stub) Enrich(_ context.Context, email string) (Result, error) {
	// Force errors for sentinel domains to exercise failure paths without external calls.
	if strings.HasSuffix(strings.ToLower(strings.TrimSpace(email)), "@error.test") {
		return Result{}, errors.New("stub enricher: forced error")
	}

	domain := ""
	if at := strings.LastIndex(email, "@"); at >= 0 && at+1 < len(email) {
		domain = email[at+1:]
	}

	return Result{
		LinkedInURL: "",
		Company:     domain,
		Title:       "",
		Description: "",
		Confidence:  "stub",
	}, nil
}
