package main

var _ Node = (*LeafNode)(nil)

type LeafNode struct {
	*node
	NextLeaf *LeafNode
	PrevLeaf *LeafNode
	values   []any
}

func NewLeafNode(
	node *node,
	values []any,
	nextLeaf *LeafNode,
	prevLeaf *LeafNode,
) *LeafNode {
	return &LeafNode{
		node:     node,
		values:   values,
		NextLeaf: nextLeaf,
		PrevLeaf: prevLeaf,
	}
}

func (l *LeafNode) Find(searchKey int64) []any {
	found := make([]any, 5)
	for i, nodeKey := range l.Keys {
		if nodeKey == searchKey {
			found = append(found, l.values[i])
		}
	}

	return found
}

func (l *LeafNode) Insert(insertKey int64, insertValue any) {
	if len(l.Keys) == 0 {
		l.Keys = append(l.Keys, insertKey)
		l.values = append(l.values, insertValue)
		return
	}

	/*
		если больше самого правого, то вставляем в конец
	*/
	if insertKey >= l.Keys[len(l.values)-1] {
		l.Keys = append(l.Keys, insertKey)
		l.values = append(l.values, insertValue)
	} else {
		/*
			вставка в отсортированный слайс
		*/
		for i, nodeKey := range l.Keys {
			if insertKey < nodeKey {
				l.Keys = insertAfter(l.Keys, insertKey, i)
				l.values = insertAfter(l.values, insertValue, i)
				break
			}
		}
	}

	if !l.isOverflow() {
		return
	}

	/*
		если не хватило места для вставки
	*/

	/*
		переполненный лист делится на два листа
	*/
	half := len(l.values) / 2
	newDivider := l.Keys[half]

	leftHalfKeys := make([]int64, half)
	rightHalfKeys := make([]int64, len(l.Keys)-half)
	leftHalfValues := make([]any, half)
	rightHalfValues := make([]any, len(l.Keys)-half)

	copy(leftHalfKeys, l.Keys[:half])
	copy(rightHalfKeys, l.Keys[half:])
	copy(leftHalfValues, l.values[:half])
	copy(rightHalfValues, l.values[half:])

	if l.isRoot() {
		l.Parent = NewInnerNode(
			NewNode(
				nil,
				nil,
				l.setRoot,
				l.maxKeys,
			),
			nil,
		)
		l.setRoot(l.Parent)
	}

	leftLeaf := NewLeafNode(
		NewNode(
			leftHalfKeys,
			l.Parent,
			l.setRoot,
			l.maxKeys,
		),
		leftHalfValues,
		nil,
		l.PrevLeaf,
	)

	rightLeaf := NewLeafNode(
		NewNode(
			rightHalfKeys,
			l.Parent,
			l.setRoot,
			l.maxKeys,
		),
		rightHalfValues,
		l.NextLeaf,
		leftLeaf,
	)

	leftLeaf.NextLeaf = rightLeaf

	// обновляем указатели в соседних листах
	if l.PrevLeaf != nil {
		l.PrevLeaf.NextLeaf = leftLeaf
	}
	if l.NextLeaf != nil {
		l.NextLeaf.PrevLeaf = rightLeaf
	}

	// поднимаем в родителя новый разделитель
	l.Parent.addDivider(
		newDivider,
		leftLeaf,
		rightLeaf,
	)
}
