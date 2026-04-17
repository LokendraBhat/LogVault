package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// ── Helpers ───────────────────────────────────────────────────────────────────

func formatSize(bytes int64) string {
	switch {
	case bytes < 1024:
		return fmt.Sprintf("%d B", bytes)
	case bytes < 1024*1024:
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	case bytes < 1024*1024*1024:
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
	default:
		return fmt.Sprintf("%.2f GB", float64(bytes)/(1024*1024*1024))
	}
}

// listDir reads dir and returns a sorted slice of EntryView.
// relBase is the path relative to logsDir (used to build URLs).
func listDir(dir string, relBase string) ([]EntryView, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var result []EntryView
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		rel := e.Name()
		if relBase != "" {
			rel = relBase + "/" + e.Name()
		}
		ev := EntryView{
			Name:     e.Name(),
			IsDir:    e.IsDir(),
			Modified: info.ModTime().Format("2006-01-02 15:04"),
		}
		if e.IsDir() {
			ev.BrowseURL = p("/browse/" + rel)
		} else {
			ev.Size = formatSize(info.Size())
			ev.DownloadURL = p("/download/" + rel)
			ev.TailURL = p("/tail/" + rel)
			ev.ViewURL = p("/view/" + rel)
		}
		result = append(result, ev)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].IsDir != result[j].IsDir {
			return result[i].IsDir
		}
		return strings.ToLower(result[i].Name) < strings.ToLower(result[j].Name)
	})
	return result, nil
}

// buildCrumbs builds the breadcrumb trail for a given subPath.
func buildCrumbs(subPath string) []Crumb {
	if subPath == "" {
		return nil
	}
	parts := strings.Split(subPath, "/")
	crumbs := make([]Crumb, len(parts))
	for i, part := range parts {
		crumbs[i] = Crumb{
			Name:   part,
			URL:    p("/browse/" + strings.Join(parts[:i+1], "/")),
			IsLast: i == len(parts)-1,
		}
	}
	return crumbs
}

// parentURL returns the browse URL for the parent of subPath.
func parentURL(subPath string) string {
	idx := strings.LastIndex(subPath, "/")
	if idx < 0 {
		return p("/browse/")
	}
	return p("/browse/" + subPath[:idx])
}

// browseURLFromFilePath returns the browse URL for the directory containing
// the given file path (relative to logsDir).
func browseURLFromFilePath(name string) string {
	idx := strings.LastIndex(name, "/")
	if idx < 0 {
		return p("/browse/")
	}
	return p("/browse/" + name[:idx])
}
