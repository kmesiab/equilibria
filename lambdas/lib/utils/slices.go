package utils

func FilterSlice[T any](s []T, condition func(T) bool) []T {
	var r []T
	for _, item := range s {
		if condition(item) {
			r = append(r, item)
		}
	}
	return r
}
