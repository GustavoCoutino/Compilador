package stack

import "errors"

type Stack[T any] struct {
	Items []T
}

func New[T any]() *Stack[T] {
	items := []T{}
	return &Stack[T]{
		Items: items,
	}
}

func (s *Stack[T]) Push(val T) {
	s.Items = append(s.Items, val)
}

func (s *Stack[T]) Pop() (T, error) {
	if len(s.Items) == 0 {
		var err T
		return err, errors.New("stack underflow: cannot pop from empty stack")
	} else {
		poppedItem := s.Items[len(s.Items)-1]
		s.Items = s.Items[:len(s.Items)-1]
		return poppedItem, nil
	}
}

func (s *Stack[T]) Top() (T, error) {
	if len(s.Items) == 0 {
		var zero T
		return zero, errors.New("error: stack has no elements")
	}
	return s.Items[len(s.Items)-1], nil
}

func (s *Stack[T]) Empty() bool {
	return len(s.Items) == 0
}