package main

import "slices"

var _ Node = (*LeafNode)(nil)

type LeafNode struct {
	*node
	NextLeaf *LeafNode
	PrevLeaf *LeafNode
	values   [][]any
}

func NewLeafNode(
	node *node,
	values [][]any,
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
	found := make([]any, 0, 5)
	for i, nodeKey := range l.Keys {
		if nodeKey == searchKey {
			found = append(found, l.values[i]...)
		}
	}

	return found
}

func (l *LeafNode) Insert(insertKey int64, insertValue any) {
	if len(l.Keys) == 0 {
		l.Keys = append(l.Keys, insertKey)
		l.values = append(l.values, []any{insertValue})
		return
	}

	// если больше самого правого, то вставляем в конец
	lastItemIndex := len(l.values) - 1
	if insertKey > l.Keys[lastItemIndex] {
		l.Keys = append(l.Keys, insertKey)
		l.values = append(l.values, []any{insertValue})
	} else if insertKey == l.Keys[lastItemIndex] {
		l.values[lastItemIndex] = append(l.values[lastItemIndex], insertValue)
	} else {
		// вставка в отсортированный слайс
		for i, nodeKey := range l.Keys {
			if insertKey < nodeKey {
				l.Keys = insertBefore(l.Keys, insertKey, i)
				l.values = insertBefore(l.values, []any{insertValue}, i)
				break
			} else if insertKey == nodeKey {
				l.values[i] = append(l.values[i], insertValue)
				break
			}
		}
	}

	if !l.isOverflow() {
		return
	}

	// переполненный лист делится на два листа
	half := len(l.values) / 2
	newDivider := l.Keys[half]

	leftHalfKeys := make([]int64, half)
	rightHalfKeys := make([]int64, len(l.Keys)-half)
	leftHalfValues := make([][]any, half)
	rightHalfValues := make([][]any, len(l.Keys)-half)

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

func (l *LeafNode) Delete(deleteKey int64) {
	var (
		found    bool
		deleteAt int
	)

	for i, nodeKey := range l.Keys {
		if nodeKey == deleteKey {
			found = true
			deleteAt = i
			break
		}
	}

	if !found {
		return
	}

	// если можем удалить из текущего листа,
	// то удаляем
	canRemove := !isUnderflow(int(l.maxKeys), len(l.Keys)) || l.isRoot()
	if canRemove {
		l.Keys = remove(l.Keys, deleteAt)
		l.values = remove(l.values, deleteAt)
		return
	}

	// если не смогли удалить,
	// то ищем замену в соседних листах

	var (
		leftSibling,
		rightSibling,
		_,
		_ = getSiblings[*LeafNode](l.node)
	)

	canTakeFromLeftLeaf := leftSibling != nil &&
		!isUnderflow(int(leftSibling.maxKeys), len(leftSibling.Keys)-1)

	canTakeFromRightLeaf := rightSibling != nil &&
		!isUnderflow(int(rightSibling.maxKeys), len(rightSibling.Keys)-1)

	if canTakeFromLeftLeaf {
		// забираем крайний правый элемент
		// из соседнего левого узла
		// на место удаленного элемента

		lastItemIndex := len(leftSibling.Keys) - 1

		deletedKey := l.Keys[deleteAt]
		replaceKey := leftSibling.Keys[lastItemIndex]
		replaceValue := leftSibling.values[lastItemIndex]

		leftSibling.Keys = remove(leftSibling.Keys, lastItemIndex)
		leftSibling.values = remove(leftSibling.values, lastItemIndex)

		l.Keys[deleteAt] = replaceKey
		l.values[deleteAt] = replaceValue

		// проверка, что следует обновить разделитель
		// в родительском внутреннем узле
		if deleteAt == 0 {
			deletedDividerIndex := slices.Index(l.Parent.Keys, deletedKey)
			if deletedDividerIndex != -1 {
				l.Parent.Keys[deletedDividerIndex] = replaceKey
			}
		}
	} else if canTakeFromRightLeaf {
		// забираем крайний левый элемент
		// из соседнего правого узла
		// на место удаленного элемента

		const (
			firstItemIndex = 0
		)

		deletedKey := l.Keys[deleteAt]
		replaceKey := rightSibling.Keys[firstItemIndex]
		replaceValue := rightSibling.values[firstItemIndex]

		rightSibling.Keys = remove(rightSibling.Keys, firstItemIndex)
		rightSibling.values = remove(rightSibling.values, firstItemIndex)

		l.Keys[deleteAt] = replaceKey
		l.values[deleteAt] = replaceValue

		// проверка, что следует обновить разделитель
		// в родительском внутреннем узле
		deletedDividerIndex := slices.Index(rightSibling.Parent.Keys, deletedKey)

		if deletedDividerIndex != -1 {
			newDivider := rightSibling.Keys[firstItemIndex]
			rightSibling.Parent.Keys[deletedDividerIndex] = newDivider
		}
	} else {
		// если не смогли взять элемент для замены
		// ни из левого,
		// ни из правого узлов,
		// то производим слияние

		if rightSibling != nil {
			// если можем смержиться с правым узлом

			l.Keys = append(l.Keys, rightSibling.Keys...)
			l.values = append(l.values, rightSibling.values...)

			var (
				dividerForRemoveIndex int
				childForRemoveIndex   int
			)

			for i, divider := range rightSibling.Parent.Keys {
				firstNextLeafKey := rightSibling.Keys[0]
				if divider >= firstNextLeafKey {
					if divider > firstNextLeafKey {
						childForRemoveIndex = 0
					} else if divider == firstNextLeafKey {
						childForRemoveIndex = i + 1
					}

					dividerForRemoveIndex = i
					break
				}
			}

			rightSibling.Parent.Children = remove(rightSibling.Parent.Children, childForRemoveIndex)
			rightSibling.Parent.removeDivider(dividerForRemoveIndex)

			// обновляем указатели на соседние листы
			l.NextLeaf = rightSibling.NextLeaf
			if l.NextLeaf != nil {
				l.NextLeaf.PrevLeaf = l
			}
		} else {
			// если правого узла нет,
			// то мержимся с левым узлом

			leftSibling.Keys = append(leftSibling.Keys, l.Keys...)
			leftSibling.values = append(l.PrevLeaf.values, l.values...)

			preLastChildIndex := len(l.Parent.Children)-2
			l.Parent.Children = remove(l.Parent.Children, preLastChildIndex)

			dividerForRemoveIndex := len(l.Parent.Keys) - 1
			l.Parent.removeDivider(dividerForRemoveIndex)

			// не нужно менять указатели,
			// они уже правильные
		}
	}
}
