package main

type Node interface {
	Find(searchKey int64) []any
	Insert(insertKey int64, insertValue any)
	Delete(deleteKey int64)
}

type node struct {
	Keys    []int64
	Parent  *InnerNode
	setRoot func(n Node)
	maxKeys int8
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

func (n *node) isRoot() bool {
	return n.Parent == nil
}

func (n *node) isOverflow() bool {
	return len(n.Keys) > int(n.maxKeys)
}

func (n *node) isUnderflow() bool {
	return isUnderflow(int(n.maxKeys), len(n.Keys))
}

func getSiblings[T Node](n *node) (
	leftSibling T,
	rightSibling T,
	leftParentDividerIndex int,
	rightParentDividerIndex int,
) {
	var (
		firstNodeKey    = n.Keys[0]
		doesntHaveLeft  = firstNodeKey < n.Parent.Keys[0]
		doesntHaveRight = firstNodeKey > n.Parent.Keys[len(n.Parent.Keys)-1]
	)

	// находим левый и правый узлы
	if doesntHaveLeft {
		secondNodeIndex := 1
		rightSibling = n.Parent.Children[secondNodeIndex].(T)
		rightParentDividerIndex = 0
	} else if doesntHaveRight {
		preLastIndex := len(n.Parent.Children) - 2
		leftSibling = n.Parent.Children[preLastIndex].(T)
		leftParentDividerIndex = len(n.Parent.Keys) - 1
	} else {
		for i := 0; i < len(n.Parent.Keys); i++ {
			parentKey := n.Parent.Keys[i]

			if parentKey > firstNodeKey {
				rightSibling = n.Parent.Children[i+1].(T)
				leftSibling = n.Parent.Children[i-1].(T)

				leftParentDividerIndex = i - 1
				rightParentDividerIndex = i

				break
			}
		}
	}

	return
}

func isUnderflow(maxKeys, keysNum int) bool {
	// делим с округлением в большую сторону
	divedCeil := (maxKeys + 2 - 1) / 2
	return keysNum < divedCeil-1
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

func (t *Tree) Delete(deleteKey int64) {
	t.root.Delete(deleteKey)
}
