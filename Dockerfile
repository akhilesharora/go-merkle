# Stage 1: Build the application
FROM golang:1.21.4 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go-merkle ./cmd/main.go

# Stage 2: Build a small image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/go-merkle .
CMD ["./go-merkle"]
