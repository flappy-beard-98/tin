package core

type Set[T comparable] map[T]struct{}

func (o *Set[T]) ToArray() []T {
	result := make([]T, 0)
	for k := range *o {
		result = append(result, k)
	}
	return result
}

func (o *Set[T]) Add(item ...T) {
	for _, i := range item {
		(*o)[i] = struct{}{}
	}
}

func (o *Set[T]) Remove(item ...T) {
	for _, i := range item {
		delete(*o, i)
	}
}

func (o *Set[T]) IsEmpty() bool {
	return len(*o) == 0
}

func (o *Set[T]) Size() int {
	return len(*o)
}

func (o *Set[T]) Clear() {
	*o = Set[T]{}
}

func (o *Set[T]) Contains(item T) bool {
	_, ok := (*o)[item]
	return ok
}

func (o *Set[T]) Equals(other Set[T]) bool {
	for item := range *o {
		if !other.Contains(item) {
			return false
		}
	}
	return true
}
