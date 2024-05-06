# Build Stage
FROM golang:1.21 AS builder 
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Run Stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]

