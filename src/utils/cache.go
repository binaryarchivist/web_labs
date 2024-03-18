package utils

var Cache = make(map[string]string)

func Set(key, value string) {
	Cache[key] = value
}

func Get(key string) (string, bool) {
	val, found := Cache[key]
	return val, found
}
