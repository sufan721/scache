package cache

type Policy interface {
	Add(key string, value string)
	Get(key string) (string, bool)
}
