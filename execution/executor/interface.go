package executor

type Executor interface {
	Check(strict bool) error
	Execute()
}
