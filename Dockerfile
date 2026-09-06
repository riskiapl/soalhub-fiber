# Stage 1: Build binary Go
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code backend
COPY . .

# Build binary statis yang efisien
RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd/api

# Stage 2: Production image ringan (Alpine)
FROM alpine:latest

WORKDIR /root/

# Install ca-certificates untuk kebutuhan SSL/TLS
RUN apk --no-cache add ca-certificates

# Copy binary hasil build dari Stage 1
COPY --from=builder /main .
COPY .env .env

EXPOSE 3000

CMD ["./main"]