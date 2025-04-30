package dltask

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

////////////////////////////////////////////////////////////////////////////////

func (c *Controller) GetFile(ctx *gin.Context) {
	fmt.Println("GetFile called")

	fileID := ctx.Param("fid")
	if fileID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file ID is required"})
		return
	}

	// Prevent path traversal by removing any "../" or similar sequences
	cleanFileID := filepath.Clean(fileID)
	if strings.Contains(cleanFileID, "..") || strings.Contains(cleanFileID, "/") {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid file ID"})
		return
	}

	filePath := filepath.Join(c.cfg.DlFolderRoot, cleanFileID)
	fmt.Println("File path:", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	fileSize := fileStat.Size()
	ctx.Header("Content-Type", "video/mp4")
	ctx.Header("Accept-Ranges", "bytes")

	rangeHeader := ctx.GetHeader("Range")
	if rangeHeader == "" {
		// Serve the entire file
		ctx.Header("Content-Length", fmt.Sprintf("%d", fileSize))
		ctx.Status(http.StatusOK)
		io.Copy(ctx.Writer, file)
		return
	}

	// Parse the Range header
	parts := strings.Split(rangeHeader, "=")
	if len(parts) != 2 || parts[0] != "bytes" {
		ctx.Status(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	byteRange := strings.Split(parts[1], "-")
	if len(byteRange) != 2 {
		ctx.Status(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	start, err := strconv.ParseInt(byteRange[0], 10, 64)
	if err != nil {
		ctx.Status(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	var end int64
	if byteRange[1] != "" {
		end, err = strconv.ParseInt(byteRange[1], 10, 64)
		if err != nil {
			ctx.Status(http.StatusRequestedRangeNotSatisfiable)
			return
		}
	} else {
		end = fileSize - 1
	}

	if start > end || end >= fileSize {
		ctx.Status(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	chunkSize := end - start + 1
	ctx.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	ctx.Header("Content-Length", fmt.Sprintf("%d", chunkSize))
	ctx.Status(http.StatusPartialContent)

	file.Seek(start, io.SeekStart)
	io.CopyN(ctx.Writer, file, chunkSize)
}
