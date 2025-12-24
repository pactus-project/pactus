package testsuite

import (
	"github.com/pactus-project/pactus/util/logger"
)

type OverrideStringer struct {
	obj  logger.LogStringer
	name string
}

func NewOverrideLogStringer(name string, obj logger.LogStringer) *OverrideStringer {
	return &OverrideStringer{
		obj:  obj,
		name: name,
	}
}

func (o *OverrideStringer) LogString() string {
	return o.name + o.obj.LogString()
}
