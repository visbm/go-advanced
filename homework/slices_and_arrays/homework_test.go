package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type integers interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~int
}

type CircularQueue[T integers] struct {
	values []T
	size   int
	first  int
	last   int
}

func NewCircularQueue[T integers](size int) CircularQueue[T] {
	return CircularQueue[T]{
		values: make([]T, size),
	}
}

func (q *CircularQueue[T]) Push(value T) bool {
	if q.Full() {
		return false
	}

	q.values[q.last] = value
	q.size++
	q.last = (q.last + 1) % cap(q.values)

	return true
}

func (q *CircularQueue[T]) Pop() bool {
	if q.Empty() {
		return false
	}

	q.size--
	q.first++
	return true
}

func (q *CircularQueue[T]) Front() T {
	if q.Empty() {
		return -1
	}

	return q.values[q.first]
}

func (q *CircularQueue[T]) Back() T {
	if q.Empty() {
		return -1
	}
	index := (q.last - 1 + q.size) % q.size
	return q.values[index]
}

func (q *CircularQueue[T]) Empty() bool {
	return q.size == 0
}

func (q *CircularQueue[T]) Full() bool {
	return q.size == cap(q.values)
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[int64](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, int64(-1), queue.Front())
	assert.Equal(t, int64(-1), queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int64{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, int64(1), queue.Front())
	assert.Equal(t, int64(3), queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(int64(4)))

	assert.True(t, reflect.DeepEqual([]int64{4, 2, 3}, queue.values))

	assert.Equal(t, int64(2), queue.Front())
	assert.Equal(t, int64(4), queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}
