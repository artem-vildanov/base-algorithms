package main

type Node interface {
	Find(searchKey int64) []any
	Insert(insertKey int64, insertValue any)
}

type node struct {
	Keys    []int64
	Parent  *InnerNode
	setRoot func(n Node)
	maxKeys int8
}

func (n *node) isRoot() bool {
	return n.Parent == nil
}

func (n *node) isOverflow() bool {
	return len(n.Keys) > int(n.maxKeys)
}

func NewNode(
	keys []int64,
	parent *InnerNode,
	setRoot func(Node),
	maxKeys int8,
) *node {
	return &node{
		Keys:    keys,
		Parent:  parent,
		setRoot: setRoot,
		maxKeys: maxKeys,
	}
}

type Tree struct {
	root Node
}

func NewTree(maxKeys int8) *Tree {
	t := new(Tree)
	t.root = NewLeafNode(
		NewNode(
			nil,
			nil,
			func(root Node) {
				t.root = root
			},
			maxKeys,
		),
		nil,
		nil,
		nil,
	)

	return t
}

func (t *Tree) Find(searchKey int64) []any {
	return t.root.Find(searchKey)
}

func (t *Tree) Insert(insertKey int64, insertValue any) {
	t.root.Insert(insertKey, insertValue)
}
