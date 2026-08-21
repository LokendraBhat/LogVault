package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// ── Session store ─────────────────────────────────────────────────────────────

type session struct{ createdAt time.Time }

var (
	sessions   = map[string]session{}
	sessionsMu sync.Mutex
	// sessionTTL default is 8 hours; overridable via SESSION_TTL_HOURS in main().
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

// sessionTTLLabel renders sessionTTL for display on the login page, e.g. "8 hours" or "1 hour".
func sessionTTLLabel() string {
	if h := sessionTTL.Hours(); h == float64(int64(h)) {
		n := int64(h)
		if n == 1 {
			return "1 hour"
		}
		return fmt.Sprintf("%d hours", n)
	}
	return sessionTTL.String()
}
