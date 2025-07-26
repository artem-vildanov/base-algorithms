package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBinaryHeap_Add(t *testing.T) {
	testCases := []struct {
		name         string
		processHeap  func(h *MinBinaryHeap)
		expectedHeap []int
	}{
		{
			name: "добавление первого элемента",
			processHeap: func(h *MinBinaryHeap) {
				h.Add(1)
			},
			expectedHeap: []int{1},
		},
		{
			name: "вставка большего элемента",
			processHeap: func(h *MinBinaryHeap) {
				h.Add(3)
				h.Add(7)
			},
			expectedHeap: []int{3, 7},
		},
		{
			name: "вставка меньшего элемента",
			processHeap: func(h *MinBinaryHeap) {
				h.Add(7)
				h.Add(3)
			},
			expectedHeap: []int{3, 7},
		},
		{
			name: "множественные вставки с каскадным всплытием",
			processHeap: func(h *MinBinaryHeap) {
				h.Add(5)
				h.Add(6)
				h.Add(7)
				h.Add(1)
			},
			expectedHeap: []int{1, 5, 7, 6},
		},
		{
			name: "вставка элементов в случайном порядке",
			processHeap: func(h *MinBinaryHeap) {
				h.Add(10)
				h.Add(4)
				h.Add(15)
				h.Add(20)
				h.Add(0)
				h.Add(7)
			},
			expectedHeap: []int{0, 4, 7, 20, 10, 15},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := NewEmptyBinaryHeap()
			tc.processHeap(h)
			assert.Equal(t, tc.expectedHeap, h.heap)
		})
	}
}

func TestBinaryHeap_ExtractRoot(t *testing.T) {
	emptyHeapErrMsg := "failed to extract root: heap is empty"

	testCases := []struct {
		name                  string
		prepareHeap           func(*MinBinaryHeap)
		expectedAfterExctract []int
		expectedExtracted     int
		expectedErrMsg        *string
	}{
		{
			name: "извлечение первого элемента",
			prepareHeap: func(h *MinBinaryHeap) {
				h.Add(1)
			},
			expectedAfterExctract: make([]int, 0),
			expectedExtracted:     1,
		},
		{
			name: "извлечение из кучи с несколькими элементами",
			prepareHeap: func(h *MinBinaryHeap) {
				h.Add(10)
				h.Add(4)
				h.Add(15)
				h.Add(0)
				h.Add(7)
			},
			expectedAfterExctract: []int{4, 7, 15, 10},
			expectedExtracted:     0,
		},
		{
			name:                  "извлечение из пустой кучи",
			prepareHeap:           func(h *MinBinaryHeap) {},
			expectedAfterExctract: []int{4, 7, 15, 10},
			expectedErrMsg:        &emptyHeapErrMsg,
		},
		{
			name: "извлечение из кучи с одинаковыми элементами",
			prepareHeap: func(h *MinBinaryHeap) {
				h.Add(3)
				h.Add(3)
				h.Add(3)
			},
			expectedAfterExctract: []int{3, 3},
			expectedExtracted:     3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := NewEmptyBinaryHeap()

			tc.prepareHeap(h)

			gotExtracted, err := h.ExtractRoot()
			if tc.expectedErrMsg == nil {
				require.NoError(t, err)
			} else {
				require.NotEmpty(t, err)
				assert.Equal(t, *tc.expectedErrMsg, err.Error())
				return
			}

			assert.Equal(t, tc.expectedExtracted, gotExtracted)
			assert.Equal(t, tc.expectedAfterExctract, h.heap)
		})
	}
}
