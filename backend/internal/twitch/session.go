package twitch

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

type OAuthSession struct {
	UserID    string
	State     string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type SessionStore struct {
	sessions map[string]*OAuthSession
	mutex    sync.RWMutex
}

var globalSessionStore *SessionStore
var sessionStoreOnce sync.Once

func GetSessionStore() *SessionStore {
	sessionStoreOnce.Do(func() {
		globalSessionStore = NewSessionStore()
	})
	return globalSessionStore
}

func NewSessionStore() *SessionStore {
	store := &SessionStore{
		sessions: make(map[string]*OAuthSession),
	}

	// Start cleanup goroutine for expired sessions
	go store.cleanupExpiredSessions()

	return store
}

func (s *SessionStore) CreateSession(userID string) (string, error) {
	// Generate secure state parameter
	state, err := generateSecureState()
	if err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}

	session := &OAuthSession{
		UserID:    userID,
		State:     state,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(10 * time.Minute), // 10 minute expiry
	}

	s.mutex.Lock()
	s.sessions[state] = session
	s.mutex.Unlock()

	return state, nil
}

func (s *SessionStore) GetSession(state string) (*OAuthSession, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	session, exists := s.sessions[state]
	if !exists {
		return nil, false
	}

	// Check if session has expired
	if time.Now().After(session.ExpiresAt) {
		s.DeleteSession(state)
		return nil, false
	}

	return session, true
}

func (s *SessionStore) DeleteSession(state string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.sessions, state)
}

func (s *SessionStore) cleanupExpiredSessions() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mutex.Lock()
		now := time.Now()
		for state, session := range s.sessions {
			if now.After(session.ExpiresAt) {
				delete(s.sessions, state)
			}
		}
		s.mutex.Unlock()
	}
}

func generateSecureState() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
