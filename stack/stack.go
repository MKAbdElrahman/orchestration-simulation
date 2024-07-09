package stack

import (
	"fmt"
	"slices"
)

/*
- last in first out LIFO
- push
- pop
- push O(1)
- pop O(1)
- search o(n)
- access o(n)
*/

var (
	ErrEmptyStack = fmt.Errorf("poping from an empty stack")
	ErrNotFound   = fmt.Errorf("consumed all the stack but value not found")
)

type stack[T comparable] struct {
	elements []T
}

func NewStack[T comparable]() *stack[T] {
	return &stack[T]{
		elements: make([]T, 0),
	}
}
func (s *stack[T]) Push(e T) {
	s.elements = append(s.elements, e)
}

func (s *stack[T]) Pop() (T, error) {

	if s.Size() == 0 {
		var zero T
		return zero, ErrEmptyStack
	}

	lastElement := s.elements[len(s.elements)-1]
	s.elements = slices.Delete(s.elements, len(s.elements)-1, len(s.elements))

	return lastElement, nil
}

func (s *stack[T]) Size() int {
	return len(s.elements)
}

func (s *stack[T]) IsEmpty() bool {
	return len(s.elements) == 0
}

func (s *stack[T]) Contains(v T) bool {
	return slices.Contains(s.elements, v)
}

func (s *stack[T]) Access(v T) error {
	for s.Size() != 0 {
		element, _ := s.Pop() // never an error bacuase of the iteration condition
		if element == v {
			return nil
		}
	}

	return ErrNotFound
}
