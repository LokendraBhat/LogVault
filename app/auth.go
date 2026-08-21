package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// ── Auth config ───────────────────────────────────────────────────────────────

type authConfig struct {
	enabled  bool
	username string
	password string
}

var auth authConfig

// ── Login lockout ─────────────────────────────────────────────────────────────
// Global (not per-IP/user) failed-login counter. Resets on process restart,
// same as sessions — no persistence by design.

var (
	maxLoginAttempts = 5
	loginLockout     = 24 * time.Hour

	loginMu          sync.Mutex
	failedLoginCount int
	lockedUntil      time.Time
)

// loginLocked reports whether logins are currently blocked. A lock that has
// naturally expired is cleared here, starting a fresh attempt count.
func loginLocked() bool {
	loginMu.Lock()
	defer loginMu.Unlock()
	if lockedUntil.IsZero() {
		return false
	}
	if time.Now().After(lockedUntil) {
		failedLoginCount = 0
		lockedUntil = time.Time{}
		return false
	}
	return true
}

func recordFailedLogin() {
	loginMu.Lock()
	defer loginMu.Unlock()
	failedLoginCount++
	if failedLoginCount >= maxLoginAttempts {
		lockedUntil = time.Now().Add(loginLockout)
	}
}

func resetFailedLogins() {
	loginMu.Lock()
	defer loginMu.Unlock()
	failedLoginCount = 0
	lockedUntil = time.Time{}
}

// lockoutMessage reports the configured lockout window, not the live
// remaining time — simpler to read, and doesn't hand an attacker an exact
// countdown to when the lock lifts.
func lockoutMessage() string {
	return fmt.Sprintf("Too many failed login attempts. Contact administrator or Try again after %dh.", int(loginLockout.Hours()))
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

// ── Login / logout handlers ───────────────────────────────────────────────────

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
	data := map[string]string{
		"Error":       "",
		"LoginAction": p("/login"),
		"LogoURL":     logoURL,
		"SessionTTL":  sessionTTLLabel(),
	}
	if loginLocked() {
		data["Error"] = lockoutMessage()
		data["Locked"] = "1"
	}
	loginTmpl.Execute(w, data)
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if loginLocked() {
		w.WriteHeader(http.StatusTooManyRequests)
		loginTmpl.Execute(w, map[string]string{
			"Error":       lockoutMessage(),
			"LoginAction": p("/login"),
			"LogoURL":     logoURL,
			"SessionTTL":  sessionTTLLabel(),
			"Locked":      "1",
		})
		return
	}

	r.ParseForm()
	user := r.FormValue("username")
	pass := r.FormValue("password")

	if user == auth.username && pass == auth.password {
		resetFailedLogins()
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

	recordFailedLogin()
	w.WriteHeader(http.StatusUnauthorized)
	loginTmpl.Execute(w, map[string]string{
		"Error":       "Invalid username or password",
		"LoginAction": p("/login"),
		"LogoURL":     logoURL,
		"SessionTTL":  sessionTTLLabel(),
	})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
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
