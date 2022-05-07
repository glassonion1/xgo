package xgo_test

import (
	"testing"
	"time"

	"github.com/glassonion1/xgo"
)

func TestContains(t *testing.T) {
	t.Run("slice of string test", func(t *testing.T) {
		t.Parallel()
		list := []string{"apple", "orange", "lemon"}
		elem := "lemon"
		got := xgo.Contains[string](list, elem)
		if !got {
			t.Errorf("testing %s: faild, this should be true", t.Name())
		}
	})
	t.Run("slice of int test", func(t *testing.T) {
		t.Parallel()
		list := []int{1, 2, 3, 4, 5}
		elem := 1
		got := xgo.Contains[int](list, elem)
		if !got {
			t.Errorf("testing %s: faild, this should be true", t.Name())
		}
	})

	t.Run("slice of float64 test", func(t *testing.T) {
		t.Parallel()
		list := []float64{1.1, 2.2, 3.3, 4.4, 5.5}
		elem := 3.3
		got := xgo.Contains[float64](list, elem)
		if !got {
			t.Errorf("testing %s: faild, this should be true", t.Name())
		}
	})

	type item struct {
		ID   string
		num  int
		Name string
	}
	t.Run("slice of struct test", func(t *testing.T) {
		t.Parallel()
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
				num:  3,
				Name: "test3",
			},
		}
		elem := item{
			ID:   "3",
			num:  3,
			Name: "test3",
		}
		got := xgo.Contains[item](list, elem)
		if !got {
			t.Errorf("testing %s: faild, this should be true", t.Name())
		}
		elem2 := item{
			ID:   "1",
			num:  2,
			Name: "test3",
		}
		got2 := xgo.Contains[item](list, elem2)
		if got2 {
			t.Errorf("testing %s: faild, this should be false", t.Name())
		}
	})

	now := time.Now()
	type order struct {
		ID        string
		CreatedAt time.Time
		UpdatedAt *time.Time
		Item      item
	}
	t.Run("slice of nested struct test", func(t *testing.T) {
		t.Parallel()
		list := []order{
			{
				ID:        "1",
				CreatedAt: now,
				Item: item{
					ID:   "1",
					Name: "test1",
				},
			},
			{
				ID:        "2",
				CreatedAt: now,
				UpdatedAt: nil,
				Item: item{
					ID:   "3",
					Name: "test3",
				},
			},
		}
		elem := order{
			ID:        "2",
			CreatedAt: now,
			Item: item{
				ID:   "3",
				Name: "test3",
			},
		}
		got := xgo.Contains[order](list, elem)
		if !got {
			t.Errorf("testing %s: faild, this should be true", t.Name())
		}

		elem2 := order{
			ID:        "2",
			CreatedAt: now,
			Item: item{
				ID:   "2",
				Name: "test2",
			},
		}
		got2 := xgo.Contains[order](list, elem2)
		if got2 {
			t.Errorf("testing %s: faild, this should be false", t.Name())
		}
	})
}
