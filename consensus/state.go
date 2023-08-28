package consensus

type consState interface {
	enter()
	decide()
	timeout(t *ticker)
	name() string
}
