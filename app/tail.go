package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ── Tail feature ──────────────────────────────────────────────────────────────

// tailPageHandler serves the tail UI page.
func tailPageHandler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, p("/tail/"))
	if name == "" {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	fp := filepath.Clean(filepath.Join(logsDir, filepath.FromSlash(name)))
	base := filepath.Clean(logsDir) + string(os.PathSeparator)
	if !strings.HasPrefix(fp, base) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	fi, err := os.Stat(fp)
	if err != nil || fi.IsDir() {
		http.NotFound(w, r)
		return
	}

	streamURL := p("/tail-stream/" + name)
	streamURLJSON, _ := json.Marshal(streamURL)

	data := TailPageData{
		FileName:     filepath.Base(name),
		FilePath:     "/app/logs/" + name,
		StreamURL:    template.JS(streamURLJSON),
		BrowseURL:    browseURLFromFilePath(name),
		AuthEnabled:  auth.enabled,
		LogoutAction: p("/logout"),
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tailTmpl.Execute(w, data)
}

// tailStreamHandler streams new log lines via Server-Sent Events.
// It seeks to the current end of the file on each new connection, so a
// page refresh clears the view and shows only lines written afterwards.
func tailStreamHandler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, p("/tail-stream/"))
	if name == "" {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	fp := filepath.Clean(filepath.Join(logsDir, filepath.FromSlash(name)))
	base := filepath.Clean(logsDir) + string(os.PathSeparator)
	if !strings.HasPrefix(fp, base) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	f, err := os.Open(fp)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer f.Close()

	// Seek to end so we only stream new content written after this connection.
	f.Seek(0, io.SeekEnd)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // disable nginx buffering
	flusher.Flush()

	reader := bufio.NewReader(f)
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for {
				line, err := reader.ReadString('\n')
				if len(line) > 0 {
					text := strings.TrimRight(line, "\r\n")
					// Replace embedded newlines so they don't break SSE framing.
					text = strings.ReplaceAll(text, "\n", " ")
					fmt.Fprintf(w, "data: %s\n\n", text)
					flusher.Flush()
				}
				if err != nil {
					break // no more data yet — wait for next tick
				}
			}
		}
	}
}
