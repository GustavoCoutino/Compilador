package queue

import (
	"testing"
)

func TestQueue_Push(t *testing.T) {
	t.Run("pushing int to empty queue", func(t *testing.T) {
		q := New[int]()
		q.Push(1)
		got, _ := q.Front()
		want := 1
		assertCorrectValue(t, got, want)
	})

	t.Run("pushing multiple ints preserves FIFO order", func(t *testing.T) {
		q := New[int]()
		q.Push(1)
		q.Push(2)
		q.Push(3)
		got, _ := q.Front()
		want := 1
		assertCorrectValue(t, got, want)
	})

	t.Run("pushing string to empty queue", func(t *testing.T) {
		q := New[string]()
		q.Push("hello")
		got, _ := q.Front()
		want := "hello"
		assertCorrectValue(t, got, want)
	})

	t.Run("back returns last pushed element", func(t *testing.T) {
		q := New[int]()
		q.Push(10)
		q.Push(20)
		got, _ := q.Back()
		want := 20
		assertCorrectValue(t, got, want)
	})
}

func TestQueue_Pop(t *testing.T) {
	t.Run("popping a queue with one element", func(t *testing.T) {
		q := New[int]()
		q.Push(1)
		item, err := q.Pop()
		if err != nil {
			t.Errorf("was not expecting error, got: %v", err)
		}
		if item != 1 {
			t.Errorf("was expecting %d, got: %d", 1, item)
		}
	})

	t.Run("popping an empty queue returns error", func(t *testing.T) {
		q := New[int]()
		_, err := q.Pop()
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("pop removes front element (FIFO)", func(t *testing.T) {
		q := New[int]()
		q.Push(1)
		q.Push(2)
		q.Pop()
		got, _ := q.Front()
		want := 2
		assertCorrectValue(t, got, want)
	})
}

func TestQueue_Front(t *testing.T) {
	t.Run("front on empty queue returns error", func(t *testing.T) {
		q := New[int]()
		_, err := q.Front()
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("front returns first element added", func(t *testing.T) {
		q := New[int]()
		q.Push(42)
		q.Push(99)
		got, err := q.Front()
		want := 42
		if err != nil {
			t.Errorf("was not expecting error, got: %v", err)
		}
		assertCorrectValue(t, got, want)
	})
}

func TestQueue_Back(t *testing.T) {
	t.Run("back on empty queue returns error", func(t *testing.T) {
		q := New[int]()
		_, err := q.Back()
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("back returns last element added", func(t *testing.T) {
		q := New[int]()
		q.Push(1)
		q.Push(2)
		q.Push(3)
		got, err := q.Back()
		want := 3
		if err != nil {
			t.Errorf("was not expecting error, got: %v", err)
		}
		assertCorrectValue(t, got, want)
	})
}

func TestQueue_Empty(t *testing.T) {
	t.Run("new queue is empty", func(t *testing.T) {
		q := New[int]()
		if !q.Empty() {
			t.Error("expected empty queue")
		}
	})

	t.Run("queue with element is not empty", func(t *testing.T) {
		q := New[int]()
		q.Push(1)
		if q.Empty() {
			t.Error("expected non-empty queue")
		}
	})

	t.Run("queue is empty after popping last element", func(t *testing.T) {
		q := New[int]()
		q.Push(1)
		q.Pop()
		if !q.Empty() {
			t.Error("expected empty queue after popping last element")
		}
	})
}

func assertCorrectValue[T comparable](t testing.TB, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}