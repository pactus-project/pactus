package version

import (
	"fmt"
)

// Version components
const (
	Maj = 0
	Min = 7
	Fix = 0
)

var (
	Version   = fmt.Sprintf("%d.%d.%d", Maj, Min, Fix)
	GitCommit string
)

func init() {
	if GitCommit != "" {
		Version += "-" + GitCommit
	}
}
