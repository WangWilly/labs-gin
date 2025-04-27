################################################################################
# Stage 1: Base image with Go environment
FROM golang:1.24-alpine AS baseimage
WORKDIR /app

################################################################################
# Stage 2: Build the application
FROM baseimage AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bin/main cmd/main.go

################################################################################
# Stage 3: Runtime container
FROM alpine:latest AS runner
WORKDIR /app
RUN apk add --no-cache ca-certificates ffmpeg

# Copy the compiled binary from the build stage
COPY --from=build /app/bin/main /app/bin/main
COPY public /app/
# Create downloads directory
RUN mkdir -p /app/public/downloads

# Set environment variables
ENV DL_FOLDER_ROOT=/app/public/downloads
ENV NUM_WORKERS=4

# Command to run the executable
CMD ["/app/bin/main"]
