package appcache

type AppCache struct {
	data map[string]string
}

func NewCache() *AppCache {
	return &AppCache{
		data: make(map[string]string),
	}
}

func (r *AppCache) SetData(key string, value string) bool {
	_, ok := r.data[key]
	if ok {
		return false
	} else {
		r.data[key] = value
		return true
	}
}

func (r *AppCache) GetData(key string) (string, bool) {

	value, ok := r.data[key]

	if ok {
		return value, ok
	} else {
		return "", false
	}
}
