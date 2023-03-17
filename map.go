package xgo

import "errors"

func Map[S, D any](input []S, f func(S) (D, error)) ([]D, error) {
	output := make([]D, len(input))
	var errs []error
	for i, v := range input {
		o, err := f(v)
		if err != nil {
			errs = append(errs, err)
		}
		output[i] = o
	}
	return output, errors.Join(errs...)
}
