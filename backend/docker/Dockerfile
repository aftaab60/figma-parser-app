FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy the source code from the mounted volume to the container's WORKDIR
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/backend-service .

# Create the final, small production image
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/bin/backend-service .
EXPOSE 3000
CMD ["./backend-service"]
