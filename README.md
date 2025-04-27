# Labs-Gin

A Go-based API service for video downloading and processing using the Gin web framework.

## Description

This project provides a RESTful API for downloading videos from various platforms. It uses a task-based system to manage downloads asynchronously and allows for progress tracking, cancellation, and streaming of downloaded content.

## Installation

### Prerequisites

- Go 1.24 or higher
- [GVM](https://github.com/moovweb/gvm) (optional, for managing Go versions)

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

## API Documentation

### Download Tasks

#### Create a download task
- **URL**: `/dlTask`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "url": "https://www.youtube.com/watch?v=example"
  }
  ```
- **Success Response**:
  - **Code**: 201 Created
  - **Content**:
    ```json
    {
      "task_id": "550e8400-e29b-41d4-a716-446655440000",
      "file_id": "550e8400-e29b-41d4-a716-446655440000.mp4",
      "status": "task submitted"
    }
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
      "task_id": "550e8400-e29b-41d4-a716-446655440000",
      "status": 75
    }
    ```
  - `status` is an integer representing the download progress (0-100)

#### Cancel a task
- **URL**: `/dlTask/:tid`
- **Method**: `DELETE`
- **URL Parameters**: `tid` - The task ID
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "task_id": "550e8400-e29b-41d4-a716-446655440000",
      "status_before_cancel": 45,
      "status": "task cancelled"
    }
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

## Environment Variables

| Name | Description | Default |
|------|-------------|---------|
| DL_FOLDER_ROOT | Directory for downloaded files | ./public/downloads |
| NUM_WORKERS | Number of concurrent download workers | 4 |

## Development Resources

- [Go Modules Documentation](https://go.dev/wiki/Modules#quick-start)
- [YouTube Downloader Library](https://github.com/kkdai/youtube)
