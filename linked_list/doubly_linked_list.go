package linkedlist

import (
	"errors"
)

type doublyNode[T any] struct {
	data T
	next *doublyNode[T]
	prev *doublyNode[T]
}

type DoublyLinkedList[T any] struct {
	head   *doublyNode[T]
	tail   *doublyNode[T]
	length int
}

func NewDoublyLinkedList[T any]() *DoublyLinkedList[T] {
	return &DoublyLinkedList[T]{}
}

func (l *DoublyLinkedList[T]) Append(data T) {
	newNode := &doublyNode[T]{data: data}
	if l.tail == nil {
		l.head = newNode
		l.tail = newNode
	} else {
		l.tail.next = newNode
		newNode.prev = l.tail
		l.tail = newNode
	}
	l.length++
}

func (l *DoublyLinkedList[T]) Prepend(data T) {
	newNode := &doublyNode[T]{data: data}
	if l.head == nil {
		l.head = newNode
		l.tail = newNode
	} else {
		newNode.next = l.head
		l.head.prev = newNode
		l.head = newNode
	}
	l.length++
}

func (l *DoublyLinkedList[T]) Pop() (T, error) {
	if l.length == 0 {
		var zero T
		return zero, errors.New("list is empty")
	}
	data := l.head.data
	l.head = l.head.next
	if l.head != nil {
		l.head.prev = nil
	} else {
		l.tail = nil
	}
	l.length--
	return data, nil
}

func (l *DoublyLinkedList[T]) PopLast() (T, error) {
	if l.length == 0 {
		var zero T
		return zero, errors.New("list is empty")
	}
	data := l.tail.data
	l.tail = l.tail.prev
	if l.tail != nil {
		l.tail.next = nil
	} else {
		l.head = nil
	}
	l.length--
	return data, nil
}

func (l *DoublyLinkedList[T]) Insert(index int, data T) error {
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

	newNode := &doublyNode[T]{data: data}
	current := l.head
	for range index {
		current = current.next
	}

	newNode.prev = current.prev
	newNode.next = current
	if current.prev != nil {
		current.prev.next = newNode
	}
	current.prev = newNode

	l.length++
	return nil
}

func (l *DoublyLinkedList[T]) Remove(index int) (T, error) {
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
	for range index {
		current = current.next
	}

	current.prev.next = current.next
	current.next.prev = current.prev

	l.length--
	return current.data, nil
}

func (l *DoublyLinkedList[T]) Get(index int) (T, error) {
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

func (l *DoublyLinkedList[T]) Length() int {
	return l.length
}

