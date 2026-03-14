package policy

import "time"

type Policy interface {
	Add(key string, value string, timeout time.Duration)
	Get(key string) (string, bool)
}
