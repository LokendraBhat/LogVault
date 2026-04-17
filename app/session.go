package main

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

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
