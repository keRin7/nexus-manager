package appcache

type AppCache struct {
	data map[string]elem
}

type elem struct {
	data      string
	layersSHA string
}

func NewCache() *AppCache {
	return &AppCache{
		data: make(map[string]elem),
	}
}

func (r *AppCache) SetData(key string, layerSHA string, value string) bool {

	_, ok := r.data[key]

	if ok {
		return false
	} else {
		r.data[key] = elem{data: value, layersSHA: layerSHA}
		return true
	}

}

func (r *AppCache) GetData(key string) (data string, sha string, ok bool) {

	value, ok := r.data[key]

	if ok {
		return value.data, value.layersSHA, ok
	} else {
		return "", "", false
	}
}
