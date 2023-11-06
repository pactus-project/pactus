package testsuite

import (
	"fmt"
)

type OverrideStringer struct {
	obj  fmt.Stringer
	name string
}

func NewOverrideStringer(name string, obj fmt.Stringer) *OverrideStringer {
	return &OverrideStringer{
		obj:  obj,
		name: name,
	}
}

func (o *OverrideStringer) String() string {
	return o.name + o.obj.String()
}
