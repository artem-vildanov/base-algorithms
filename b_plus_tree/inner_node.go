package main

import (
	"fmt"
	"log"
)

type InnerNode struct {
	Keys     []int64
	Children []any // *Node | *Leaf
	Parent   *InnerNode
	setRoot func(n Node)
	maxKeys int8
}

func NewInnerNode(
	keys []int64, 
	children []any, 
	parent *InnerNode,
	setRoot func(n Node),
	maxKeys int8,
) *InnerNode {
	return &InnerNode{
		Keys:     keys,
		Children: children,
		Parent:   parent,
		setRoot: setRoot,
		maxKeys: maxKeys,
	}
}

func (n *InnerNode) Find(searchKey int64) (any, error) {
	leaf, err := n.findLeaf(searchKey)
	if err != nil {
		return nil, fmt.Errorf("Node.findLeaf: %w", err)
	}

	value, err := leaf.Find(searchKey)
	if err != nil {
		return nil, fmt.Errorf("LeafNode.Find: %w", err)
	}

	return value, nil
}

func (n *InnerNode) findLeaf(searchKey int64) (*LeafNode, error) {
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
		return castedNode, nil
	default:
		return nil, fmt.Errorf("unexpected child type: %T", castedNode)
	}
}

func (n *InnerNode) Insert(insertKey int64, insertValue any) {
	leaf, err := n.findLeaf(insertKey)
	if err != nil {
		log.Printf("Node.findLeaf: %s", err.Error())
		return
	}

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
				n.Children = insertAfter(n.Children, rightChild, i)
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
	leftHalf := len(n.Keys) / 2
	rightHalf := len(n.Keys) - leftHalf
	newDivider := n.Keys[rightHalf] // среднее значение

	/*
		если это корень - то создаем новый корень
	*/
	if n.isRoot() {
		n.Parent = NewInnerNode(
			nil, 
			nil, 
			nil, 
			n.setRoot, 
			n.maxKeys,
		)
		n.setRoot(n.Parent)
	}

	leftNode := NewInnerNode(
		n.Keys[:leftHalf],
		n.Children[:leftHalf+1], // len(children) = len(keys) + 1
		n.Parent,
		n.setRoot,
		n.maxKeys,
	)

	rightNode := NewInnerNode(
		/*
			среднее значение не попадает в новый узел,
			а добавляется в родителя
		*/
		n.Keys[rightHalf+1:],
		n.Children[rightHalf:],
		n.Parent,
		n.setRoot,
		n.maxKeys,
	)

	/*
		добавляем в родителя разделитель
	*/
	n.Parent.addDivider(newDivider, leftNode, rightNode)
}

func (n *InnerNode) isRoot() bool {
	return n.Parent == nil
}

func (n *InnerNode) isOverflow() bool {
	return len(n.Keys) > int(n.maxKeys)
}
