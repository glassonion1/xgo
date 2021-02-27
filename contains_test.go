package xgo_test

import (
	"testing"
	"time"

	"github.com/glassonion1/xgo"
)

func TestContains(t *testing.T) {
	now := time.Now()
	type item struct {
		ID   string
		num  int
		Name string
	}
	type order struct {
		ID        string
		CreatedAt time.Time
		UpdatedAt *time.Time
		Item      item
	}
	type args struct {
		list interface{}
		elem interface{}
	}
	tests := []struct {
		name string
		in   args
		want bool
	}{
		{
			name: "slice of struct",
			in: args{
				list: []item{
					item{
						ID:   "1",
						Name: "test1",
					},
					item{
						ID:   "2",
						Name: "test2",
					},
					item{
						ID:   "3",
						num:  3,
						Name: "test3",
					},
				},
				elem: item{
					ID:   "3",
					num:  3,
					Name: "test3",
				},
			},
			want: true,
		},
		{
			name: "slice of struct 2",
			in: args{
				list: []order{
					order{
						ID:        "1",
						CreatedAt: now,
						Item: item{
							ID:   "1",
							Name: "test1",
						},
					},
					order{
						ID:        "2",
						CreatedAt: now,
						UpdatedAt: nil,
						Item: item{
							ID:   "3",
							Name: "test3",
						},
					},
				},
				elem: order{
					ID:        "2",
					CreatedAt: now,
					Item: item{
						ID:   "3",
						Name: "test3",
					},
				},
			},
			want: true,
		},
		{
			name: "slice of struct 3",
			in: args{
				list: []order{
					order{
						ID:        "1",
						CreatedAt: now,
						Item: item{
							ID:   "1",
							Name: "test1",
						},
					},
					order{
						ID:        "2",
						CreatedAt: now,
						UpdatedAt: &now,
						Item: item{
							ID:   "3",
							Name: "test3",
						},
					},
				},
				elem: order{
					ID:        "2",
					CreatedAt: now,
					UpdatedAt: &now,
					Item: item{
						ID:   "3",
						Name: "test3",
					},
				},
			},
			want: true,
		},
		{
			name: "slice of pointer",
			in: args{
				list: []*item{
					&item{
						ID:   "1",
						Name: "test1",
					},
					&item{
						ID:   "2",
						Name: "test2",
					},
					&item{
						ID:   "3",
						num:  3,
						Name: "test3",
					},
				},
				elem: &item{
					ID:   "3",
					num:  3,
					Name: "test3",
				},
			},
			want: true,
		},
		{
			name: "int32",
			in: args{
				list: []int32{1, 2, 3, 4, 5},
				elem: 4,
			},
			want: true,
		},
		{
			name: "int",
			in: args{
				list: []int{1, 2, 3, 4, 5},
				elem: 1,
			},
			want: true,
		},
		{
			name: "float64",
			in: args{
				list: []float64{1.1, 2.2, 3.3, 4.4, 5.5},
				elem: 3.3,
			},
			want: true,
		},
		{
			name: "string",
			in: args{
				list: []string{"apple", "orange", "lemon"},
				elem: "lemon",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := xgo.Contains(tt.in.list, tt.in.elem)
			if tt.want != got {
				t.Errorf("testing %s: faild want: %v, got: %v", tt.name, tt.want, got)
			}
		})
	}
}
