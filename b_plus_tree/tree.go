package main

type Node interface {
	Find(searchKey int64) (any, error)
	Insert(insertKey int64, insertValue any)
}

type Tree struct {
	root Node
}

func NewTree(maxKeys int8) *Tree {
	t := new(Tree)
	t.root = NewLeafNode(
		nil,
		nil,
		nil,
		nil,
		func(root Node) {
			t.root = root
		},
		maxKeys,
	)

	return t
}

func (t *Tree) Find(searchKey int64) (any, error) {
	return t.root.Find(searchKey)
}

func (t *Tree) Insert(insertKey int64, insertValue any) {
	t.root.Insert(insertKey, insertValue)
}
