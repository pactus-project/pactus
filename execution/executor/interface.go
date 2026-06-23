package executor

import "github.com/pactus-project/pactus/sandbox"

type Executor interface {
	Check(sbx sandbox.SandboxReader, strict bool) error
	Execute(sbx sandbox.Sandbox)
}
