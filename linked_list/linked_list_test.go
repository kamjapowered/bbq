package linkedlist

import (
	"slices"
	"testing"
)

type list[T any] interface {
	Append(data T)
	Prepend(data T)
	Pop() (T, error)
	PopLast() (T, error)
	Insert(index int, data T) error
	Remove(index int) (T, error)
	Get(index int) (T, error)
	Length() int
}

func testAppend(t *testing.T, list list[int], expected []int) {
	for _, v := range expected {
		list.Append(v)
	}
	for i, v := range expected {
		if got, _ := list.Get(i); got != v {
			t.Errorf("Append: Expected %d at index %d, got %d", v, i, got)
		}
	}
}

func testPrepend(t *testing.T, list list[int], expected []int) {
	for _, v := range slices.Backward(expected) {
		list.Prepend(v)
	}
	for i, v := range expected {
		if got, _ := list.Get(i); got != v {
			t.Errorf("Prepend: Expected %d at index %d, got %d", v, i, got)
		}
	}
}

func testPop(t *testing.T, list list[int], expected []int) {
	for _, v := range expected {
		got, err := list.Pop()
		if err != nil {
			t.Errorf("Pop: Unexpected error: %v", err)
		}
		if got != v {
			t.Errorf("Pop: Expected %d, got %d", v, got)
		}
	}
	if list.Length() != 0 {
		t.Errorf("Pop: Expected empty list, got length %d", list.Length())
	}
	if _, err := list.Pop(); err == nil {
		t.Errorf("Pop: Expected error on empty list, got nil")
	}
}

func testPopLast(t *testing.T, list list[int], expected []int) {
	for _, v := range expected {
		got, err := list.PopLast()
		if err != nil {
			t.Errorf("PopLast: Unexpected error: %v", err)
		}
		if got != v {
			t.Errorf("PopLast: Expected %d, got %d", v, got)
		}
	}
	if list.Length() != 0 {
		t.Errorf("PopLast: Expected empty list, got length %d", list.Length())
	}
	if _, err := list.PopLast(); err == nil {
		t.Errorf("PopLast: Expected error on empty list, got nil")
	}
}

func testInsert(t *testing.T, list list[int], initial []int, index int, value int, expected []int) {
	for _, v := range initial {
		list.Append(v)
	}
	err := list.Insert(index, value)
	if err != nil {
		t.Errorf("Insert: Unexpected error: %v", err)
	}
	for i, v := range expected {
		if got, _ := list.Get(i); got != v {
			t.Errorf("Insert: Expected %d at index %d, got %d", v, i, got)
		}
	}
}

func testRemove(t *testing.T, list list[int], initial []int, index int, expectedValue int, expectedList []int) {
	for _, v := range initial {
		list.Append(v)
	}
	got, err := list.Remove(index)
	if err != nil {
		t.Errorf("Remove: Unexpected error: %v", err)
	}
	if got != expectedValue {
		t.Errorf("Remove: Expected removed value %d, got %d", expectedValue, got)
	}
	for i, v := range expectedList {
		if val, _ := list.Get(i); val != v {
			t.Errorf("Remove: Expected %d at index %d, got %d", v, i, val)
		}
	}
}

func TestMain(t *testing.T) {
	factories := map[string]func() list[int]{
		"LinkedList":       func() list[int] { return NewLinkedList[int]() },
		"DoublyLinkedList": func() list[int] { return NewDoublyLinkedList[int]() },
	}

	expected := []int{1, 2, 3}

	for name, factory := range factories {
		t.Run(name+"/append", func(t *testing.T) {
			list := factory()
			testAppend(t, list, expected)
		})

		t.Run(name+"/prepend", func(t *testing.T) {
			list := factory()
			testPrepend(t, list, expected)
		})

		t.Run(name+"/pop", func(t *testing.T) {
			list := factory()
			for _, v := range expected {
				list.Append(v)
			}
			testPop(t, list, expected)
		})

		t.Run(name+"/popLast", func(t *testing.T) {
			list := factory()
			for _, v := range expected {
				list.Append(v)
			}
			reversed := slices.Clone(expected)
			slices.Reverse(reversed)
			testPopLast(t, list, reversed)
		})

		t.Run(name+"/insert", func(t *testing.T) {
			list := factory()
			initial := []int{1, 3}
			expectedList := []int{1, 2, 3}
			testInsert(t, list, initial, 1, 2, expectedList)
		})

		t.Run(name+"/remove", func(t *testing.T) {
			list := factory()
			initial := []int{1, 2, 3}
			expectedList := []int{1, 3}
			testRemove(t, list, initial, 1, 2, expectedList)
		})
	}
}

