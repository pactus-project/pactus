package sqlitestorage

// SQLite settings.
type settings struct {
	lockingMode string
}

func defaultSettings() settings {
	return settings{
		lockingMode: "EXCLUSIVE",
	}
}

// Option configures SQLite settings.
type Option func(*settings)

// WithLockingMode sets SQLite locking mode (e.g. NORMAL, EXCLUSIVE).
func WithLockingMode(lockMode bool) Option {
	return func(p *settings) {
		if lockMode {
			p.lockingMode = "EXCLUSIVE"
		} else {
			p.lockingMode = "NORMAL"
		}
	}
}
