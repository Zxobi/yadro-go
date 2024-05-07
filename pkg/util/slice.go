package util

func SliceToAny[T any](slice []T) []any {
	s := make([]interface{}, len(slice))
	for i, v := range slice {
		s[i] = v
	}

	return s
}
