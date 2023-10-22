package utils

func CompL[T any, V any](objs []T, lambda func(x T) V) []V {
	results := make([]V, 0, len(objs))
	for _, obj := range objs {
		results = append(results, lambda(obj))
	}
	return results
}
