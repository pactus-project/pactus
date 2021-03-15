package status

import (
	"fmt"

	"github.com/sasha-s/go-deadlock"
)

type Status struct {
	lk deadlock.RWMutex

	isProposed     bool
	isPrepared     bool
	isPreCommitted bool
	isCommitted    bool
}

func NewStatus() *Status {
	return &Status{}
}

func (s *Status) IsProposed() bool {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.isProposed
}

func (s *Status) IsPrepared() bool {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.isPrepared
}

func (s *Status) IsPreCommitted() bool {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.isPreCommitted
}

func (s *Status) IsCommitted() bool {
	s.lk.Lock()
	defer s.lk.Unlock()

	return s.isCommitted
}

func (s *Status) SetProposed(isProposed bool) {
	s.lk.Lock()
	defer s.lk.Unlock()

	s.isProposed = isProposed
}

func (s *Status) SetPrepared(isPrepared bool) {
	s.lk.Lock()
	defer s.lk.Unlock()

	s.isPrepared = isPrepared
}

func (s *Status) SetPreCommitted(isPreCommitted bool) {
	s.lk.Lock()
	defer s.lk.Unlock()

	s.isPreCommitted = isPreCommitted
}

func (s *Status) SetCommitted(isCommitted bool) {
	s.lk.Lock()
	defer s.lk.Unlock()

	s.isCommitted = isCommitted
}

func (s *Status) String() string {
	s.lk.Lock()
	defer s.lk.Unlock()

	isProposed := "-"
	if s.isProposed {
		isProposed = "X"
	}
	isPrepared := "-"
	if s.isPrepared {
		isPrepared = "X"
	}
	isPreCommitted := "-"
	if s.isPreCommitted {
		isPreCommitted = "X"
	}
	isCommitted := "-"
	if s.isCommitted {
		isCommitted = "X"
	}
	return fmt.Sprintf("%s%s%s%s", isProposed, isPrepared, isPreCommitted, isCommitted)
}
