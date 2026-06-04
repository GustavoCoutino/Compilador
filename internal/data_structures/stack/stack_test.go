package stack

import (
	"testing"
)

func TestStack_Push(t *testing.T) {
	t.Run("pushing int to empty stack", func(t *testing.T) {
		st := New[int]()
		st.Push(1)
		got, _ := st.Top()
		want := 1
		assertCorrectValue(t, got, want)
	})

	t.Run("pushing element to stack with numerical values", func(t *testing.T) {
		st := New[int]()
		st.Push(1)
		st.Push(2)
		got, _ := st.Top()
		want := 2
		assertCorrectValue(t, got, want)
	})

	t.Run("pushing multiple ints preserves LIFO order", func(t *testing.T) {
		s := New[int]()
		s.Push(1)
		s.Push(2)
		s.Push(3)
		got, _ := s.Top()
		want := 3
		assertCorrectValue(t, got, want)
	})

	t.Run("pushing string to empty stack", func(t *testing.T) {
		st := New[string]()
		st.Push("World")
		got, _ := st.Top()
		want := "World"
		assertCorrectValue(t, got, want)
	})

	t.Run("pushing string to stack with strings", func(t *testing.T) {
		st := New[string]()
		st.Push("Hello")
		st.Push("World")
		got, _ := st.Top()
		want := "World"
		assertCorrectValue(t, got, want)
	})
}

func TestStack_Pop(t *testing.T) {
	t.Run("popping a stack with integers", func(t *testing.T) {
		st := New[int]()
		st.Push(1)
		item, err := st.Pop()
		if err != nil {
			t.Error()
		}
		if item != 1 {
			t.Errorf("was expecting %d, got: %d", 1, item)
		}
	})

	t.Run("popping an empty stack", func(t *testing.T) {
		st := New[int]()
		_, err := st.Pop()
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestStack_Top(t *testing.T) {
	t.Run("top on empty stack returns error", func(t *testing.T) {
		st := New[int]()
		_, err := st.Top()
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("top returns last element added", func(t *testing.T) {
		st := New[int]()
		st.Push(42)
		got, err := st.Top()
		want := 42
		if err != nil {
			t.Error("was not expecting error")
		}
		assertCorrectValue(t, got, want)
	})
}

func TestStack_Empty(t *testing.T) {
	t.Run("new stack is empty", func(t *testing.T) {
		st := New[int]()
		if !st.Empty() {
			t.Error("expected empty stack")
		}
	})

	t.Run("stack with element is not empty", func(t *testing.T) {
		st := New[int]()
		st.Push(1)
		if st.Empty() {
			t.Error("expected non-empty stack")
		}
	})
}

func assertCorrectValue[T comparable](t testing.TB, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}

}