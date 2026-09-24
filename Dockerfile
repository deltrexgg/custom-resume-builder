# =========================
# Build stage
# =========================
FROM golang:1.26.1-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o resume-builder .


# =========================
# Runtime stage
# =========================
FROM debian:bookworm-slim

WORKDIR /app

# Install LibreOffice + certificates
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        libreoffice \
        ca-certificates \
        fonts-dejavu \
        fonts-liberation \
    && rm -rf /var/lib/apt/lists/*

# Copy application
COPY --from=builder /app/resume-builder .

# Runtime directories used by the application
COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/prompt ./prompt
COPY --from=builder /app/docx ./docx
COPY --from=builder /app/pdf ./pdf

EXPOSE 8090

CMD ["./resume-builder"]