# syntax=docker/dockerfile:1

FROM --platform=linux/amd64 golang:1.25 AS builder

WORKDIR /src
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/enricher ./cmd/enricher

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /out/enricher /enricher

ENTRYPOINT ["/enricher"]
CMD ["foundry"]
