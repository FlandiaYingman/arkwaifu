package cols

func Filter[T any](in []T, predicate func(element T) bool) (out []T) {
	for _, element := range in {
		if predicate(element) {
			out = append(out, element)
		}
	}
	return
}

func FilterErr[T any](in []T, predicate func(element T) (bool, error)) (out []T, err error) {
	for _, element := range in {
		ok, err := predicate(element)
		if err != nil {
			return nil, err
		}
		if ok {
			out = append(out, element)
		}
	}
	return
}
