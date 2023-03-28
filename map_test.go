package xgo_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/glassonion1/xgo"
	"github.com/google/go-cmp/cmp"
)

func TestMap(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := []string{"1", "2", "3", "4", "5"}

	// Define mapping function
	f := func(i int) (string, error) {
		return strconv.Itoa(i), nil
	}

	// Call Map function
	output, err := xgo.Map(input, f)

	// Verify output and error
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if diff := cmp.Diff(output, expected); diff != "" {
		t.Errorf("Output does not match expected (-got +want):\n%s", diff)
	}

	// Test error case
	input2 := []int{1, 2, 3, 4, 5}
	expectedErr := errors.New("failed to map")
	f2 := func(i int) (string, error) {
		if i == 3 {
			return "", expectedErr
		}
		return strconv.Itoa(i), nil
	}

	output2, err2 := xgo.Map(input2, f2)
	if !errors.Is(err2, expectedErr) {
		t.Errorf("Expected error %v but got %v", expectedErr, err2)
	}
	if output2 != nil {
		t.Errorf("Expected nil output but got %v", output2)
	}
}
