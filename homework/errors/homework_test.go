package main

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MyError struct {
}

func (e *MyError) Error() string {
	return "my error"
}

type MultiError struct {
	errsA []error
}

func (e *MultiError) Error() string {
	if len(e.errsA) == 0 {
		return ""
	}

	b := strings.Builder{}
	b.WriteString(strconv.Itoa(len(e.errsA)) + " errors occured:\n")
	for _, err := range e.errsA {
		b.WriteString("\t* " + err.Error())
	}
	b.WriteString("\n")
	return b.String()
}

func (e *MultiError) As(target any) bool {
	for _, err := range e.errsA {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}

func (e *MultiError) Is(target error) bool {
	for _, err := range e.errsA {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func (e *MultiError) Unwrap() error {
	if len(e.errsA) <= 1 {
		return nil
	}
	return &MultiError{errsA: e.errsA[0 : len(e.errsA)-1]}
}

func Append(err error, errs ...error) *MultiError {
	if err == nil && len(errs) == 0 {
		return nil
	}

	var myErr *MultiError
	if errors.As(err, &myErr) {
		myErr.errsA = append(myErr.errsA, errs...)
		return myErr
	}

	errorList := make([]error, 0, len(errs)+1)
	if err != nil {
		errorList = append(errorList, err)
	}

	errorList = append(errorList, errs...)

	return &MultiError{
		errsA: errorList,
	}
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))
	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)

	//Is
	err = Append(err, sql.ErrNoRows)
	assert.True(t, errors.Is(err, sql.ErrNoRows))
	assert.False(t, errors.Is(err, sql.ErrConnDone))

	//As
	err = Append(err, &MyError{})
	var target *MyError
	assert.True(t, errors.As(err, &target))

	//Unwrap
	err = errors.Unwrap(errors.Unwrap(err))
	assert.EqualError(t, err, expectedMessage)
}
