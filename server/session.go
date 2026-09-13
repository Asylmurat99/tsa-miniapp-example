package main

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

// session is what the mini app backend knows about the visitor. The key is
// called partnerUserID on purpose: it is the alias the platform issues to
// this mini app, not the subscriber's identifier inside Telecom.
type session struct {
	ID            string
	Auth          string // "customer" or "guest"
	PartnerUserID string // empty for a guest
	Scope         []string
}

// sessionStore keeps sessions in memory. A restart logs everyone out; the
// page then simply asks the app to reopen the mini app.
type sessionStore struct {
	mu   sync.Mutex
	byID map[string]*session
}

func newSessionStore() *sessionStore {
	return &sessionStore{byID: map[string]*session{}}
}

func (s *sessionStore) create(auth, partnerUserID string, scope []string) *session {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		panic(err)
	}
	sess := &session{ID: hex.EncodeToString(raw[:]), Auth: auth, PartnerUserID: partnerUserID, Scope: scope}
	s.mu.Lock()
	s.byID[sess.ID] = sess
	s.mu.Unlock()
	return sess
}

func (s *sessionStore) get(id string) (*session, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.byID[id]
	return sess, ok
}
