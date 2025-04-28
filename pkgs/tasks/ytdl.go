package tasks

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/http/httpproxy"

	ytdl "github.com/WangWilly/go-youtube-dl/downloader"
)

////////////////////////////////////////////////////////////////////////////////

type DownloadTask struct {
	taskID    string
	targetUrl string
	filePath  string
	progress  int64
	ctx       context.Context
	cancel    context.CancelFunc
}

////////////////////////////////////////////////////////////////////////////////

func NewTask(url string, filepath string) *DownloadTask {
	ctx, cancel := context.WithCancel(context.Background())
	task := &DownloadTask{
		taskID:    uuid.New().String(),
		targetUrl: url,
		filePath:  filepath,
		progress:  0,
		ctx:       ctx,
		cancel:    cancel,
	}
	return task
}

func NewTaskWithCtx(ctx context.Context, url string, filepath string) *DownloadTask {
	ctx, cancel := context.WithCancel(ctx)
	task := &DownloadTask{
		taskID:    uuid.New().String(),
		targetUrl: url,
		filePath:  filepath,
		progress:  0,
		ctx:       ctx,
		cancel:    cancel,
	}
	return task
}

////////////////////////////////////////////////////////////////////////////////

func (t *DownloadTask) GetID() string {
	return t.taskID
}

func (t *DownloadTask) GetProgress() int64 {
	return t.progress
}

func (t *DownloadTask) Execute() {
	t.progress = 30

	if err := exec.CommandContext(
		t.ctx,
		"yt-dlp",
		"-o", t.filePath,
		"-f", "mp4",
		t.targetUrl,
	).Run(); err != nil {
		t.progress = -1
		if t.ctx.Err() == context.Canceled {
			fmt.Printf("Download canceled: %s\n", t.filePath)
		} else {
			fmt.Printf("Error executing command: %v\n", err)
		}
		return
	}
	t.progress = 100

	fmt.Printf("Download complete: %s\n", t.filePath)
}

func (t *DownloadTask) Cancel() {
	fmt.Printf("Canceling download: %s, ", t.filePath)
	fmt.Printf("Canceling task: %s\n", t.taskID)
	t.cancel()
}

////////////////////////////////////////////////////////////////////////////////
// utils

func GetDownloader(outputDir string) *ytdl.Downloader {
	proxyFunc := httpproxy.FromEnvironment().ProxyFunc()
	httpTransport := &http.Transport{
		// Proxy: http.ProxyFromEnvironment() does not work. Why?
		Proxy: func(r *http.Request) (uri *url.URL, err error) {
			return proxyFunc(r.URL)
		},
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	downloader := &ytdl.Downloader{
		OutputDir: outputDir,
	}
	downloader.HTTPClient = &http.Client{Transport: httpTransport}

	return downloader
}

// // copyWithContext copies from src to dst while respecting context cancellation
// func (t *DownloadTask) copyWithContext(dst io.Writer, src io.Reader) error {
// 	// Buffer size of 32KB
// 	buf := make([]byte, 32*1024)
// 	progressLeft := 100 - t.progress

// 	for {
// 		// Check if context is done before reading
// 		select {
// 		case <-t.ctx.Done():
// 			return context.Canceled
// 		default:
// 		}

// 		nr, readErr := src.Read(buf)
// 		if nr > 0 {
// 			nw, writeErr := dst.Write(buf[0:nr])
// 			if writeErr != nil {
// 				return writeErr
// 			}
// 			if nw != nr {
// 				return io.ErrShortWrite
// 			}
// 			t.progress += int64(nr) * progressLeft / 100
// 			if t.progress >= 100 {
// 				t.progress = 100
// 			}
// 			fmt.Printf("Download progress: %d%%\n", t.progress)
// 		}

// 		if readErr != nil {
// 			if readErr == io.EOF {
// 				return nil
// 			}
// 			return readErr
// 		}
// 	}
// }
