################################################################################
# Stage 1: Build the application
FROM golang:1.24-alpine AS BUILD
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-w -s" -o bin/main cmd/main.go

################################################################################
# Stage 2: Runtime container
FROM alpine:latest AS RUNNER
WORKDIR /app

RUN apk add --no-cache ca-certificates ffmpeg 
# https://github.com/yt-dlp/yt-dlp-wiki/blob/master/Installation.md
RUN apk -U add yt-dlp

# Copy the compiled binary from the build stage
COPY --from=BUILD /app/bin/main /app/bin/main
COPY public /app/

# Set environment variables
ENV DL_FOLDER_ROOT=/app/public/downloads

# Create downloads directory
RUN mkdir -p /app/public/downloads

# Command to run the executable
CMD ["/app/bin/main"]
