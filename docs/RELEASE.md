# Release (Foundry)

This repo is public. Do not put real Foundry URLs, dataset RIDs, tokens, or API keys in git.

## Preconditions

- A pipeline-mode Compute Module is created in Foundry
- Input dataset and output dataset exist (and their schemas are decided)
- You know the alias names you will configure in the pipeline (e.g. `input`, `output`)

## Required Configuration

### Environment Variables

Foundry injects:

- `BUILD2_TOKEN` (file path)
- `RESOURCE_ALIAS_MAP` (file path)

You must provide (via compute module configuration):

- `FOUNDRY_URL`: Foundry base URL (example format: `https://<your-stack>.palantirfoundry.com`)
- `GEMINI_API_KEY`: inject as a secret (do not hardcode)
- `GEMINI_MODEL`: Gemini model name (do not hardcode in code; configure per environment)

Optional knobs:

- `WORKERS` (int)
- `MAX_RETRIES` (int)
- `REQUEST_TIMEOUT` (duration)
- `FAIL_FAST` (bool)
- `RATE_LIMIT_RPS` (float)
- `GEMINI_CAPTURE_AUDIT` (bool)
- `GEMINI_BASE_URL` (string; for hermetic testing only, not recommended in Foundry)

### Egress Policy

The container needs outbound network access to the Gemini API endpoints (and any related Google endpoints required by the SDK/tooling).

At minimum, expect to allowlist domains like:

- `generativelanguage.googleapis.com` (Gemini API)

Confirm exact domains from the client library / runtime behavior before locking the policy.

## Image Publishing

There are typically two publishing targets:

1. GitHub artifacts (CI visibility and reproducibility)
2. Your org's Palantir registry / artifact store used by Foundry compute modules

Define and document the publishing mechanism used in your environment (exact commands and auth differ by org).

Recommended end state:

- GitHub Actions builds the Docker image and produces a versioned artifact
- A gated workflow publishes the same image tag to the Palantir registry consumed by Foundry

## Rollout Checklist

- First run against a tiny fixture (5-10 emails)
- Validate output schema matches the dataset schema exactly
- Validate logs contain no tokens/API keys and avoid raw emails
- Validate failure policy (fail-fast vs partial output) matches pipeline expectations
- Confirm egress policy is minimal and sufficient

## Rollback

- Revert compute module image tag to the previously known-good tag
- Re-run pipeline build
