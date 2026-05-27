package stack

import (
	"testing"
)

type ustack[T comparable] interface {
	Push(T)
	Pop() (T, error)
	Peek() (T, error)
	IsEmpty() bool
	Length() int
	Clear()
}

var factories = map[string]func() ustack[int]{
	"UniqueStack":     func() ustack[int] { return NewUniqueStack[int]() },
	"StackArray":      func() ustack[int] { return NewStackArray[int]() },
	"StackLinkedList": func() ustack[int] { return NewStackLinkedList[int]() },
}

func mustPop(t *testing.T, s ustack[int]) int {
	t.Helper()
	v, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop returned error: %v", err)
	}
	return v
}

func TestPushPopPeekAndLength(t *testing.T) {
	for name, newStack := range factories {
		t.Run(name, func(t *testing.T) {
			s := newStack()

			values := []int{1, 2, 3}
			for _, v := range values {
				s.Push(v)
			}

			if got, want := s.Length(), len(values); got != want {
				t.Fatalf("Length = %d, want %d", got, want)
			}

			if top, _ := s.Peek(); top != 3 {
				t.Fatalf("Peek = %d, want 3", top)
			}

			// Pop should return 3-2-1 (LIFO)
			for i := len(values) - 1; i >= 0; i-- {
				if got := mustPop(t, s); got != values[i] {
					t.Fatalf("Pop = %d, want %d", got, values[i])
				}
			}

			if !s.IsEmpty() {
				t.Fatal("stack should be empty after popping everything")
			}
		})
	}
}

func TestDuplicatePromotion(t *testing.T) {
	for name, newStack := range factories {
		if name != "UniqueStack" {
			continue
		}

		t.Run(name, func(t *testing.T) {
			s := newStack()

			// Push a sequence containing duplicates.
			input := []int{1, 2, 3, 2, 1, 4}
			for _, v := range input {
				s.Push(v)
			}

			// Expected pop order after duplicate-to-top behaviour:
			want := []int{4, 1, 2, 3}
			for _, v := range want {
				if got := mustPop(t, s); got != v {
					t.Fatalf("Pop = %d, want %d", got, v)
				}
			}

			if !s.IsEmpty() {
				t.Fatal("stack should be empty after popping everything")
			}
		})
	}
}

func TestClearAndErrorCases(t *testing.T) {
	for name, newStack := range factories {
		t.Run(name, func(t *testing.T) {
			s := newStack()
			s.Push(42)
			s.Push(99)
			s.Clear()

			if !s.IsEmpty() || s.Length() != 0 {
				t.Fatal("Clear() did not reset the stack")
			}

			if _, err := s.Pop(); err == nil {
				t.Fatal("expected error when popping from empty stack, got nil")
			}

			if _, err := s.Peek(); err == nil {
				t.Fatal("expected error when peeking into empty stack, got nil")
			}
		})
	}
}
