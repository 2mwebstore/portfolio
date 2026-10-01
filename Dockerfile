# --- Stage 1: build the Vue frontend ---
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# --- Stage 2: build the Go backend ---
# mattn/go-sqlite3 needs cgo, so we build on a glibc image with gcc.
FROM golang:1.22-bookworm AS backend-builder
WORKDIR /app/backend
ENV CGO_ENABLED=1
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN go build -o server .

# --- Stage 3: runtime ---
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app

COPY --from=backend-builder /app/backend/server ./server
COPY --from=frontend-builder /app/frontend/dist ./web

# Persistent data (sqlite db + uploaded images) — docker-compose.yml mounts ./data here.
ENV DATA_DIR=/app/data
RUN mkdir -p /app/data

EXPOSE 8080
CMD ["./server"]
