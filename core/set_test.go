package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdd(t *testing.T) {
	exp := "a"
	sut := Set[string]{}

	sut.Add(exp)

	act := sut.ToArray()

	assert.Len(t, act, 1)
	assert.Equal(t, exp, act[0])
}

func TestRemove(t *testing.T) {
	exp := "a"
	sut := Set[string]{}

	sut.Add(exp)
	sut.Remove(exp)

	act := sut.IsEmpty()

	assert.True(t, act)
}

func TestContains(t *testing.T) {
	exp := "a"
	sut := Set[string]{}

	sut.Add(exp)

	act := sut.Contains(exp)

	assert.True(t, act)
}

func TestClear(t *testing.T) {
	sut := Set[string]{}
	sut.Add("a")
	sut.Clear()
	act := sut.IsEmpty()
	assert.True(t, act)
}

func TestSize(t *testing.T) {
	exp := 3
	sut := Set[string]{}
	sut.Add("a", "b", "c")

	act := sut.Size()
	assert.Equal(t, exp, act)
}

func TestEquals(t *testing.T) {
	exp := "a"
	sut := Set[string]{}
	sut.Add(exp)

	other := Set[string]{}
	other.Add(exp)

	act := sut.Equals(other)

	assert.True(t, act)
}

func TestNotEquals(t *testing.T) {
	exp := "a"
	sut := Set[string]{}
	sut.Add(exp)

	other := Set[string]{}

	act := sut.Equals(other)

	assert.False(t, act)
}
