package sort

import (
	"github.com/vistormu/go-dsa/constraints"
)

func QuickSort[T constraints.Ordered](arr []T) {
	if len(arr) < 2 {
		return
	}

	quickSort(arr, 0, len(arr)-1)
}

func quickSort[T constraints.Ordered](arr []T, low, high int) {
	for low < high {
		pivotIndex := partition(arr, low, high)

		if pivotIndex-low < high-pivotIndex {
			quickSort(arr, low, pivotIndex-1)
			low = pivotIndex + 1
		} else {
			quickSort(arr, pivotIndex+1, high)
			high = pivotIndex - 1
		}
	}
}

func partition[T constraints.Ordered](arr []T, low, high int) int {
	pivot := arr[high]
	i := low

	for j := low; j < high; j++ {
		if arr[j] < pivot {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}

	arr[i], arr[high] = arr[high], arr[i]

	return i
}
