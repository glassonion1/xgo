package xgo

func Map[S, D any](input []S, f func(S) (D, error)) ([]D, error) {
	output := make([]D, len(input))
	for i, v := range input {
		o, err := f(v)
		if err != nil {
			return nil, err
		}
		output[i] = o
	}
	return output, nil
}
