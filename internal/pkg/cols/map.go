package cols

func Map[In any, Out any](in []In, f func(i In) (o Out)) (out []Out) {
	for _, element := range in {
		out = append(out, f(element))
	}
	return
}

func MapErr[In any, Out any](in []In, f func(i In) (o Out, err error)) (out []Out, err error) {
	for _, i := range in {
		var o Out
		o, err = f(i)
		if err != nil {
			return
		}
		out = append(out, o)
	}
	return
}
