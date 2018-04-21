package tnetstrings

import (
	"bufio"
	"bytes"
	"io"
	"reflect"
	"strconv"
	"testing"
)

func TestDecoder_Decode_string(t *testing.T) {
	testCases := []struct {
		title string
		in    string
		out   string
		err   error
	}{
		{
			title: "empty",
			in:    "0:,",
			out:   "",
		},
		{
			title: "just",
			in:    "13:Hello, World!,",
			out:   "Hello, World!",
		},
		{
			title: "bigger size",
			in:    "1000:foo,",
			err:   io.ErrUnexpectedEOF,
		},
		{
			title: "smaller size",
			in:    "2:foo,",
			err:   ErrInvalidTypeChar('o'),
		},
	}

	for _, tc := range testCases {
		d := Decoder{
			Reader: bufio.NewReader(bytes.NewReader([]byte(tc.in))),
		}
		var s string
		if err := d.Decode(&s); err != nil && !reflect.DeepEqual(tc.err, err) {
			t.Errorf("[%s] expected: %v, got: %v", tc.title, tc.err, err)
		}
		if tc.out != s {
			t.Errorf("[%s] expected: %s, got: %s", tc.title, tc.out, s)
		}
	}
}

func TestDecoder_Decode_int(t *testing.T) {
	testCases := []struct {
		title string
		in    string
		out   int
		err   error
	}{
		{
			title: "empty",
			in:    "0:#",
			err:   &strconv.NumError{Func: "ParseInt", Num: "", Err: strconv.ErrSyntax},
		},
		{
			title: "positive",
			in:    "5:12345#",
			out:   12345,
		},
		{
			title: "negative",
			in:    "6:-12345#",
			out:   -12345,
		},
		{
			title: "bigger size",
			in:    "1000:1#",
			err:   io.ErrUnexpectedEOF,
		},
		{
			title: "smaller size",
			in:    "2:123#",
			err:   ErrInvalidTypeChar('3'),
		},
	}

	for _, tc := range testCases {
		d := Decoder{
			Reader: bufio.NewReader(bytes.NewReader([]byte(tc.in))),
		}
		var i int
		if err := d.Decode(&i); err != nil && !reflect.DeepEqual(tc.err, err) {
			t.Errorf("[%s] expected: %v, got: %v", tc.title, tc.err, err)
		}
		if tc.out != i {
			t.Errorf("[%s] expected: %d, got: %d", tc.title, tc.out, i)
		}
	}
}

func TestDecoder_Decode_uint(t *testing.T) {
	testCases := []struct {
		title string
		in    string
		out   uint
		err   error
	}{
		{
			title: "empty",
			in:    "0:#",
			err:   &strconv.NumError{Func: "ParseUint", Num: "", Err: strconv.ErrSyntax},
		},
		{
			title: "positive",
			in:    "5:12345#",
			out:   12345,
		},
		{
			title: "negative",
			in:    "6:-12345#",
			err:   &strconv.NumError{Func: "ParseUint", Num: "-12345", Err: strconv.ErrSyntax},
		},
		{
			title: "bigger size",
			in:    "1000:1#",
			err:   io.ErrUnexpectedEOF,
		},
		{
			title: "smaller size",
			in:    "2:123#",
			err:   ErrInvalidTypeChar('3'),
		},
	}

	for _, tc := range testCases {
		d := Decoder{
			Reader: bufio.NewReader(bytes.NewReader([]byte(tc.in))),
		}
		var i uint
		if err := d.Decode(&i); err != nil && !reflect.DeepEqual(tc.err, err) {
			t.Errorf("[%s] expected: %v, got: %v", tc.title, tc.err, err)
		}
		if tc.out != i {
			t.Errorf("[%s] expected: %d, got: %d", tc.title, tc.out, i)
		}
	}
}

func TestDecoder_Decode_float32(t *testing.T) {
	testCases := []struct {
		title string
		in    string
		out   float32
		err   error
	}{
		{
			title: "empty",
			in:    "0:^",
			err:   &strconv.NumError{Func: "ParseFloat", Num: "", Err: strconv.ErrSyntax},
		},
		{
			title: "positive",
			in:    "6:.12345^",
			out:   .12345,
		},
		{
			title: "negative",
			in:    "7:-.12345^",
			out:   -.12345,
		},
		{
			title: "bigger size",
			in:    "1000:.1^",
			err:   io.ErrUnexpectedEOF,
		},
		{
			title: "smaller size",
			in:    "2:.123^",
			err:   ErrInvalidTypeChar('2'),
		},
	}

	for _, tc := range testCases {
		d := Decoder{
			Reader: bufio.NewReader(bytes.NewReader([]byte(tc.in))),
		}
		var f float32
		if err := d.Decode(&f); err != nil && !reflect.DeepEqual(tc.err, err) {
			t.Errorf("[%s] expected: %v, got: %v", tc.title, tc.err, err)
		}
		if tc.out != f {
			t.Errorf("[%s] expected: %f, got: %f", tc.title, tc.out, f)
		}
	}
}

func TestDecoder_Decode_bool(t *testing.T) {
	testCases := []struct {
		title string
		in    string
		out   bool
		err   error
	}{
		{
			title: "empty",
			in:    "0:!",
			err:   &strconv.NumError{Func: "ParseBool", Num: "", Err: strconv.ErrSyntax},
		},
		{
			title: "true",
			in:    "4:true!",
			out:   true,
		},
		{
			title: "false",
			in:    "5:false!",
			out:   false,
		},
		{
			title: "bigger size",
			in:    "1000:true!",
			err:   io.ErrUnexpectedEOF,
		},
		{
			title: "smaller size",
			in:    "2:false!",
			err:   ErrInvalidTypeChar('l'),
		},
	}

	for _, tc := range testCases {
		d := Decoder{
			Reader: bufio.NewReader(bytes.NewReader([]byte(tc.in))),
		}
		var b bool
		if err := d.Decode(&b); err != nil && !reflect.DeepEqual(tc.err, err) {
			t.Errorf("[%s] expected: %v, got: %v", tc.title, tc.err, err)
		}
		if tc.out != b {
			t.Errorf("[%s] expected: %t, got: %t", tc.title, tc.out, b)
		}
	}
}

func TestDecoder_Decode_null(t *testing.T) {
	testCases := []struct {
		title string
		in    string
		out   *int
		err   error
	}{
		{
			title: "empty",
			in:    "0:~",
			out:   nil,
		},
		{
			title: "data",
			in:    "3:abc~",
			out:   nil,
		},
		{
			title: "bigger size",
			in:    "1000:~",
			err:   io.ErrUnexpectedEOF,
		},
	}

	for _, tc := range testCases {
		d := Decoder{
			Reader: bufio.NewReader(bytes.NewReader([]byte(tc.in))),
		}
		var p *int
		if err := d.Decode(&p); err != nil && !reflect.DeepEqual(tc.err, err) {
			t.Errorf("[%s] expected: %v, got: %v", tc.title, tc.err, err)
		}
		if tc.out != p {
			t.Errorf("[%s] expected: %v, got: %v", tc.title, tc.out, p)
		}
	}
}

func TestDecoder_Decode_map(t *testing.T) {
	testCases := []struct {
		title string
		in    string
		out   map[string]interface{}
		err   error
	}{
		{
			title: "empty",
			in:    "0:}",
			out:   map[string]interface{}{},
		},
		{
			title: "just",
			in:    "12:3:foo,3:bar,}",
			out: map[string]interface{}{
				"foo": "bar",
			},
		},
		{
			title: "bigger size",
			in:    "1000:3:foo,3:bar,}",
			err:   io.ErrUnexpectedEOF,
		},
		{
			title: "smaller size",
			in:    "2:3:foo,3:bar,}",
			err:   ErrInvalidTypeChar('f'),
		},
	}

	for _, tc := range testCases {
		d := Decoder{
			Reader: bufio.NewReader(bytes.NewReader([]byte(tc.in))),
		}
		var m map[string]interface{}
		if err := d.Decode(&m); err != nil && !reflect.DeepEqual(tc.err, err) {
			t.Errorf("[%s] expected: %v, got: %v", tc.title, tc.err, err)
		}
		if !reflect.DeepEqual(tc.out, m) {
			t.Errorf("[%s] expected: %#v, got: %#v", tc.title, tc.out, m)
		}
	}
}

func TestDecoder_Decode_struct(t *testing.T) {
	type s struct {
		NoTag          string
		Ignored        string `tnetstrings:"-"`
		Named          string `tnetstrings:"foo"`
		OmitEmpty      string `tnetstrings:",omitempty"`
		NamedOmitEmpty string `tnetstrings:"bar,omitempty"`
	}

	testCases := []struct {
		title string
		in    string
		out   s
		err   error
	}{
		{
			title: "empty",
			in:    "0:}",
			out:   s{},
		},
		{
			title: "just",
			in:    "62:5:NoTag,3:bar,3:foo,3:foo,9:OmitEmpty,9:OmitEmpty,3:bar,3:bar,}",
			out: s{
				NoTag:          "bar",
				Named:          "foo",
				OmitEmpty:      "OmitEmpty",
				NamedOmitEmpty: "bar",
			},
		},
		{
			title: "bigger size",
			in:    "1000:3:foo,3:bar,}",
			err:   io.ErrUnexpectedEOF,
		},
		{
			title: "smaller size",
			in:    "2:3:foo,3:bar,}",
			err:   ErrInvalidTypeChar('f'),
		},
	}

	for _, tc := range testCases {
		d := Decoder{
			Reader: bufio.NewReader(bytes.NewReader([]byte(tc.in))),
		}
		var s s
		if err := d.Decode(&s); err != nil && !reflect.DeepEqual(tc.err, err) {
			t.Errorf("[%s] expected: %v, got: %v", tc.title, tc.err, err)
		}
		if !reflect.DeepEqual(tc.out, s) {
			t.Errorf("[%s] expected: %#v, got: %#v", tc.title, tc.out, s)
		}
	}
}

func TestDecoder_Decode_array(t *testing.T) {
	testCases := []struct {
		title string
		in    string
		out   [2]string
		err   error
	}{
		{
			title: "empty",
			in:    "0:]",
			out:   [2]string{},
		},
		{
			title: "just",
			in:    "12:3:foo,3:bar,]",
			out: [2]string{
				"foo",
				"bar",
			},
		},
		{
			title: "fewer elements",
			in:    "6:3:foo,]",
			out: [2]string{
				"foo",
			},
		},
		{
			title: "more elements",
			in:    "18:3:foo,3:bar,3:baz,]",
			out: [2]string{
				"foo",
				"bar",
			},
		},
		{
			title: "bigger size",
			in:    "1000:3:foo,3:bar,]",
			err:   io.ErrUnexpectedEOF,
		},
		{
			title: "smaller size",
			in:    "2:3:foo,3:bar,]",
			err:   ErrInvalidTypeChar('f'),
		},
	}

	for _, tc := range testCases {
		d := Decoder{
			Reader: bufio.NewReader(bytes.NewReader([]byte(tc.in))),
		}
		var a [2]string
		if err := d.Decode(&a); err != nil && !reflect.DeepEqual(tc.err, err) {
			t.Errorf("[%s] expected: %v, got: %v", tc.title, tc.err, err)
		}
		if !reflect.DeepEqual(tc.out, a) {
			t.Errorf("[%s] expected: %#v, got: %#v", tc.title, tc.out, a)
		}
	}
}

func TestDecoder_Decode_slice(t *testing.T) {
	testCases := []struct {
		title string
		in    string
		out   []string
		err   error
	}{
		{
			title: "empty",
			in:    "0:]",
			out:   []string{},
		},
		{
			title: "just",
			in:    "12:3:foo,3:bar,]",
			out: []string{
				"foo",
				"bar",
			},
		},
		{
			title: "fewer elements",
			in:    "6:3:foo,]",
			out: []string{
				"foo",
			},
		},
		{
			title: "more elements",
			in:    "18:3:foo,3:bar,3:baz,]",
			out: []string{
				"foo",
				"bar",
				"baz",
			},
		},
		{
			title: "bigger size",
			in:    "1000:3:foo,3:bar,]",
			err:   io.ErrUnexpectedEOF,
		},
		{
			title: "smaller size",
			in:    "2:3:foo,3:bar,]",
			err:   ErrInvalidTypeChar('f'),
		},
	}

	for _, tc := range testCases {
		d := Decoder{
			Reader: bufio.NewReader(bytes.NewReader([]byte(tc.in))),
		}
		var a []string
		if err := d.Decode(&a); err != nil && !reflect.DeepEqual(tc.err, err) {
			t.Errorf("[%s] expected: %v, got: %v", tc.title, tc.err, err)
		}
		if !reflect.DeepEqual(tc.out, a) {
			t.Errorf("[%s] expected: %#v, got: %#v", tc.title, tc.out, a)
		}
	}
}

func TestDecoder_Decode_interface(t *testing.T) {
	testCases := []struct {
		title string
		in    string
		out   interface{}
		err   error
	}{
		{
			title: "string",
			in:    "3:abc,",
			out:   "abc",
		},
		{
			title: "integer",
			in:    "3:123#",
			out:   int64(123),
		},
		{
			title: "float",
			in:    "4:.123^",
			out:   float64(.123),
		},
		{
			title: "boolean",
			in:    "4:true!",
			out:   true,
		},
		{
			title: "null",
			in:    "0:~",
			out:   nil,
		},
		{
			title: "dictionary",
			in:    "12:3:foo,3:bar,}",
			out: map[string]interface{}{
				"foo": "bar",
			},
		},
		{
			title: "list",
			in:    "12:3:foo,3:bar,]",
			out: []interface{}{
				"foo",
				"bar",
			},
		},
	}

	for _, tc := range testCases {
		d := Decoder{
			Reader: bufio.NewReader(bytes.NewReader([]byte(tc.in))),
		}
		var i interface{}
		if err := d.Decode(&i); err != nil && !reflect.DeepEqual(tc.err, err) {
			t.Errorf("[%s] expected: %v, got: %v", tc.title, tc.err, err)
		}
		if !reflect.DeepEqual(tc.out, i) {
			t.Errorf("[%s] expected: %#v (%s), got: %#v (%s)", tc.title, tc.out, reflect.TypeOf(tc.out), i, reflect.TypeOf(i))
		}
	}
}
