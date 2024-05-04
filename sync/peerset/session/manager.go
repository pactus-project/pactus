package session

import (
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/util"
)

type Manager struct {
	sessionTimeout time.Duration
	sessions       map[int]*Session
	nextSessionID  int
}

func NewManager(sessionTimeout time.Duration) *Manager {
	return &Manager{
		sessionTimeout: sessionTimeout,
		sessions:       make(map[int]*Session),
	}
}

func (sm *Manager) Stats() Stats {
	total := len(sm.sessions)
	open := 0
	completed := 0
	unCompleted := 0
	for _, ssn := range sm.sessions {
		switch ssn.Status {
		case Open:
			open++

		case Completed:
			completed++

		case Uncompleted:
			unCompleted++
		}
	}

	return Stats{
		Total:       total,
		Open:        open,
		Completed:   completed,
		Uncompleted: unCompleted,
	}
}

func (sm *Manager) OpenSession(pid peer.ID, from, count uint32) *Session {
	ssn := NewSession(sm.nextSessionID, pid, from, count)
	sm.sessions[ssn.SessionID] = ssn
	sm.nextSessionID++

	return ssn
}

func (sm *Manager) FindSession(sid int) *Session {
	ssn, ok := sm.sessions[sid]
	if ok {
		return ssn
	}

	return nil
}

func (sm *Manager) NumberOfSessions() int {
	return len(sm.sessions)
}

func (sm *Manager) HasOpenSession(pid peer.ID) bool {
	for _, ssn := range sm.sessions {
		if ssn.PeerID == pid && ssn.Status == Open {
			return true
		}
	}

	return false
}

func (sm *Manager) HasAnyOpenSession() bool {
	for _, ssn := range sm.sessions {
		if ssn.Status == Open {
			return true
		}
	}

	return false
}

func (sm *Manager) UpdateSessionLastActivity(sid int) {
	ssn := sm.sessions[sid]
	if ssn != nil {
		ssn.LastActivity = time.Now()
	}
}

func (sm *Manager) SetExpiredSessionsAsUncompleted() {
	for _, ssn := range sm.sessions {
		if sm.sessionTimeout < util.Now().Sub(ssn.LastActivity) {
			ssn.Status = Uncompleted
		}
	}
}

func (sm *Manager) SetSessionUncompleted(sid int) {
	ssn := sm.sessions[sid]
	if ssn != nil {
		ssn.Status = Uncompleted
	}
}

func (sm *Manager) SetSessionCompleted(sid int) *Session {
	ssn := sm.sessions[sid]
	if ssn != nil {
		ssn.Status = Completed
	}

	return ssn
}

func (sm *Manager) RemoveAllSessions() {
	sm.sessions = make(map[int]*Session)
}

func (sm *Manager) Sessions() []*Session {
	sessions := make([]*Session, 0, len(sm.sessions))

	for _, ssn := range sm.sessions {
		sessions = append(sessions, ssn)
	}

	return sessions
}
