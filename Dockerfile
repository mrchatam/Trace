# syntax=docker/dockerfile:1

# Trace — local-first project graph CLI + MCP (CGO / tree-sitter required).
# Default image: trace CLI. Target trace-mcp for stdio MCP only.

FROM golang:1.24-bookworm AS builder
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc g++ git ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/
COPY api/ api/
COPY scripts/embed-gui.sh scripts/embed-gui.sh
COPY internal/httpapi/embeddist/ internal/httpapi/embeddist/
COPY web/package.json web/package-lock.json web/
COPY web/tsconfig.json web/tsconfig.app.json web/tsconfig.node.json web/vite.config.ts web/index.html web/
COPY web/src/ web/src/
COPY web/public/ web/public/

# Refresh embed if web sources changed; skip npm when embeddist already valid.
RUN test -f internal/httpapi/embeddist/index.html

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 go build -ldflags="-s -w" -o /out/trace ./cmd/trace && \
    CGO_ENABLED=1 go build -ldflags="-s -w" -o /out/trace-mcp ./cmd/trace-mcp

FROM debian:bookworm-slim AS trace
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /out/trace /usr/local/bin/trace
ENTRYPOINT ["trace"]
CMD ["--help"]

FROM debian:bookworm-slim AS trace-mcp
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /out/trace-mcp /usr/local/bin/trace-mcp
ENTRYPOINT ["trace-mcp"]
