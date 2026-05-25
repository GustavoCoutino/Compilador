package queue

import (
	"errors"
)

type Queue[T any] struct {
	Items []T
}

func New[T any]() *Queue[T] {
	items := []T{}
	return &Queue[T]{
		Items: items,
	}
}

func (q *Queue[T]) Push(val T) {
	q.Items = append(q.Items, val)
}

func (q *Queue[T]) Pop() (T, error) {
	if len(q.Items) == 0 {
		var err T
		return err, errors.New("queue underflow: cannot pop from empty queue")
	} else {
		poppedItem := q.Items[0]
		q.Items = q.Items[1:]
		return poppedItem, nil
	}
}

func (q *Queue[T]) Front() (T, error) {
	if len(q.Items) == 0 {
		var zero T
		return zero, errors.New("error: queue has no elements")
	}
	return q.Items[0], nil
}

func (q *Queue[T]) Back() (T, error) {
	if len(q.Items) == 0 {
		var zero T
		return zero, errors.New("error: queue has no elements")
	}
	return q.Items[len(q.Items)-1], nil
}

func (q *Queue[T]) Empty() bool {
	return len(q.Items) == 0
}