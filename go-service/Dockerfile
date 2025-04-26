# Stage 1: Build the Go binary
FROM golang:1.24 AS builder
WORKDIR /app

# 1) Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# 2) Copy the rest of your source
COPY . .

# 3) Build the static binary
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -o go-service ./cmd/server/main.go

# Stage 2: Production image
FROM alpine:latest
RUN apk add --no-cache ca-certificates

WORKDIR /root/

# 1) Copy the binary
COPY --from=builder /app/go-service ./

# 2) Copy your mapping.json
COPY --from=builder /app/mapping.json ./

# 3) Copy your .env so godotenv.Load() will find it
COPY --from=builder /app/.env ./.env

# 3) Copy your Firebase service account key
COPY Service-Account-Cred-firebase.json ./
# 4) Tell the Firebase SDK where to find it
ENV GOOGLE_APPLICATION_CREDENTIALS="/root/Service-Account-Cred-firebase.json"

EXPOSE 8080

CMD ["./go-service"]
