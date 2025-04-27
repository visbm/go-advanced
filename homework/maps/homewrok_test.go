package main

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"
	"reflect"
	"testing"
)

// go test -v homework_test.go

type OrderedMap[K constraints.Ordered, V any] struct {
	root *node[K, V]
	size int
}

type node[K constraints.Ordered, V any] struct {
	key   K
	val   V
	left  *node[K, V]
	right *node[K, V]
}

func NewOrderedMap[K constraints.Ordered, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{}
}

func (m *OrderedMap[K, V]) Insert(key K, value V) {
	if m.root == nil {
		m.root = &node[K, V]{
			key:   key,
			val:   value,
			left:  nil,
			right: nil,
		}
		m.size++
		return
	}

	if m.root.insert(key, value) {
		m.size++
	}

}

func (n *node[K, V]) insert(key K, value V) bool {
	if n.key == key {
		n.val = value
		return false
	}

	var child **node[K, V]
	if key < n.key {
		child = &n.left
	} else if key > n.key {
		child = &n.right
	} else {
		return false
	}

	if *child == nil {
		*child = &node[K, V]{
			key: key,
			val: value,
		}
		return true
	}

	return (*child).insert(key, value)
}

func (m *OrderedMap[K, V]) Erase(key K) {
	var changed bool
	m.root, changed = m.root.erase(key)
	if changed {
		m.size--
	}
}

func (n *node[K, V]) erase(key K) (*node[K, V], bool) {
	if n == nil {
		return nil, false
	}
	if key < n.key {
		var changed bool
		n.left, changed = n.left.erase(key)
		return n, changed
	}
	if key > n.key {
		var changed bool
		n.right, changed = n.right.erase(key)
		return n, changed
	}

	if n.left == nil {
		return n.right, true
	}
	if n.right == nil {
		return n.left, true
	}

	min := n.right
	for min.left != nil {
		min = min.left
	}
	n.key, n.val = min.key, min.val
	var changed bool
	n.right, changed = n.right.erase(min.key)
	return n, changed
}

func (m *OrderedMap[K, V]) Contains(key K) bool {
	return m.root.contains(key)
}

func (n *node[K, V]) contains(key K) bool {
	if n.key == key {
		return true
	}
	if key < n.key {
		if n.left == nil {
			return false
		}
		return n.left.contains(key)
	}
	if key > n.key {
		if n.right == nil {
			return false
		}
		return n.right.contains(key)
	}
	return false
}

func (m *OrderedMap[K, V]) Size() int {
	return m.size
}

func (m *OrderedMap[K, V]) ForEach(action func(K, V)) {
	m.root.ForEach(action)
}

func (n *node[K, V]) ForEach(action func(K, V)) {
	if n == nil {
		return
	}
	n.left.ForEach(action)
	action(n.key, n.val)
	n.right.ForEach(action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap[int, int]()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
