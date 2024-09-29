package helpers

func MapKeys[K comparable, V int64 | string | interface{}](i map[K]V) []K {
	keys := []K{}
	for k := range i {
		keys = append(keys, k)
	}
	return keys
}
