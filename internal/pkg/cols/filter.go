package cols

func Filter[T any](in []T, predicate func(element T) bool) (out []T) {
	for _, element := range in {
		if predicate(element) {
			out = append(out, element)
		}
	}
	return
}
