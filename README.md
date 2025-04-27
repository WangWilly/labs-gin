# Labs-Gin

A Go-based API service for video downloading and processing using the Gin web framework.

## Description

This project provides a RESTful API for downloading videos from various platforms. It uses a task-based system to manage downloads asynchronously and allows for progress tracking, cancellation, and streaming of downloaded content.

## Installation

### Prerequisites

- Go 1.24 or higher
- [GVM](https://github.com/moovweb/gvm) (optional, for managing Go versions)
- https://github.com/yt-dlp/yt-dlp-wiki/blob/master/Installation.md

### Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/WangWilly/labs-gin.git
   cd labs-gin
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the development server:
   ```bash
   ./scripts/dev.sh
   ```

### Docker Setup

1. Build the Docker image using the provided build script:
   ```bash
   ./scripts/build.sh
   ```
   This will create a Docker image named `labs-gin-app:latest`.

2. Run the service using Docker Compose:
   ```bash
   cd deployments
   docker compose up -d
   ```
   This will start the service in detached mode, listening on port 8080.

3. To stop the service:
   ```bash
   docker compose down
   ```

4. Monitor logs:
   ```bash
   docker compose logs -f
   ```

### Docker Environment Variables

When running with Docker Compose, you can configure the following environment variables in the `deployments/docker-compose.yml` file:

| Name | Description | Default |
|------|-------------|---------|
| PORT | The port on which the service listens | 8080 |
| TASK_MENAGER_NUM_WORKERS | Number of concurrent download workers | 4 |

The downloaded files will be persisted in the `./public/downloads` directory on your host machine through Docker volume mapping.

## API Documentation

### Download Tasks

#### Create a Download Task
- **Endpoint**: `/dlTask`
- **HTTP Method**: `POST`
- **Description**: Initiates a new video download task for the specified URL.
- **Request Parameters**:
  | Parameter | Type | Required | Description |
  |-----------|------|----------|-------------|
  | url | string | Yes | The video URL to download |

- **Request Body Example**:
  ```json
  {
    "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
  }
  ```

- **Response**: 
  - **Status Code**: 201 Created
  - **Content Type**: application/json
  - **Body**:
    ```json
    {
      "task_id": "b3a63526-24c0-4fe8-a068-f8ae28349788",
      "file_id": "550e8400-e29b-41d4-a716-446655440000.mp4",
      "status": "task submitted"
    }
    ```
  - **Response Fields**:
    | Field | Type | Description |
    |-------|------|-------------|
    | task_id | string | Unique identifier for tracking and managing the download task |
    | file_id | string | Identifier of the resulting video file (used for accessing the file) |
    | status | string | Current status of the task |

- **Example Request**:
  ```bash
  curl -X POST http://localhost:8080/dlTask \
    -H "Content-Type: application/json" \
    -d '{"url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"}'
  ```

#### Get task status
- **URL**: `/dlTask/:tid`
- **Method**: `GET`
- **URL Parameters**: `tid` - The task ID
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "task_id": "b3a63526-24c0-4fe8-a068-f8ae28349788",
      "status": 75
    }
    ```
  - `status` is an integer representing the download progress (0-100)

- **curl Example**:
  ```bash
  curl -X GET http://localhost:8080/dlTask/b3a63526-24c0-4fe8-a068-f8ae28349788
  ```

#### Cancel a task
- **URL**: `/dlTask/:tid`
- **Method**: `DELETE`
- **URL Parameters**: `tid` - The task ID
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "task_id": "b3a63526-24c0-4fe8-a068-f8ae28349788",
      "status_before_cancel": 45,
      "status": "task cancelled"
    }
    ```

- **curl Example**:
  ```bash
  curl -X DELETE http://localhost:8080/dlTask/b3a63526-24c0-4fe8-a068-f8ae28349788
  ```

### File Access

#### Stream or download a file
- **URL**: `/dlTaskFile/:fid`
- **Method**: `GET`
- **URL Parameters**: `fid` - The file ID
- **Success Response**:
  - **Code**: 200 OK or 206 Partial Content
  - **Content**: The requested video file
  - **Headers**:
    - `Content-Type`: video/mp4
    - `Accept-Ranges`: bytes
    - `Content-Length`: [file size]

- **curl Examples**:
  ```bash
  # Download the entire file
  curl -X GET http://localhost:8080/dlTaskFile/550e8400-e29b-41d4-a716-446655440000.mp4 --output video.mp4
  
  # Stream a portion of the file (partial content)
  curl -X GET http://localhost:8080/dlTaskFile/550e8400-e29b-41d4-a716-446655440000.mp4 \
    -H "Range: bytes=0-1048576" --output video_part.mp4
  ```

## Environment Variables

| Name | Description | Default |
|------|-------------|---------|
| DL_FOLDER_ROOT | Directory for downloaded files | ./public/downloads |
| TASK_MENAGER_NUM_WORKERS | Number of concurrent download workers | 4 |

## Development Resources

- [Go Modules Documentation](https://go.dev/wiki/Modules#quick-start)
- [YouTube Downloader Library](https://github.com/kkdai/youtube)
- https://github.com/yt-dlp/yt-dlp-wiki/blob/master/Installation.md
