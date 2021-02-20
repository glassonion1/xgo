package xgo_test

import (
	"testing"

	"github.com/glassonion1/xgo"
	"github.com/google/go-cmp/cmp"
)

func TestSplitChunk(t *testing.T) {
	type args struct {
		length    int
		chunkSize int
	}

	tests := []struct {
		name string
		in   args
		want []xgo.Chunk
	}{
		{
			name: "splits chunks by evne",
			in: args{
				length:    6,
				chunkSize: 2,
			},
			want: []xgo.Chunk{
				xgo.Chunk{
					From: 0,
					To:   2,
				},
				xgo.Chunk{
					From: 2,
					To:   4,
				},
				xgo.Chunk{
					From: 4,
					To:   6,
				},
			},
		},
		{
			name: "splits chunks by even, length has a remainder",
			in: args{
				length:    5,
				chunkSize: 2,
			},
			want: []xgo.Chunk{
				xgo.Chunk{
					From: 0,
					To:   2,
				},
				xgo.Chunk{
					From: 2,
					To:   4,
				},
				xgo.Chunk{
					From: 4,
					To:   5,
				},
			},
		},
		{
			name: "splits chunks by odd, length has a remainder",
			in: args{
				length:    10,
				chunkSize: 3,
			},
			want: []xgo.Chunk{
				xgo.Chunk{
					From: 0,
					To:   3,
				},
				xgo.Chunk{
					From: 3,
					To:   6,
				},
				xgo.Chunk{
					From: 6,
					To:   9,
				},
				xgo.Chunk{
					From: 9,
					To:   10,
				},
			},
		},
		{
			name: "can not split chunks",
			in: args{
				length:    5,
				chunkSize: 10,
			},
			want: []xgo.Chunk{
				xgo.Chunk{
					From: 0,
					To:   5,
				},
			},
		},
		{
			name: "chunk size is 0",
			in: args{
				length:    3,
				chunkSize: 0,
			},
			want: []xgo.Chunk{
				xgo.Chunk{
					From: 0,
					To:   1,
				},
				xgo.Chunk{
					From: 1,
					To:   2,
				},
				xgo.Chunk{
					From: 2,
					To:   3,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cgot := xgo.SplitChunks(tt.in.length, tt.in.chunkSize)
			var got []xgo.Chunk
			for i := range cgot {
				got = append(got, i)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("testing %s mismatch (-want +got):\n%s\n", tt.name, diff)
			}
		})
	}

	t.Run("tests slice", func(t *testing.T) {
		list := []int{2, 4, 6, 8, 10, 12, 14, 16}
		wants := [][]int{
			[]int{2, 4, 6},
			[]int{8, 10, 12},
			[]int{14, 16},
		}
		index := 0
		for chunk := range xgo.SplitChunks(len(list), 3) {
			got := list[chunk.From:chunk.To]
			if diff := cmp.Diff(wants[index], got); diff != "" {
				t.Errorf("testing %s mismatch (-want +got):\n%s\n", "tests slice", diff)
			}
			index++
		}
	})
}
