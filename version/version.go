package version

// Version components
const (
	Maj = "0"
	Min = "7"
	Fix = "0"
)

var (
	// Version is the current version of Zarb in string
	Version = "0.7.0"

	// GitCommit is the current HEAD set using ldflags.
	GitCommit string
)

func init() {
	if GitCommit != "" {
		Version += "-" + GitCommit
	}
}
