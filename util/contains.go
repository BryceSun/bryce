package util

func Load[T any](m map[string]any, key string) (rt T) {
	v := m[key]
	if v == nil {
		return rt
	}
	rt = v.(T)
	return rt
}

func Keys[K comparable, V any](m map[K]V) []K {
	s := make([]K, 0, len(m))
	for k := range m {
		s = append(s, k)
	}
	return s
}

func Values[K comparable, V any](m map[K]V) []V {
	s := make([]V, 0, len(m))
	for _, v := range m {
		s = append(s, v)
	}
	return s
}
