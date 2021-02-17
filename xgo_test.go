package xgo_test

import (
	"testing"

	"github.com/glassonion1/xgo"
)

func TestIsFirstUpper(t *testing.T) {

	tests := []struct {
		name string
		in   string
		want bool
	}{
		{
			name: "uppercase",
			in:   "ABCD",
			want: true,
		},
		{
			name: "uppercase",
			in:   "Abcd",
			want: true,
		},
		{
			name: "uppercase",
			in:   "A",
			want: true,
		},
		{
			name: "empty",
			in:   "",
			want: false,
		},
		{
			name: "lowercase",
			in:   "abc",
			want: false,
		},
		{
			name: "lowercase",
			in:   "a",
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := xgo.IsFirstUpper(tt.in)
			if tt.want != got {
				t.Errorf("testing %s: faild want: %v, got: %v", tt.name, tt.want, got)
			}
		})
	}
}
