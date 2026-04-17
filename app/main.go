package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

// ── Constants & globals ───────────────────────────────────────────────────────

const logsDir = "/app/logs"

// basePath is the URL prefix (e.g. "/logvault"). Always no trailing slash.
var basePath string

// getEnv returns the environment variable value or fallback if unset/empty.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// p returns an absolute URL path with the basePath prefix applied.
func p(path string) string {
	if basePath == "" {
		return path
	}
	return basePath + path
}

// ── Entry point ───────────────────────────────────────────────────────────────

func main() {
	port := getEnv("PORT", "8080")

	// Normalise BASE_PATH: always no trailing slash, must start with / if set.
	bp := os.Getenv("BASE_PATH")
	bp = strings.TrimRight(bp, "/")
	if bp != "" && !strings.HasPrefix(bp, "/") {
		bp = "/" + bp
	}
	basePath = bp
	if basePath != "" {
		fmt.Printf("Base path: %s\n", basePath)
	}

	user := os.Getenv("AUTH_USER")
	pass := os.Getenv("AUTH_PASSWORD")
	if user != "" && pass != "" {
		auth = authConfig{enabled: true, username: user, password: pass}
		fmt.Println("Auth enabled — login page active")
	} else {
		fmt.Println("Auth disabled — set AUTH_USER and AUTH_PASSWORD to enable")
	}

	os.MkdirAll(logsDir, 0755)

	mux := http.NewServeMux()

	// Root redirect.
	rootPath := basePath + "/"
	if basePath == "" {
		rootPath = "/"
	}
	mux.HandleFunc(rootPath, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != rootPath {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, p("/browse/"), http.StatusFound)
	})

	mux.HandleFunc(p("/login"), func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			loginPostHandler(w, r)
		} else {
			loginPageHandler(w, r)
		}
	})
	mux.HandleFunc(p("/logout"), logoutHandler)
	mux.HandleFunc(p("/browse/"), requireAuth(browseHandler))
	mux.HandleFunc(p("/download/"), requireAuth(downloadHandler))
	mux.HandleFunc(p("/tail/"), requireAuth(tailPageHandler))
	mux.HandleFunc(p("/tail-stream/"), requireAuth(tailStreamHandler))
	mux.HandleFunc(p("/view/"), requireAuth(viewHandler))
	mux.HandleFunc(p("/health"), healthHandler)

	fmt.Printf("LogVault running on :%s — serving %s\n", port, logsDir)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
