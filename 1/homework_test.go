package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

// go test -v homework_test.go

func ToLittleEndian(number uint32) uint32 {
	bytes := make([]byte, 4)

	unsPointerNumber := unsafe.Pointer(&number)
	for i := 0; i < 4; i++ {
		c := *(*byte)(unsafe.Add(unsPointerNumber, i))
		bytes[i] = c
	}

	newNumber := uint32(bytes[0])<<24 | uint32(bytes[1])<<16 | uint32(bytes[2])<<8 | uint32(bytes[3])

	return newNumber
}

func ToLittleEndianGeneric[T ~uint16 | ~uint32 | ~uint64](number T) T {
	size := unsafe.Sizeof(number)
	bytes := make([]byte, size)

	unsPointerNumber := unsafe.Pointer(&number)
	for i := 0; i < int(size); i++ {
		c := *(*byte)(unsafe.Add(unsPointerNumber, i))
		bytes[i] = c
	}

	var result T
	for i := range bytes {
		result = (result << 8) | T(bytes[i])
	}

	return result
}

func TestСonversion(t *testing.T) {
	tests := map[string]struct {
		number uint32
		result uint32
	}{
		"test case #1": {
			number: 0x00000000,
			result: 0x00000000,
		},
		"test case #2": {
			number: 0xFFFFFFFF,
			result: 0xFFFFFFFF,
		},
		"test case #3": {
			number: 0x00FF00FF,
			result: 0xFF00FF00,
		},
		"test case #4": {
			number: 0x0000FFFF,
			result: 0xFFFF0000,
		},
		"test case #5": {
			number: 0x01020304,
			result: 0x04030201,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndianGeneric[uint32](test.number)
			assert.Equal(t, test.result, result)
		})
	}
}
