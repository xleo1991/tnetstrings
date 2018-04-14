package tnetstrings

import (
	"bytes"
	"reflect"
	"testing"
)

func TestEncoder_Encode(t *testing.T) {
	intVar := 1

	testCases := []struct {
		title string
		in    interface{}
		out   string
		err   error
	}{
		{
			title: "string",
			in:    "hello, world",
			out:   "12:hello, world,",
		},
		{
			title: "non-ascii string",
			in:    "日本語",
			out:   "9:日本語,",
		},
		{
			title: "true",
			in:    true,
			out:   "4:true!",
		},
		{
			title: "false",
			in:    false,
			out:   "5:false!",
		},
		{
			title: "positive int",
			in:    1,
			out:   "1:1#",
		},
		{
			title: "negative int",
			in:    -1,
			out:   "2:-1#",
		},
		{
			title: "positive float",
			in:    1.0,
			out:   "8:1.000000^",
		},
		{
			title: "nil",
			in:    nil,
			out:   "0:~",
		},
		{
			title: "empty map",
			in:    map[string]interface{}{},
			out:   "0:}",
		},
		{
			title: "map",
			in: map[string]interface{}{
				"foo": nil,
				"bar": 1,
			},
			out: "19:3:bar,1:1#3:foo,0:~}",
		},
		{
			title: "empty struct",
			in:    struct{}{},
			out:   "0:}",
		},
		{
			title: "struct with unexported field",
			in: struct {
				field int
			}{
				field: 1,
			},
			out: "0:}",
		},
		{
			title: "struct with exported field",
			in: struct {
				Field int
			}{
				Field: 1,
			},
			out: "12:5:Field,1:1#}",
		},
		{
			title: "struct with named field",
			in: struct {
				Field int `tnetstrings:"myName"`
			}{
				Field: 1,
			},
			out: "13:6:myName,1:1#}",
		},
		{
			title: "struct with named, omitempty field",
			in: struct {
				Field *int `tnetstrings:"myName,omitempty"`
			}{
				Field: nil,
			},
			out: "0:}",
		},
		{
			title: "struct with omitempty field",
			in: struct {
				Field *int `tnetstrings:",omitempty"`
			}{
				Field: &intVar,
			},
			out: "12:5:Field,1:1#}",
		},
		{
			title: "struct with ignored field",
			in: struct {
				Field int `tnetstrings:"-"`
			}{
				Field: 1,
			},
			out: "0:}",
		},
		{
			title: "struct with weirdly named field",
			in: struct {
				Field int `tnetstrings:"-,"`
			}{
				Field: 1,
			},
			out: "8:1:-,1:1#}",
		},
		{
			title: "empty array",
			in:    [0]string{},
			out:   "0:]",
		},
		{
			title: "array",
			in:    [3]string{"foo", "bar", "baz"},
			out:   "18:3:foo,3:bar,3:baz,]",
		},
		{
			title: "empty byte array",
			in:    [0]uint8{},
			out:   "0:,",
		},
		{
			title: "byte array",
			in:    [3]uint8{'a', 'b', 'c'},
			out:   "3:abc,",
		},
		{
			title: "empty slice",
			in:    []string{},
			out:   "0:]",
		},
		{
			title: "slice",
			in:    []string{"foo", "bar", "baz"},
			out:   "18:3:foo,3:bar,3:baz,]",
		},
		{
			title: "empty byte slice",
			in:    []uint8{},
			out:   "0:,",
		},
		{
			title: "byte slice",
			in:    []uint8{'a', 'b', 'c'},
			out:   "3:abc,",
		},
		{
			title: "unsupported type (complex128)",
			in:    1i,
			err:   ErrUnsupportedType{Type: reflect.TypeOf(0i)},
		},
	}

	for _, tc := range testCases {
		var buf bytes.Buffer
		e := Encoder{
			Writer: &buf,
		}

		if err := e.Encode(tc.in); err != tc.err {
			t.Error(err)
		}
		if tc.out != buf.String() {
			t.Errorf("[%s] expected: %s, got: %s", tc.title, tc.out, buf.String())
		}
		buf.Reset()
		if err := e.Encode(&tc.in); err != tc.err {
			t.Error(err)
		}
		if tc.out != buf.String() {
			t.Errorf("[pointer %s] expected: %s, got: %s", tc.title, tc.out, buf.String())
		}
	}
}
