package linkedlist

import (
	"errors"
)

type node[T any] struct {
	data T
	next *node[T]
}

type LinkedList[T any] struct {
	head   *node[T]
	tail   *node[T]
	length int
}

func NewLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{}
}

func (l *LinkedList[T]) Append(data T) {
	newNode := &node[T]{data: data}

	if l.head == nil {
		l.head = newNode
		l.tail = newNode
	} else {
		l.tail.next = newNode
		l.tail = newNode
	}

	l.length++
}

func (l *LinkedList[T]) Prepend(data T) {
	newNode := &node[T]{data: data}

	if l.head == nil {
		l.head = newNode
		l.tail = newNode
	} else {
		newNode.next = l.head
		l.head = newNode
	}

	l.length++
}

func (l *LinkedList[T]) Pop() (T, error) {
	if l.length == 0 {
		var zero T
		return zero, errors.New("list is empty")
	}

	data := l.head.data
	l.head = l.head.next

	if l.head == nil {
		l.tail = nil
	}

	l.length--
	return data, nil
}

func (l *LinkedList[T]) PopLast() (T, error) {
	if l.length == 0 {
		var zero T
		return zero, errors.New("list is empty")
	}

	data := l.tail.data

	if l.head == l.tail {
		l.head = nil
		l.tail = nil
	} else {
		current := l.head
		for current.next != l.tail {
			current = current.next
		}
		current.next = nil
		l.tail = current
	}

	l.length--
	return data, nil
}

func (l *LinkedList[T]) Insert(index int, data T) error {
	if index < 0 || index > l.length {
		return errors.New("index out of bounds")
	}

	if index == 0 {
		l.Prepend(data)
		return nil
	}

	if index == l.length {
		l.Append(data)
		return nil
	}

	newNode := &node[T]{data: data}

	current := l.head
	for range index - 1 {
		current = current.next
	}
	newNode.next = current.next
	current.next = newNode

	l.length++

	return nil
}

func (l *LinkedList[T]) Remove(index int) (T, error) {
	if index < 0 || index >= l.length {
		var zero T
		return zero, errors.New("index out of bounds")
	}

	if index == 0 {
		return l.Pop()
	}

	if index == l.length-1 {
		return l.PopLast()
	}

	current := l.head
	for range index - 1 {
		current = current.next
	}
	data := current.next.data
	current.next = current.next.next

	l.length--
	return data, nil
}

func (l *LinkedList[T]) Get(index int) (T, error) {
	if index < 0 || index >= l.length {
		var zero T
		return zero, errors.New("index out of bounds")
	}

	current := l.head
	for range index {
		current = current.next
	}
	return current.data, nil
}

func (l *LinkedList[T]) Length() int {
	return l.length
}
