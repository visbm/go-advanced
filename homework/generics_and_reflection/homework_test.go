package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(person Person) string {
	vp := reflect.ValueOf(person)
	vt := reflect.TypeOf(person)

	var s string

	for i := 0; i < vt.NumField(); i++ {
		f := vt.Field(i)
		v := vp.Field(i)

		tag := f.Tag.Get("properties")
		if tag == "" {
			continue
		}

		tags := strings.Split(tag, ",")
		if len(tags) >= 2 {
			if tags[1] == "omitempty" && v.IsZero() {
				continue
			}
		}
		if i == vt.NumField()-1 {
			s += fmt.Sprint(tags[0], "=", v)
			break
		}
		s += fmt.Sprint(tags[0], "=", v, "\n")

	}
	return s
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
