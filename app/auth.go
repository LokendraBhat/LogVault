package main

import "net/http"

// ── Auth config ───────────────────────────────────────────────────────────────

type authConfig struct {
	enabled  bool
	username string
	password string
}

var auth authConfig

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
	loginTmpl.Execute(w, map[string]string{
		"Error":       "",
		"LoginAction": p("/login"),
		"LogoURL":     logoURL,
		"SessionTTL":  sessionTTLLabel(),
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
