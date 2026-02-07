# palantir-compute-module-pipeline-search

Pipeline-mode Foundry Compute Module (Go) that:

1. Reads a dataset of email addresses
2. Enriches each email via Gemini (Google Search grounding + URL context + structured output)
3. Writes an output dataset

This is a one-shot batch container triggered by Foundry pipeline builds (not a long-lived service).

It should also be runnable locally (without Foundry) against a local input file for faster iteration and personal one-off batches.

## Development

Verify (format, checks, tests):

```
./godelw verify
```

Run locally (no Foundry required, Gemini required):

```
export GEMINI_API_KEY=...
export GEMINI_MODEL=gemini-2.5-flash
go run ./cmd/enricher local --input /path/to/emails.csv --output /path/to/enriched.csv
```

Run Foundry-like end-to-end locally (mock dataset API + mock Gemini + real container):

```
docker compose -f docker-compose.test.yml up --abort-on-container-exit --build
```

Run black-box E2E via Venom (generates an HTML report under `out/venom/`):

```
bash test/scripts/run_venom_e2e_local.sh
```

Run the local harness with your own input CSV (not committed):

```
docker compose -f docker-compose.local.yml up --abort-on-container-exit --build
```

## Docs

- `docs/DESIGN.md`: architecture, interfaces, local testing approach
- `docs/RELEASE.md`: Foundry/pipeline configuration and publishing steps (incl. egress policy)
