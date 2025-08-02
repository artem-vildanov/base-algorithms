package main

var _ Node = (*InnerNode)(nil)

type InnerNode struct {
	*node
	Children []Node // *Node | *Leaf
}

func NewInnerNode(
	node *node,
	children []Node,
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

func (n *InnerNode) Insert(insertKey int64, insertValue any) {
	leaf := n.findLeaf(insertKey)
	leaf.Insert(insertKey, insertValue)
}

func (n *InnerNode) Delete(deleteKey int64) {
	leaf := n.findLeaf(deleteKey)
	leaf.Delete(deleteKey)
}

func (n *InnerNode) findLeaf(searchKey int64) *LeafNode {
	var node Node

	// если искомый ключ больше или равен
	// чем крайний правый ключ, то возвращаем последний узел
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

func (n *InnerNode) addDivider(divider int64, leftChild, rightChild Node) {
	if len(n.Keys) == 0 {
		n.Keys = append(n.Keys, divider)
		n.Children = append(n.Children, leftChild, rightChild)
		return
	}

	// если добавляем новый максимальный ключ
	if divider >= n.Keys[len(n.Keys)-1] {
		// вставляем ключ в конец
		// заменяем последнего ребенка
		// вставляем в конец нового
		n.Keys = append(n.Keys, divider)
		n.Children[len(n.Children)-1] = leftChild
		n.Children = append(n.Children, rightChild)
	} else {
		// вставка в сортированный слайс
		for i, nodeKey := range n.Keys {
			if divider < nodeKey {
				// вставляем ключ
				// замена левого ребенка
				// вставляем правого ребенка
				n.Keys = insertBefore(n.Keys, divider, i)
				n.Children[i] = leftChild
				n.Children = insertBefore(n.Children, rightChild, i+1)
				break
			}
		}
	}

	if !n.isOverflow() {
		return
	}

	// делим узел на два

	half := len(n.Keys) / 2
	newDivider := n.Keys[half]
	leftHalfKeys := make([]int64, half)
	rightHalfKeys := make([]int64, half)

	copy(leftHalfKeys, n.Keys[:half])
	copy(rightHalfKeys, n.Keys[half+1:])

	// если это корень - то создаем новый корень
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
			// средний ключ не попадает в новый узел,
			// а добавляется в родителя
			rightHalfKeys,
			n.Parent,
			n.setRoot,
			n.maxKeys,
		),
		n.Children[half+1:],
	)

	// добавляем в родителя разделитель
	n.Parent.addDivider(newDivider, leftNode, rightNode)
}

func (n *InnerNode) removeDivider(dividerIndex int) {
	n.Keys = remove(n.Keys, dividerIndex)

	if !n.isUnderflow() || n.isRoot() {
		return
	}

	// если при удалении разделителя
	// в узле остается количество разделителей
	// меньшее порогового значения,
	// то производим ребалансировку

	var (
		leftSibling,
		rightSibling,
		leftParentDividerIndex,
		rightParentDividerIndex = getSiblings[*InnerNode](n.node)
	)

	// можем забрать ключ из узла слева
	canTakeFromLeft := leftSibling != nil &&
		!isUnderflow(int(leftSibling.maxKeys), len(leftSibling.Keys)-1)

	// можем забрать ключ из узла справа
	canTakeFromRight := rightSibling != nil &&
		!isUnderflow(int(rightSibling.maxKeys), len(rightSibling.Keys)-1)

	if canTakeFromLeft {
		// удаляем разделитель
		n.Keys = remove(n.Keys, dividerIndex)

		// забираем родительский разделитель
		n.Keys = append(
			[]int64{n.Parent.Keys[leftParentDividerIndex]},
			n.Keys...,
		)

		// на место родительского разделителя ставим
		// крайний правый ключ из левого узла
		lastLeftNodeKeyIndex := len(leftSibling.Keys) - 1
		n.Parent.Keys[leftParentDividerIndex] = leftSibling.Keys[lastLeftNodeKeyIndex]
		leftSibling.Keys = remove(leftSibling.Keys, lastLeftNodeKeyIndex)

		// забираем крайний правый указатель
		// левого узла
		lastLeftNodeChild := leftSibling.Children[len(leftSibling.Children)-1]
		n.Children = append([]Node{lastLeftNodeChild}, n.Children...)
	} else if canTakeFromRight {
		// удаляем разделитель
		n.Keys = remove(n.Keys, dividerIndex)

		// забираем родительский разделитель
		n.Keys = append(n.Keys, n.Parent.Keys[rightParentDividerIndex])

		// на место родительского разделителя ставим
		// крайний левый ключ из правого узла
		n.Parent.Keys[rightParentDividerIndex] = rightSibling.Keys[0]
		rightSibling.Keys = remove(rightSibling.Keys, 0)

		// забираем крайний левый указатель
		// правого узла
		firstRightNodeChild := rightSibling.Children[0]
		n.Children = append(n.Children, firstRightNodeChild)
	} else {
		// если не смогли взять элемент для замены
		// ни из левого,
		// ни из правого узлов,
		// то производим слияние

		if rightSibling != nil { // если можем смержиться с правым узлом
			n.Keys = append(
				n.Keys,
				append(
					[]int64{n.Parent.Keys[rightParentDividerIndex]},
					rightSibling.Keys...,
				)...,
			)
			n.Children = append(n.Children, rightSibling.Children...)

			// удаляем указатель на правый узел
			n.Parent.Children = remove(n.Parent.Children, rightParentDividerIndex+1)
			n.Parent.removeDivider(rightParentDividerIndex)
		} else { // если можем смержиться с левым узлом
			n.Keys = append(
				leftSibling.Keys,
				append(
					[]int64{n.Parent.Keys[leftParentDividerIndex]},
					n.Keys...,
				)...,
			)
			n.Children = append(leftSibling.Children, n.Children...)

			// удаляем указатель на левый узел
			n.Parent.Children = remove(n.Parent.Children, leftParentDividerIndex)
			n.Parent.removeDivider(leftParentDividerIndex)
		}
	}
}
