package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// ── Constants & globals ───────────────────────────────────────────────────────

const logsDir = "/app/logs"

// basePath is the URL prefix (e.g. "/logvault"). Always no trailing slash.
var basePath string

// logoURL, if set, replaces the default logo icon everywhere it appears.
// It is either the LOGO_URL env var verbatim, or — when LOGO_PATH is set — the
// /logo route that serves that mounted file.
var logoURL string

// logoPath is the mounted logo file served at /logo, if LOGO_PATH is set and valid.
var logoPath string

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
	// Self-check mode: exec-form healthcheck for the shell-less scratch image.
	// `docker exec` can't run `wget`/`sh` there, so the binary checks itself.
	if len(os.Args) > 1 && os.Args[1] == "-healthcheck" {
		runHealthCheck()
		return
	}

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

	logoURL = os.Getenv("LOGO_URL")
	if lp := os.Getenv("LOGO_PATH"); lp != "" {
		if fi, err := os.Stat(lp); err != nil || fi.IsDir() {
			fmt.Printf("LOGO_PATH %q not found — ignoring\n", lp)
		} else {
			logoPath = lp
			logoURL = p("/logo")
		}
	}

	if ttlHours := os.Getenv("SESSION_TTL_HOURS"); ttlHours != "" {
		if h, err := strconv.Atoi(ttlHours); err == nil && h > 0 {
			sessionTTL = time.Duration(h) * time.Hour
		} else {
			fmt.Printf("Invalid SESSION_TTL_HOURS %q — using default (%s)\n", ttlHours, sessionTTL)
		}
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
	mux.HandleFunc(p("/logo"), logoHandler)

	fmt.Printf("LogVault running on :%s — serving %s\n", port, logsDir)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

// runHealthCheck GETs its own /health endpoint and exits 0/1 accordingly.
// Invoked as `/logvault -healthcheck` from HEALTHCHECK, since the scratch
// image has no shell or wget for the usual CMD-SHELL form.
func runHealthCheck() {
	port := getEnv("PORT", "8080")
	bp := strings.TrimRight(os.Getenv("BASE_PATH"), "/")
	if bp != "" && !strings.HasPrefix(bp, "/") {
		bp = "/" + bp
	}
	url := fmt.Sprintf("http://127.0.0.1:%s%s/health", port, bp)

	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		os.Exit(1)
	}
}
