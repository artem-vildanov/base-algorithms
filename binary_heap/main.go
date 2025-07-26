package main

import (
	"errors"
)

type MinBinaryHeap struct {
	heap []int
}

func NewEmptyBinaryHeap() *MinBinaryHeap {
	return &MinBinaryHeap{
		heap: make([]int, 0),
	}
}

func (h *MinBinaryHeap) GetRoot() int {
	return h.heap[0]
}

func (h *MinBinaryHeap) ExtractRoot() (int, error) {
	if len(h.heap) == 0 {
		return 0, errors.New("failed to extract root: heap is empty")
	}

	if len(h.heap) == 1 {
		extracted := h.heap[0]
		h.heap = []int{}
		return extracted, nil
	}

	var (
		root         = h.heap[0]
		currentIndex = 0
	)

	h.heap[0] = h.heap[len(h.heap)-1] // меняем первый элемент на последний
	h.heap = h.heap[:len(h.heap)-1]   // отрезаем последний элемент

	for {
		currentValue := h.heap[currentIndex]
		leftChildIndex, rightChildIndex := h.getChildrenIndeces(currentIndex)

		leftChildIsLess := leftChildIndex < len(h.heap) &&
			h.heap[leftChildIndex] < currentValue

		rightChildIsLess := rightChildIndex < len(h.heap) &&
			h.heap[rightChildIndex] < currentValue

		if leftChildIsLess {
			h.heap[currentIndex] = h.heap[leftChildIndex]
			h.heap[leftChildIndex] = currentValue
			currentIndex = leftChildIndex
		} else if rightChildIsLess {
			h.heap[currentIndex] = h.heap[rightChildIndex]
			h.heap[rightChildIndex] = currentValue
			currentIndex = rightChildIndex
		} else {
			break
		}
	}

	return root, nil
}

func (h *MinBinaryHeap) Add(newValue int) {
	if len(h.heap) == 0 {
		h.heap = append(h.heap, newValue)
		return
	}

	newValueIndex := len(h.heap)
	h.heap = append(h.heap, newValue)

	parentIndex := h.getParentIndex(newValueIndex)
	parentValue := h.heap[parentIndex]

	if parentValue <= newValue {
		return
	}

	for {
		h.heap[parentIndex] = newValue
		h.heap[newValueIndex] = parentValue

		newValueIndex = parentIndex

		parentIndex = h.getParentIndex(newValueIndex)
		parentValue = h.heap[parentIndex]

		if parentValue <= newValue {
			break
		}
	}
}

func (h *MinBinaryHeap) getParentIndex(i int) (
	parentIndex int,
) {
	parentIndex = (i - 1) / 2
	return
}

func (h *MinBinaryHeap) getChildrenIndeces(i int) (
	leftChildIndex int,
	rightChildIndex int,
) {
	leftChildIndex = 2*i + 1
	rightChildIndex = 2*i + 2
	return
}

func main() {

}
