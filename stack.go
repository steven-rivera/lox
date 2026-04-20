package main

type Stack[T any] []T

func (s *Stack[T]) Push(v T) {
    *s = append(*s, v)
}

func (s *Stack[T]) Pop() (T, bool) {
    if s.Size() == 0 {
        var zero T
        return zero, false
    }
    idx := s.Size() - 1
    v := (*s)[idx]
    *s = (*s)[:idx]
    return v, true
}

func (s *Stack[T]) Peek() T {
	return (*s)[s.Size()-1]
}

func (s *Stack[T]) Get(index int) T {
	return (*s)[index]
}

func (s *Stack[T]) Size() int {
	return len(*s)
}

func (s *Stack[T]) IsEmpty() bool {
    return s.Size() == 0
}