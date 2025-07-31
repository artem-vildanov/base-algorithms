package main

var _ Node = (*InnerNode)(nil)

type InnerNode struct {
	*node
	Children []any // *Node | *Leaf
}

func NewInnerNode(
	node *node,
	children []any,
) *InnerNode {
	return &InnerNode{
		node:     node,
		Children: children,
	}
}

func (n *InnerNode) Find(searchKey int64) []any {
	leaf := n.findLeaf(searchKey)
	value := leaf.Find(searchKey)
	return value
}

func (n *InnerNode) findLeaf(searchKey int64) *LeafNode {
	var node any

	/*
		если искомый ключ больше или равен
		чем крайний правый ключ, то возвращаем последний узел
	*/
	lastDivider := n.Keys[len(n.Keys)-1]
	if searchKey >= lastDivider {
		node = n.Children[len(n.Children)-1]
	} else {
		for i, divider := range n.Keys {
			if divider < searchKey {
				continue
			}

			node = n.Children[i]
			break
		}
	}

	switch castedNode := node.(type) {
	case *InnerNode:
		return castedNode.findLeaf(searchKey)
	case *LeafNode:
		return castedNode
	default:
		// не попадем сюда
		return nil
	}
}

func (n *InnerNode) Insert(insertKey int64, insertValue any) {
	leaf := n.findLeaf(insertKey)
	leaf.Insert(insertKey, insertValue)
}

func (n *InnerNode) addDivider(divider int64, leftChild any, rightChild any) {
	if len(n.Keys) == 0 {
		n.Keys = append(n.Keys, divider)
		n.Children = append(n.Children, leftChild, rightChild)
		return
	}

	/*
		если добавляем новый максимальный ключ
	*/
	if divider >= n.Keys[len(n.Keys)-1] {
		/*
			вставляем ключ в конец
			заменяем последнего ребенка
			вставляем в конец нового
		*/
		n.Keys = append(n.Keys, divider)
		n.Children[len(n.Children)-1] = leftChild
		n.Children = append(n.Children, rightChild)
	} else {
		/*
			вставка в сортированный слайс
		*/
		for i, nodeKey := range n.Keys {
			if divider < nodeKey {
				/*
					вставляем ключ
					замена левого ребенка
					вставляем правого ребенка
				*/
				n.Keys = insertAfter(n.Keys, divider, i)
				n.Children[i] = leftChild
				n.Children = insertAfter(n.Children, rightChild, i+1)
				break
			}
		}
	}

	if !n.isOverflow() {
		return
	}
	/*
		нет места в узле
	*/

	/*
		делим узел на два
	*/
	half := len(n.Keys) / 2
	newDivider := n.Keys[half]
	leftHalfKeys := make([]int64, half)
	rightHalfKeys := make([]int64, half)

	copy(leftHalfKeys, n.Keys[:half])
	copy(rightHalfKeys, n.Keys[half+1:])

	/*
		если это корень - то создаем новый корень
	*/
	if n.isRoot() {
		n.Parent = NewInnerNode(
			NewNode(
				nil,
				nil,
				n.setRoot,
				n.maxKeys,
			),
			nil,
		)
		n.setRoot(n.Parent)
	}

	leftNode := NewInnerNode(
		NewNode(
			leftHalfKeys,
			n.Parent,
			n.setRoot,
			n.maxKeys,
		),
		n.Children[:half+1], // len(children) = len(keys) + 1
	)

	rightNode := NewInnerNode(
		NewNode(
			/*
				средний ключ не попадает в новый узел,
				а добавляется в родителя
			*/
			rightHalfKeys,
			n.Parent,
			n.setRoot,
			n.maxKeys,
		),
		n.Children[half+1:],
	)

	/*
		добавляем в родителя разделитель
	*/
	n.Parent.addDivider(newDivider, leftNode, rightNode)
}