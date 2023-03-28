package xgo

import (
	"errors"
	"fmt"
)

func Map[S, D any](input []S, f func(S) (D, error)) ([]D, error) {
	var output []D
	var errs []error
	for i, v := range input {
		o, err := f(v)
		if err != nil {
			errs = append(errs, fmt.Errorf("index: %d, err: %w", i, err))
			continue
		}
		output = append(output, o)
	}
	return output, errors.Join(errs...)
}
