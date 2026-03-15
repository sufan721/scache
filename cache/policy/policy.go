package policy

import "time"

//实现接口

type Policy interface {
	Add(key string, value string, timeout time.Duration)
	Get(key string) (string, bool)
}
