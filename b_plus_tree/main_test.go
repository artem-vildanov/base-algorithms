package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Insert(t *testing.T) {
	const maxKeys = 4

	t.Run("лист является корнем", func(t *testing.T) {
		tree := NewTree(maxKeys)

		tree.Insert(1, "hello world")
		tree.Insert(2, 123)
		tree.Insert(3, false)
		tree.Insert(4, "zxcv")

		expectKeys := []int64{1, 2, 3, 4}
		expectValues := []any{"hello world", 123, false, "zxcv"}

		casted, ok := tree.root.(*LeafNode)
		require.Equal(t, true, ok)

		assert.Equal(t, expectKeys, casted.Keys)
		assert.Equal(t, expectValues, casted.values)
	})

	t.Run("переполнение листа, один внутренний узел", func(t *testing.T) {
		tree := NewTree(maxKeys)

		tree.Insert(1, "hello world")
		tree.Insert(2, 123)
		tree.Insert(3, false)
		tree.Insert(4, "zxcv")
		tree.Insert(5, "qwer")

		casted, ok := tree.root.(*InnerNode)
		require.Equal(t, true, ok)

		assert.Equal(t, []int64{3}, casted.Keys)
		require.Equal(t, 2, len(casted.Children))

		leftChild, ok := casted.Children[0].(*LeafNode)
		require.Equal(t, true, ok)

		rightChild, ok := casted.Children[1].(*LeafNode)
		require.Equal(t, true, ok)

		assert.Equal(t, []int64{1, 2}, leftChild.Keys)
		assert.Equal(t, []any{"hello world", 123}, leftChild.values)

		assert.Equal(t, []int64{3, 4, 5}, rightChild.Keys)
		assert.Equal(t, []any{false, "zxcv", "qwer"}, rightChild.values)
	})

	t.Run("каскадное переполнение листов и внутренних узлов", func(t *testing.T) {
		tree := NewTree(maxKeys)

		value10 := "asdf"
		value11 := "qwer"
		value12 := "zxcv"
		value14 := 123
		value15 := 1
		value16 := true
		value30 := "zxcvzxcvzxcv"
		value1 := "zxcvzxcvzxcvzxc"
		value2 := "lkasjdfklj"
		value9 := "uioiuoiuoiu"
		value8 := "ertyertyerty"
		value3 := "hdjfgdjfg"
		value7 := false
		value4 := 898989
		value5 := "zxvzxyeurty"

		tree.Insert(10, value10)
		tree.Insert(11, value11)
		tree.Insert(12, value12)
		tree.Insert(14, value14)
		tree.Insert(15, value15)
		tree.Insert(16, value16)
		tree.Insert(30, value30)
		tree.Insert(1, value1)
		tree.Insert(2, value2)
		tree.Insert(9, value9)
		tree.Insert(8, value8)
		tree.Insert(3, value3)
		tree.Insert(7, value7)
		tree.Insert(4, value4)
		tree.Insert(5, value5)

		/*
			проверка внутренних узлов
		*/
		castedRoot, ok := tree.root.(*InnerNode)
		require.True(t, ok)
		require.Equal(t, 2, len(castedRoot.Children))
		require.Equal(t, len(castedRoot.Keys), 1)
		assert.Equal(t, int64(9), castedRoot.Keys[0])

		castedLeftChild, ok := castedRoot.Children[0].(*InnerNode)
		require.True(t, ok)
		require.Equal(t, 2, len(castedLeftChild.Keys))
		assert.Equal(t, int64(3), castedLeftChild.Keys[0])
		assert.Equal(t, int64(5), castedLeftChild.Keys[1])

		castedRightChild, ok := castedRoot.Children[1].(*InnerNode)
		require.True(t, ok)
		require.Equal(t, 2, len(castedRightChild.Keys))
		assert.Equal(t, int64(12), castedRightChild.Keys[0])
		assert.Equal(t, int64(15), castedRightChild.Keys[1])

		/*
			проверка листов левого узла
		*/
		require.Equal(t, 3, len(castedLeftChild.Children))

		gotLeftLeaf, ok := castedLeftChild.Children[0].(*LeafNode)
		require.True(t, ok)
		gotMiddleLeaf, ok := castedLeftChild.Children[1].(*LeafNode)
		require.True(t, ok)
		gotRightLeaf, ok := castedLeftChild.Children[2].(*LeafNode)
		require.True(t, ok)

		assert.Equal(t, []int64{1, 2}, gotLeftLeaf.Keys)
		assert.Equal(t, []any{value1, value2}, gotLeftLeaf.values)
		assert.Equal(t, []int64{3, 4}, gotMiddleLeaf.Keys)
		assert.Equal(t, []any{value3, value4}, gotMiddleLeaf.values)
		assert.Equal(t, []int64{5, 7, 8}, gotRightLeaf.Keys)
		assert.Equal(t, []any{value5, value7, value8}, gotRightLeaf.values)

		assert.Equal(t, (*LeafNode)(nil), gotLeftLeaf.PrevLeaf)
		assert.Equal(t, gotLeftLeaf.NextLeaf.Keys, gotMiddleLeaf.Keys)
		assert.Equal(t, gotMiddleLeaf.PrevLeaf.Keys, gotLeftLeaf.Keys)
		assert.Equal(t, gotMiddleLeaf.NextLeaf.Keys, gotRightLeaf.Keys)
		assert.Equal(t, gotRightLeaf.PrevLeaf.Keys, gotMiddleLeaf.Keys)
		assert.Equal(t, []int64{9, 10, 11}, gotRightLeaf.NextLeaf.Keys)

		/*
			проверка листов правого узла
		*/
		require.Equal(t, 3, len(castedRightChild.Children))

		gotLeftLeaf, ok = castedRightChild.Children[0].(*LeafNode)
		require.True(t, ok)
		gotMiddleLeaf, ok = castedRightChild.Children[1].(*LeafNode)
		require.True(t, ok)
		gotRightLeaf, ok = castedRightChild.Children[2].(*LeafNode)
		require.True(t, ok)

		assert.Equal(t, []int64{9, 10, 11}, gotLeftLeaf.Keys)
		assert.Equal(t, []any{value9, value10, value11}, gotLeftLeaf.values)
		assert.Equal(t, []int64{12, 14}, gotMiddleLeaf.Keys)
		assert.Equal(t, []any{value12, value14}, gotMiddleLeaf.values)
		assert.Equal(t, []int64{15, 16, 30}, gotRightLeaf.Keys)
		assert.Equal(t, []any{value15, value16, value30}, gotRightLeaf.values)

		assert.Equal(t, []int64{5,7,8}, gotLeftLeaf.PrevLeaf.Keys)
		assert.Equal(t, gotLeftLeaf.NextLeaf.Keys, gotMiddleLeaf.Keys)
		assert.Equal(t, gotMiddleLeaf.PrevLeaf.Keys, gotLeftLeaf.Keys)
		assert.Equal(t, gotMiddleLeaf.NextLeaf.Keys, gotRightLeaf.Keys)
		assert.Equal(t, gotRightLeaf.PrevLeaf.Keys, gotMiddleLeaf.Keys)
		assert.Equal(t, (*LeafNode)(nil), gotRightLeaf.NextLeaf)
	})
}
