package score

import "github.com/pactus-project/pactus/types/certificate"

type scoreData struct {
	ShouldVote int
	NotVote    int
}

type ScoreManager struct {
	certs   map[uint32]*certificate.Certificate
	vals    map[int32]*scoreData
	maxCert uint32
}

func NewScoreManager(maxCert uint32) *ScoreManager {
	return &ScoreManager{
		certs:   make(map[uint32]*certificate.Certificate),
		vals:    make(map[int32]*scoreData),
		maxCert: maxCert,
	}
}

func (sm *ScoreManager) SetCertificate(cert *certificate.Certificate) {
	lastHeight := cert.Height()
	sm.certs[lastHeight] = cert

	for _, num := range cert.Committers() {
		data, ok := sm.vals[num]
		if !ok {
			data = new(scoreData)
			sm.vals[num] = data
		}

		data.ShouldVote++
	}

	for _, num := range cert.Absentees() {
		data := sm.vals[num]
		sm.vals[num] = data

		data.NotVote++
	}

	oldHeight := lastHeight - uint32(sm.maxCert)
	oldCert, ok := sm.certs[oldHeight]
	if ok {
		for _, num := range oldCert.Committers() {
			data := sm.vals[num]
			data.ShouldVote--
		}

		for _, num := range oldCert.Absentees() {
			data := sm.vals[num]
			data.NotVote--
		}

		delete(sm.certs, oldHeight)
	}
}

func (sm *ScoreManager) AvailabilityScore(valNum int32) float64 {
	data, ok := sm.vals[valNum]
	if ok {
		if data.ShouldVote == 0 {
			return 1.0
		} else {
			return 1 - (float64(data.NotVote) / float64(data.ShouldVote))
		}
	}

	return 1.0
}
