package util

type (
	Daemon interface {
		Start() error
		Stop() error
	}
)
