package internal

type Cache interface {
	Set(key string, value any)
	Get(key string) (any, bool)
}

type localCache map[string]any

func LocalCache() Cache {
	return localCache(make(map[string]any))
}

func (l localCache) Set(key string, value any) {
	l[key] = value
}

func (l localCache) Get(key string) (any, bool) {
	v, ok := l[key]
	return v, ok
}
