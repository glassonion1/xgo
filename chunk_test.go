package xgo_test

import (
	"testing"

	"github.com/glassonion1/xgo"
	"github.com/google/go-cmp/cmp"
)

func TestChunkIndex(t *testing.T) {

	type args struct {
		length    int
		chunkSize int
	}

	tests := []struct {
		name string
		in   args
		want []xgo.ChunkedIndex
	}{
		{
			name: "splits chunks by evne",
			in: args{
				length:    6,
				chunkSize: 2,
			},
			want: []xgo.ChunkedIndex{
				{
					From: 0,
					To:   2,
				},
				{
					From: 2,
					To:   4,
				},
				{
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
			want: []xgo.ChunkedIndex{
				{
					From: 0,
					To:   2,
				},
				{
					From: 2,
					To:   4,
				},
				{
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
			want: []xgo.ChunkedIndex{
				{
					From: 0,
					To:   3,
				},
				{
					From: 3,
					To:   6,
				},
				{
					From: 6,
					To:   9,
				},
				{
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
			want: []xgo.ChunkedIndex{
				{
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
			want: []xgo.ChunkedIndex{
				{
					From: 0,
					To:   1,
				},
				{
					From: 1,
					To:   2,
				},
				{
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
			cgot := xgo.ChunkIndex(tt.in.length, tt.in.chunkSize)
			var got []xgo.ChunkedIndex
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
			{2, 4, 6},
			{8, 10, 12},
			{14, 16},
		}
		index := 0
		for chunk := range xgo.ChunkIndex(len(list), 3) {
			got := list[chunk.From:chunk.To]
			if diff := cmp.Diff(wants[index], got); diff != "" {
				t.Errorf("testing %s mismatch (-want +got):\n%s\n", "tests slice", diff)
			}
			index++
		}
	})
}

func TestChunk(t *testing.T) {
	t.Run("tests slice of int", func(t *testing.T) {
		list := []int{2, 4, 6, 8, 10, 12, 14, 16}
		want := [][]int{
			{2, 4, 6},
			{8, 10, 12},
			{14, 16},
		}
		got := xgo.Chunk[int](list, 3)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("testing %s mismatch (-want +got):\n%s\n", "tests slice", diff)
		}
	})

	t.Run("tests slice of struct", func(t *testing.T) {
		type item struct {
			ID   string
			Name string
		}
		list := []item{
			{
				ID:   "1",
				Name: "test1",
			},
			{
				ID:   "2",
				Name: "test2",
			},
			{
				ID:   "3",
				Name: "test3",
			},
			{
				ID:   "4",
				Name: "test4",
			},
			{
				ID:   "5",
				Name: "test5",
			},
			{
				ID:   "6",
				Name: "test6",
			},
		}
		want := [][]item{
			{
				{
					ID:   "1",
					Name: "test1",
				},
				{
					ID:   "2",
					Name: "test2",
				},
			},
			{
				{
					ID:   "3",
					Name: "test3",
				},
				{
					ID:   "4",
					Name: "test4",
				},
			},
			{
				{
					ID:   "5",
					Name: "test5",
				},
				{
					ID:   "6",
					Name: "test6",
				},
			},
		}
		got := xgo.Chunk[item](list, 2)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("testing %s mismatch (-want +got):\n%s\n", "tests slice", diff)
		}
	})

	t.Run("tests slice of pointer", func(t *testing.T) {
		type item struct {
			ID   string
			Name string
		}
		list := []*item{
			{
				ID:   "1",
				Name: "test1",
			},
			{
				ID:   "2",
				Name: "test2",
			},
			{
				ID:   "3",
				Name: "test3",
			},
			{
				ID:   "4",
				Name: "test4",
			},
			{
				ID:   "5",
				Name: "test5",
			},
			{
				ID:   "6",
				Name: "test6",
			},
		}
		want := [][]*item{
			{
				{
					ID:   "1",
					Name: "test1",
				},
				{
					ID:   "2",
					Name: "test2",
				},
			},
			{
				{
					ID:   "3",
					Name: "test3",
				},
				{
					ID:   "4",
					Name: "test4",
				},
			},
			{
				{
					ID:   "5",
					Name: "test5",
				},
				{
					ID:   "6",
					Name: "test6",
				},
			},
		}
		got := xgo.Chunk[*item](list, 2)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("testing %s mismatch (-want +got):\n%s\n", "tests slice", diff)
		}
	})
}
