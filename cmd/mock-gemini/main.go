package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/palantir/palantir-compute-module-pipeline-search/internal/mockgemini"
)

func main() {
	addr := defaultString("MOCK_GEMINI_ADDR", ":8081")

	fs := flag.NewFlagSet("mock-gemini", flag.ExitOnError)
	fs.StringVar(&addr, "addr", addr, "Listen address")
	_ = fs.Parse(os.Args[1:])

	srv := mockgemini.New()

	_, _ = fmt.Fprintf(os.Stdout, "mock-gemini listening on %s\n", addr)
	if err := http.ListenAndServe(addr, srv.Handler()); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func defaultString(envVar string, fallback string) string {
	v := strings.TrimSpace(os.Getenv(envVar))
	if v == "" {
		return fallback
	}
	return v
}
