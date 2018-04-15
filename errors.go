package tnetstrings

import (
	"errors"
	"fmt"
	"reflect"
)

// ErrUnsupportedType means the argument type is not eligible to encode/decode.
type ErrUnsupportedType struct {
	reflect.Type
}

func (e ErrUnsupportedType) Error() string {
	return fmt.Sprintf("unsupported type: %s", e.Type)
}

// ErrNonStringKey means passed map has non-string key which is not accepted in tnetstrings.
var ErrNonStringKey = errors.New("non string key")

// ErrSizeLimitExceeded means SIZE is longer than 9 digits.
var ErrSizeLimitExceeded = errors.New("size limit exceeded")

// ErrInvalidSizeChar means there's a non-digit character before `:`.
type ErrInvalidSizeChar uint8

func (e ErrInvalidSizeChar) Error() string {
	return fmt.Sprintf("invalid size char: %s", string(e))
}

// ErrInvalidTypeChar means the payload doesn't end with one of `,#^!~}]`.
type ErrInvalidTypeChar uint8

func (e ErrInvalidTypeChar) Error() string {
	return fmt.Sprintf("invalid type char: %s", string(e))
}
