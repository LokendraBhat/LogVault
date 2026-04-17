package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ── Handlers ──────────────────────────────────────────────────────────────────

func browseHandler(w http.ResponseWriter, r *http.Request) {
	subPath := strings.TrimPrefix(r.URL.Path, p("/browse/"))
	subPath = strings.Trim(subPath, "/")

	targetDir := filepath.Clean(filepath.Join(logsDir, filepath.FromSlash(subPath)))
	base := filepath.Clean(logsDir)
	if targetDir != base && !strings.HasPrefix(targetDir, base+string(os.PathSeparator)) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	info, err := os.Stat(targetDir)
	if err != nil || !info.IsDir() {
		http.NotFound(w, r)
		return
	}

	entries, err := listDir(targetDir, subPath)
	if err != nil {
		entries = []EntryView{}
	}

	data := PageData{
		Entries:      entries,
		Count:        len(entries),
		SubPath:      subPath,
		ParentURL:    parentURL(subPath),
		BrowseRoot:   p("/browse/"),
		Crumbs:       buildCrumbs(subPath),
		Port:         getEnv("PORT", "8080"),
		AuthEnabled:  auth.enabled,
		LogoutAction: p("/logout"),
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	browserTmpl.Execute(w, data)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, p("/download/"))
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

	fi, _ := f.Stat()
	if fi.IsDir() {
		http.Error(w, "Cannot download a directory", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(name)))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fi.Size()))
	io.Copy(w, f)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok","time":"%s"}`, time.Now().UTC().Format(time.RFC3339))
}
