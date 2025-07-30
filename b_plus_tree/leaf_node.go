package main

import "fmt"

type KeyToValue struct {
	Key   int64
	Value any
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

type LeafNode struct {
	Values   []*KeyToValue
	NextLeaf *LeafNode
	PrevLeaf *LeafNode
	Parent   *InnerNode
	setRoot  func(n Node)
	maxKeys  int8
}

func NewLeafNode(
	parent *InnerNode,
	values []*KeyToValue,
	nextLeaf *LeafNode,
	prevLeaf *LeafNode,
	setRoot func(n Node),
	maxKeys int8,
) *LeafNode {
	return &LeafNode{
		Parent:   parent,
		Values:   values,
		NextLeaf: nextLeaf,
		PrevLeaf: prevLeaf,
		setRoot:  setRoot,
		maxKeys:  maxKeys,
	}
}

func (l *LeafNode) Find(searchKey int64) (any, error) {
	for _, keyToValue := range l.Values {
		if keyToValue.Key == searchKey {
			return keyToValue.Value, nil
		}
	}

	return nil, fmt.Errorf("not found by key: %d", searchKey)
}

func (l *LeafNode) Insert(insertKey int64, insertValue any) {
	newKeyToVal := &KeyToValue{
		Key:   insertKey,
		Value: insertValue,
	}

	if len(l.Values) == 0 {
		l.Values = append(l.Values, newKeyToVal)
		return
	}

	// если больше самого правого, то вставляем в конец
	if insertKey >= l.Values[len(l.Values)-1].Key {
		l.Values = append(l.Values, newKeyToVal)
	} else {
		for i, keyToVal := range l.Values {
			if insertKey < keyToVal.Key {
				// вставляем в отсортированный слайс
				l.Values = insertAfter(l.Values, newKeyToVal, i)
				break
			}
		}
	}

	if !l.isOverflow() {
		return
	}

	// если не хватило места для вставки

	// переполненный лист делится на два листа
	leftHalf := len(l.Values) / 2
	rightHalf := len(l.Values) - leftHalf

	if l.isRoot() {
		l.Parent = NewInnerNode(
			nil,
			nil,
			nil,
			l.setRoot,
			l.maxKeys,
		)
		l.setRoot(l.Parent)
	}

	leftLeaf := NewLeafNode(
		l.Parent,
		l.Values[:leftHalf],
		nil,
		l.PrevLeaf,
		l.setRoot,
		l.maxKeys,
	)

	rightLeaf := NewLeafNode(
		l.Parent,
		l.Values[rightHalf:],
		l.NextLeaf,
		leftLeaf,
		l.setRoot,
		l.maxKeys,
	)

	leftLeaf.NextLeaf = rightLeaf

	// обновляем указатели в соседних листах
	l.PrevLeaf.NextLeaf = leftLeaf
	l.NextLeaf.PrevLeaf = rightLeaf

	// поднимаем в родителя новый разделитель
	l.Parent.addDivider(
		rightLeaf.Values[0].Key,
		leftLeaf,
		rightLeaf,
	)
}

func (l *LeafNode) isOverflow() bool {
	return len(l.Values) > int(l.maxKeys)
}

func (l *LeafNode) isRoot() bool {
	return l.Parent == nil
}
