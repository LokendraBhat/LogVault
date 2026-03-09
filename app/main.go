package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const logsDir = "/app/logs"

// basePath is the URL prefix (e.g. "/logvault"). Always no trailing slash.
var basePath string

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

// ── Session store ─────────────────────────────────────────────────────────────

type session struct{ createdAt time.Time }

var (
	sessions   = map[string]session{}
	sessionsMu sync.Mutex
	sessionTTL = 8 * time.Hour
)

func newSession() string {
	b := make([]byte, 16)
	rand.Read(b)
	token := hex.EncodeToString(b)
	sessionsMu.Lock()
	sessions[token] = session{createdAt: time.Now()}
	sessionsMu.Unlock()
	return token
}

func validSession(token string) bool {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	s, ok := sessions[token]
	if !ok {
		return false
	}
	if time.Since(s.createdAt) > sessionTTL {
		delete(sessions, token)
		return false
	}
	return true
}

func deleteSession(token string) {
	sessionsMu.Lock()
	delete(sessions, token)
	sessionsMu.Unlock()
}

// ── Auth config ───────────────────────────────────────────────────────────────

type authConfig struct {
	enabled  bool
	username string
	password string
}

var auth authConfig

// ── Entry type ────────────────────────────────────────────────────────────────

type Entry struct {
	Name     string
	IsDir    bool
	Size     string
	Modified string
	RelPath  string
}

// ── Templates ─────────────────────────────────────────────────────────────────

var loginTmpl = template.Must(template.New("login").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>LogVault · Sign In</title>
  <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;600;700&family=Syne:wght@400;700;800&display=swap" rel="stylesheet">
  <style>
    :root{--bg:#0a0a0f;--surface:#111118;--border:#1e1e2e;--accent:#00ff9d;--accent2:#7c3aed;--text:#e2e8f0;--muted:#4a5568}
    *{margin:0;padding:0;box-sizing:border-box}
    body{background:var(--bg);color:var(--text);font-family:'JetBrains Mono',monospace;min-height:100vh;display:flex;align-items:center;justify-content:center}
    body::before{content:'';position:fixed;inset:0;background-image:linear-gradient(rgba(0,255,157,.03) 1px,transparent 1px),linear-gradient(90deg,rgba(0,255,157,.03) 1px,transparent 1px);background-size:40px 40px;pointer-events:none}
    .card{position:relative;z-index:1;width:100%;max-width:400px;margin:1rem;background:var(--surface);border:1px solid var(--border);border-radius:14px;padding:2.5rem;box-shadow:0 0 60px rgba(0,0,0,.6),0 0 0 1px rgba(0,255,157,.05);animation:rise .5s ease both}
    @keyframes rise{from{opacity:0;transform:translateY(20px)}to{opacity:1;transform:translateY(0)}}
    .logo-row{display:flex;align-items:center;gap:.85rem;margin-bottom:.4rem}
    .logo-icon{width:40px;height:40px;background:linear-gradient(135deg,var(--accent),var(--accent2));border-radius:8px;display:flex;align-items:center;justify-content:center;font-size:1.1rem;box-shadow:0 0 18px rgba(0,255,157,.25)}
    h1{font-family:'Syne',sans-serif;font-weight:800;font-size:1.7rem;background:linear-gradient(135deg,#fff 30%,var(--accent));-webkit-background-clip:text;-webkit-text-fill-color:transparent}
    .tagline{font-size:.72rem;color:var(--muted);margin-bottom:2rem;letter-spacing:.05em}
    label{display:block;font-size:.7rem;color:var(--muted);letter-spacing:.1em;text-transform:uppercase;margin-bottom:.4rem}
    input{width:100%;padding:.7rem .9rem;background:rgba(255,255,255,.04);border:1px solid var(--border);border-radius:6px;color:var(--text);font-family:'JetBrains Mono',monospace;font-size:.85rem;outline:none;transition:border-color .2s;margin-bottom:1.1rem}
    input:focus{border-color:rgba(0,255,157,.4);box-shadow:0 0 0 3px rgba(0,255,157,.06)}
    button{width:100%;padding:.8rem;background:linear-gradient(135deg,rgba(0,255,157,.15),rgba(124,58,237,.15));border:1px solid rgba(0,255,157,.35);border-radius:7px;color:var(--accent);font-family:'JetBrains Mono',monospace;font-size:.85rem;font-weight:600;cursor:pointer;letter-spacing:.05em;transition:all .2s;margin-top:.4rem}
    button:hover{background:rgba(0,255,157,.12);box-shadow:0 0 18px rgba(0,255,157,.18)}
    .error{background:rgba(255,71,87,.1);border:1px solid rgba(255,71,87,.3);color:#ff6b7a;border-radius:6px;padding:.65rem .9rem;font-size:.78rem;margin-bottom:1.2rem}
    .divider{height:1px;background:var(--border);margin:1.5rem 0}
    .hint{font-size:.68rem;color:var(--muted);text-align:center}
  </style>
</head>
<body>
  <div class="card">
    <div class="logo-row"><div class="logo-icon">⬡</div><h1>LogVault</h1></div>
    <p class="tagline">Secure log access — sign in to continue</p>
    {{if .Error}}<div class="error">⚠ {{.Error}}</div>{{end}}
    <form method="POST" action="{{.LoginAction}}">
      <label>Username</label>
      <input type="text" name="username" autocomplete="username" autofocus placeholder="enter username">
      <label>Password</label>
      <input type="password" name="password" autocomplete="current-password" placeholder="••••••••">
      <button type="submit">→ Sign In</button>
    </form>
    <div class="divider"></div>
    <p class="hint">Session expires after 8 hours</p>
  </div>
</body>
</html>`))

var funcMap = template.FuncMap{"notDir": func(b bool) bool { return !b }}

var browserTmpl = template.Must(template.New("browser").Funcs(funcMap).Parse(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>LogVault{{if .SubPath}} · /{{.SubPath}}{{end}}</title>
  <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;600;700&family=Syne:wght@400;700;800&display=swap" rel="stylesheet">
  <style>
    :root{--bg:#0a0a0f;--surface:#111118;--border:#1e1e2e;--accent:#00ff9d;--accent2:#7c3aed;--text:#e2e8f0;--muted:#4a5568;--folder:#f59e0b}
    *{margin:0;padding:0;box-sizing:border-box}
    body{background:var(--bg);color:var(--text);font-family:'JetBrains Mono',monospace;min-height:100vh;overflow-x:hidden}
    body::before{content:'';position:fixed;inset:0;background-image:linear-gradient(rgba(0,255,157,.03) 1px,transparent 1px),linear-gradient(90deg,rgba(0,255,157,.03) 1px,transparent 1px);background-size:40px 40px;pointer-events:none;z-index:0}
    .container{max-width:980px;margin:0 auto;padding:2.5rem 2rem;position:relative;z-index:1}
    .topbar{display:flex;align-items:center;justify-content:space-between;margin-bottom:2rem;animation:slideIn .5s ease both}
    .logo-row{display:flex;align-items:center;gap:.85rem}
    .logo-icon{width:38px;height:38px;background:linear-gradient(135deg,var(--accent),var(--accent2));border-radius:8px;display:flex;align-items:center;justify-content:center;font-size:1rem;box-shadow:0 0 16px rgba(0,255,157,.25)}
    h1{font-family:'Syne',sans-serif;font-weight:800;font-size:1.8rem;background:linear-gradient(135deg,#fff 30%,var(--accent));-webkit-background-clip:text;-webkit-text-fill-color:transparent}
    .logout-btn{padding:.4rem .9rem;background:transparent;border:1px solid rgba(255,71,87,.35);border-radius:5px;color:#ff6b7a;font-family:'JetBrains Mono',monospace;font-size:.72rem;cursor:pointer;transition:all .2s;letter-spacing:.05em}
    .logout-btn:hover{background:rgba(255,71,87,.08);border-color:#ff4757}
    .breadcrumb{display:flex;align-items:center;gap:.4rem;flex-wrap:wrap;padding:.6rem 1rem;background:var(--surface);border:1px solid var(--border);border-radius:6px;font-size:.75rem;margin-bottom:1.25rem;animation:slideIn .5s .05s ease both}
    .breadcrumb a{color:var(--accent);text-decoration:none}
    .breadcrumb a:hover{text-decoration:underline}
    .breadcrumb .sep{color:var(--muted)}
    .breadcrumb .current{color:var(--text)}
    .status-bar{display:flex;align-items:center;gap:.75rem;margin-bottom:1.5rem;padding:.65rem 1rem;background:var(--surface);border:1px solid var(--border);border-radius:6px;font-size:.72rem;color:var(--muted);animation:slideIn .5s .1s ease both}
    .dot{width:6px;height:6px;border-radius:50%;background:var(--accent);box-shadow:0 0 8px var(--accent);animation:pulse 2s infinite}
    @keyframes pulse{0%,100%{opacity:1}50%{opacity:.4}}
    .status-bar span{color:var(--text)}
    .tag{display:inline-block;padding:.15rem .45rem;background:rgba(124,58,237,.15);border:1px solid rgba(124,58,237,.3);border-radius:3px;color:#a78bfa;font-size:.65rem;letter-spacing:.08em}
    .log-table{background:var(--surface);border:1px solid var(--border);border-radius:10px;overflow:hidden;animation:slideIn .5s .15s ease both}
    .table-header{display:grid;grid-template-columns:1fr 110px 160px 120px;padding:.7rem 1.25rem;background:rgba(255,255,255,.02);border-bottom:1px solid var(--border);font-size:.62rem;text-transform:uppercase;letter-spacing:.12em;color:var(--muted)}
    .log-row{display:grid;grid-template-columns:1fr 110px 160px 120px;padding:.9rem 1.25rem;border-bottom:1px solid var(--border);align-items:center;transition:background .15s;position:relative}
    .log-row:last-child{border-bottom:none}
    .log-row:hover{background:rgba(0,255,157,.025)}
    .log-row::before{content:'';position:absolute;left:0;top:0;bottom:0;width:2px;background:var(--accent);opacity:0;transition:opacity .15s}
    .log-row:hover::before{opacity:1}
    .log-row.is-dir::before{background:var(--folder)}
    .entry-name{display:flex;align-items:center;gap:.6rem;font-size:.83rem;font-weight:600;overflow:hidden}
    .entry-icon{flex-shrink:0;width:28px;height:28px;border-radius:5px;display:flex;align-items:center;justify-content:center;font-size:.85rem}
    .file-icon{background:rgba(0,255,157,.07);border:1px solid rgba(0,255,157,.18)}
    .dir-icon{background:rgba(245,158,11,.1);border:1px solid rgba(245,158,11,.25)}
    .entry-link{white-space:nowrap;overflow:hidden;text-overflow:ellipsis;text-decoration:none;color:inherit;transition:color .15s}
    .entry-link:hover{color:var(--accent)}
    .dir-link{color:#fcd34d}
    .dir-link:hover{color:#fbbf24}
    .up-link{color:#a78bfa}
    .file-size{font-size:.75rem;color:var(--muted)}
    .file-modified{font-size:.72rem;color:var(--muted)}
    .download-btn{display:inline-flex;align-items:center;gap:.35rem;padding:.35rem .8rem;background:transparent;border:1px solid rgba(0,255,157,.35);color:var(--accent);border-radius:5px;font-family:'JetBrains Mono',monospace;font-size:.7rem;text-decoration:none;transition:all .15s;letter-spacing:.04em}
    .download-btn:hover{background:rgba(0,255,157,.08);box-shadow:0 0 10px rgba(0,255,157,.15)}
    .empty-state{text-align:center;padding:5rem 2rem;color:var(--muted)}
    .empty-icon{font-size:3rem;margin-bottom:1rem;opacity:.3}
    .empty-state h3{font-family:'Syne',sans-serif;font-size:1.1rem;color:var(--text);margin-bottom:.5rem}
    .empty-state p{font-size:.78rem}
    .footer{margin-top:1.8rem;text-align:center;font-size:.68rem;color:var(--muted)}
    @keyframes slideIn{from{opacity:0;transform:translateY(10px)}to{opacity:1;transform:translateY(0)}}
    @media(max-width:600px){.table-header,.log-row{grid-template-columns:1fr 90px 110px}.table-header>*:last-child,.log-row>*:last-child{display:none}}
  </style>
</head>
<body>
  <div class="container">
    <div class="topbar">
      <div class="logo-row"><div class="logo-icon">⬡</div><h1>LogVault</h1></div>
      {{if .AuthEnabled}}
      <form method="POST" action="{{.LogoutAction}}" style="margin:0">
        <button class="logout-btn" type="submit">⏏ Sign out</button>
      </form>
      {{end}}
    </div>

    <div class="breadcrumb">
      <a href="{{.BrowseRoot}}">⌂ root</a>
      {{range .Crumbs}}
        <span class="sep">/</span>
        {{if .IsLast}}<span class="current">{{.Name}}</span>
        {{else}}<a href="{{.URL}}">{{.Name}}</a>{{end}}
      {{end}}
    </div>

    <div class="status-bar">
      <div class="dot"></div>
      <span>{{.Count}} item{{if ne .Count 1}}s{{end}}</span>
      &nbsp;·&nbsp;<span class="tag">/app/logs{{if .SubPath}}/{{.SubPath}}{{end}}</span>
    </div>

    {{if .Entries}}
    <div class="log-table">
      <div class="table-header">
        <div>Name</div><div>Size</div><div>Modified</div><div>Action</div>
      </div>
      {{if .SubPath}}
      <div class="log-row">
        <div class="entry-name">
          <div class="entry-icon dir-icon">↑</div>
          <a class="entry-link up-link" href="{{.ParentURL}}">..</a>
        </div>
        <div></div><div></div><div></div>
      </div>
      {{end}}
      {{range .Entries}}
      <div class="log-row{{if .IsDir}} is-dir{{end}}">
        <div class="entry-name">
          <div class="entry-icon {{if .IsDir}}dir-icon{{else}}file-icon{{end}}">{{if .IsDir}}📁{{else}}📄{{end}}</div>
          {{if .IsDir}}
            <a class="entry-link dir-link" href="{{.BrowseURL}}">{{.Name}}</a>
          {{else}}
            <span class="entry-link" style="cursor:default">{{.Name}}</span>
          {{end}}
        </div>
        <div class="file-size">{{if notDir .IsDir}}{{.Size}}{{else}}—{{end}}</div>
        <div class="file-modified">{{.Modified}}</div>
        <div>
          {{if notDir .IsDir}}
          <a class="download-btn" href="{{.DownloadURL}}">↓ Download</a>
          {{end}}
        </div>
      </div>
      {{end}}
    </div>
    {{else}}
    <div class="empty-state">
      <div class="empty-icon">📭</div>
      <h3>Empty directory</h3>
      <p>No files or folders here.</p>
    </div>
    {{end}}

    <div class="footer">LogVault · Minimal Log Server · Port {{.Port}}</div>
  </div>
</body>
</html>`))

// ── Data types ────────────────────────────────────────────────────────────────

type Crumb struct {
	Name   string
	URL    string
	IsLast bool
}

type EntryView struct {
	Name        string
	IsDir       bool
	Size        string
	Modified    string
	BrowseURL   string
	DownloadURL string
}

type PageData struct {
	Entries      []EntryView
	Count        int
	SubPath      string
	ParentURL    string
	BrowseRoot   string
	Crumbs       []Crumb
	Port         string
	AuthEnabled  bool
	LogoutAction string
}

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

func parentURL(subPath string) string {
	idx := strings.LastIndex(subPath, "/")
	if idx < 0 {
		return p("/browse/")
	}
	return p("/browse/" + subPath[:idx])
}

// ── Auth middleware ───────────────────────────────────────────────────────────

func sessionToken(r *http.Request) string {
	c, err := r.Cookie("lv_session")
	if err != nil {
		return ""
	}
	return c.Value
}

func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !auth.enabled {
			next(w, r)
			return
		}
		if !validSession(sessionToken(r)) {
			http.Redirect(w, r, p("/login"), http.StatusFound)
			return
		}
		next(w, r)
	}
}

// ── Handlers ──────────────────────────────────────────────────────────────────

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	if !auth.enabled {
		http.Redirect(w, r, p("/browse/"), http.StatusFound)
		return
	}
	if validSession(sessionToken(r)) {
		http.Redirect(w, r, p("/browse/"), http.StatusFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	loginTmpl.Execute(w, map[string]string{
		"Error":       "",
		"LoginAction": p("/login"),
	})
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user := r.FormValue("username")
	pass := r.FormValue("password")

	if user == auth.username && pass == auth.password {
		token := newSession()
		cookiePath := "/"
		if basePath != "" {
			cookiePath = basePath + "/"
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "lv_session",
			Value:    token,
			Path:     cookiePath,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   int(sessionTTL.Seconds()),
		})
		http.Redirect(w, r, p("/browse/"), http.StatusFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	loginTmpl.Execute(w, map[string]string{
		"Error":       "Invalid username or password",
		"LoginAction": p("/login"),
	})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	deleteSession(sessionToken(r))
	cookiePath := "/"
	if basePath != "" {
		cookiePath = basePath + "/"
	}
	http.SetCookie(w, &http.Cookie{
		Name: "lv_session", Value: "", Path: cookiePath,
		HttpOnly: true, MaxAge: -1,
	})
	http.Redirect(w, r, p("/login"), http.StatusFound)
}

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

// ── Main ──────────────────────────────────────────────────────────────────────

func main() {
	port := getEnv("PORT", "8080")

	// Normalise BASE_PATH: always no trailing slash, must start with / if set
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

	// Root redirect
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
	mux.HandleFunc(p("/logout"), requireAuth(logoutHandler))
	mux.HandleFunc(p("/browse/"), requireAuth(browseHandler))
	mux.HandleFunc(p("/download/"), requireAuth(downloadHandler))
	mux.HandleFunc(p("/health"), healthHandler)

	fmt.Printf("LogVault running on :%s — serving %s\n", port, logsDir)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
