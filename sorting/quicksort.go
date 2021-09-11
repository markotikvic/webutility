package sorting

import "sort"

// QuickSort quicksorts que.
func QuickSort(que sort.Interface, low, high int, reverse bool) {
	if low >= high {
		return
	}

	index := quickSortPartition(que, low, high, reverse)
	QuickSort(que, low, index-1, reverse)
	QuickSort(que, index+1, high, reverse)
}

func quickSortPartition(que sort.Interface, low, high int, reverse bool) int {
	i := low - 1
	for j := low; j <= high-1; j++ {
		if !reverse {
			if que.Less(j, high) {
				i++
				que.Swap(i, j)
			}
		} else {
			if !que.Less(j, high) {
				i++
				que.Swap(i, j)
			}
		}
	}
	que.Swap(i+1, high)
	return i + 1
}

func BubbleSort(arr []int64) {
	for i := 0; i < len(arr)-1; i++ {
		for j := i; j < len(arr); j++ {
			if arr[i] > arr[j] {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
}
