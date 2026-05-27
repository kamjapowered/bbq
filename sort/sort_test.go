package sort

import (
	"github.com/vistormu/go-dsa/constraints"
	"testing"
)

type sortFn[T constraints.Ordered] struct {
	name string
	fn   func([]T)
}

func sort[T constraints.Ordered](arr []T, fn func([]T)) {
	fn(arr)
}

func testSort[T constraints.Ordered](t *testing.T, arr []T, expected []T, fn func([]T)) {
	t.Run("Sorting", func(t *testing.T) {
		arrCopy := make([]T, len(arr))
		copy(arrCopy, arr)
		sort(arrCopy, fn)
		for i, v := range arrCopy {
			if v != expected[i] {
				t.Errorf("Expected %v at index %d, got %v", expected[i], i, v)
			}
		}
	})
}

func TestIntSort(t *testing.T) {
	fns := []sortFn[int]{
		{"QuickSort", QuickSort[int]},
	}

	arr := []int{5, 2, 9, 1, 5, 6}
	expected := []int{1, 2, 5, 5, 6, 9}

	for _, fn := range fns {
		t.Run(fn.name, func(t *testing.T) {
			testSort(t, arr, expected, fn.fn)
		})
	}
}

func TestStringSort(t *testing.T) {
	fns := []sortFn[string]{
		{"QuickSort", QuickSort[string]},
	}

	arr := []string{"banana", "apple", "cherry", "date"}
	expected := []string{"apple", "banana", "cherry", "date"}

	for _, fn := range fns {
		t.Run(fn.name, func(t *testing.T) {
			testSort(t, arr, expected, fn.fn)
		})
	}
}

func TestFloatSort(t *testing.T) {
	fns := []sortFn[float64]{
		{"QuickSort", QuickSort[float64]},
	}

	arr := []float64{3.1, 2.4, 5.6, 1.2}
	expected := []float64{1.2, 2.4, 3.1, 5.6}

	for _, fn := range fns {
		t.Run(fn.name, func(t *testing.T) {
			testSort(t, arr, expected, fn.fn)
		})
	}
}
