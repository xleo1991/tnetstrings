package tnetstrings

import (
	"bufio"
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

// ErrTypeMismatch means no decoding is defined.
type ErrTypeMismatch struct {
	Tag  uint8
	Type reflect.Type
}

func (t ErrTypeMismatch) Error() string {
	return fmt.Sprintf("type mismatch: %s, %v", string(t.Tag), t.Type)
}

// Decoder is a streaming tnetstrings decoder.
type Decoder struct {
	*bufio.Reader
}
