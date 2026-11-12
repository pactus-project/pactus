package session

import (
	"time"

	"github.com/pactus-project/pactus/sync/peerset/peer"
	"github.com/pactus-project/pactus/types"
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

func (sm *Manager) OpenSession(pid peer.ID, from types.Height, count uint32) *Session {
	ssn := NewSession(sm.nextSessionID, pid, from, count)
	sm.sessions[ssn.SessionID] = ssn
	sm.nextSessionID++

	return ssn
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

// maxSessionDurationFactor bounds the duration of a session, as a number of session timeouts.
// A peer can keep a session active by answering without blocks. As long as a session is open,
// no other blocks are requested, so a session that lasts too long is closed, whatever its activity.
// An honest peer sends the blocks of a session within a few seconds.
const maxSessionDurationFactor = 12

// SetExpiredSessionsAsUncompleted marks as uncompleted the open sessions without recent activity,
// and the ones that last too long.
// Completed sessions have downloaded their blocks, and they don't expire.
func (sm *Manager) SetExpiredSessionsAsUncompleted() {
	maxDuration := maxSessionDurationFactor * sm.sessionTimeout
	for _, ssn := range sm.sessions {
		if ssn.Status != Open {
			continue
		}

		if sm.sessionTimeout < time.Since(ssn.LastActivity) || maxDuration < time.Since(ssn.StartedAt) {
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

// SetSessionOnlyFullNodes records that the blocks of the session
// must be requested from full nodes only.
func (sm *Manager) SetSessionOnlyFullNodes(sid int) {
	ssn := sm.sessions[sid]
	if ssn != nil {
		ssn.OnlyFullNodes = true
	}
}

// RemoveSession removes the session with the given id.
func (sm *Manager) RemoveSession(sid int) {
	delete(sm.sessions, sid)
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
