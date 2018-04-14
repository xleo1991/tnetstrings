package tnetstrings

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

const limit = 10

var SizeLimitExceeded = errors.New("size limit exceeded")

type TypeMismatch struct {
	Type uint8
	Kind reflect.Kind
}

func (t TypeMismatch) Error() string {
	return fmt.Sprintf("type mismatch: %s, %v", string(t.Type), t.Kind)
}

type InvalidType struct {
	reflect.Type
}

func (i InvalidType) Error() string {
	return fmt.Sprintf("invalid type: %v", i.Type)
}

type SyntaxError struct {
	Offset uint64
}

func (s SyntaxError) Error() string {
	return fmt.Sprintf("syntax error: %d", s.Offset)
}

type Decoder struct {
	*bufio.Reader
}

func (d *Decoder) Decode(val interface{}) error {
	size, err := d.size()
	if err != nil {
		return err
	}
	data := make([]uint8, size)
	if _, err := io.ReadFull(d, data[:]); err != nil {
		return err
	}
	t, err := d.ReadByte()
	if err != nil {
		return err
	}
	rv := reflect.ValueOf(val)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return InvalidType{Type: rv.Type()}
	}
	rv = reflect.Indirect(rv)

	k := rv.Kind()
	switch (pair{t: t, k: k}) {
	case pair{t: ',', k: reflect.Interface}, pair{t: ',', k: reflect.String}:
		rv.Set(reflect.ValueOf(string(data)))
	case pair{t: '#', k: reflect.Interface}:
		i, err := strconv.ParseInt(string(data), 0, 64)
		if err != nil {
			return err
		}
		rv.Set(reflect.ValueOf(i))
	case pair{t: '#', k: reflect.Int}, pair{t: '#', k: reflect.Int8}, pair{t: '#', k: reflect.Int16}, pair{t: '#', k: reflect.Int32}, pair{t: '#', k: reflect.Int64}:
		i, err := strconv.ParseInt(string(data), 0, int(rv.Type().Size()))
		if err != nil {
			return err
		}
		rv.SetInt(i)
	case pair{t: '#', k: reflect.Uint}, pair{t: '#', k: reflect.Uint8}, pair{t: '#', k: reflect.Uint16}, pair{t: '#', k: reflect.Uint32}, pair{t: '#', k: reflect.Uint64}:
		i, err := strconv.ParseUint(string(data), 0, int(rv.Type().Size()))
		if err != nil {
			return err
		}
		rv.SetUint(i)
	case pair{t: '^', k: reflect.Interface}:
		f, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return err
		}
		rv.Set(reflect.ValueOf(f))
	case  pair{t: '^', k: reflect.Float32}, pair{t: '^', k: reflect.Float64}:
		f, err := strconv.ParseFloat(string(data), int(rv.Type().Size()))
		if err != nil {
			return err
		}
		rv.SetFloat(f)
	case pair{t: '!', k: reflect.Interface}, pair{t: '!', k: reflect.Bool}:
		b, err := strconv.ParseBool(string(data))
		if err != nil {
			return err
		}
		rv.Set(reflect.ValueOf(b))
	case pair{t: '~', k: reflect.Interface}, pair{t: '~', k: reflect.Ptr}:
		rv.Set(reflect.Zero(rv.Type()))
	case pair{t: '}', k: reflect.Interface}, pair{t: '}', k: reflect.Map}:
		var m reflect.Value
		if k == reflect.Interface {
			m = reflect.ValueOf(map[string]interface{}{})
		} else {
			m = reflect.MakeMap(rv.Type())
		}
		d := Decoder{Reader: bufio.NewReader(bytes.NewReader(data))}
		var key string
		var val interface{}
		for d.More() {
			if err := d.Decode(&key); err != nil {
				return err
			}
			if err := d.Decode(&val); err != nil {
				return err
			}
			m.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
		}
		rv.Set(m)
	case pair{t: '}', k: reflect.Struct}:
		d := Decoder{Reader: bufio.NewReader(bytes.NewReader(data))}
		for d.More() {
			var key string
			if err := d.Decode(&key); err != nil {
				return err
			}
			if err := d.Decode(rv.FieldByName(key).Addr().Interface()); err != nil {
				return err
			}
		}
	case pair{t: ']', k: reflect.Array}:
		d := Decoder{Reader: bufio.NewReader(bytes.NewReader(data))}
		for i := 0; i < rv.Len(); i++ {
			if !d.More() {
				rv.Index(i).Set(reflect.Zero(rv.Type().Elem()))
				continue
			}
			if err := d.Decode(rv.Index(i).Addr().Interface()); err != nil {
				return err
			}
		}
	case pair{t: ']', k: reflect.Interface}, pair{t: ']', k: reflect.Slice}:
		var s reflect.Value
		if k == reflect.Interface {
			s = reflect.MakeSlice(reflect.TypeOf([]interface{}{}), 0, strings.Count(string(data), ":"))
		} else {
			s = reflect.MakeSlice(rv.Type(), 0, strings.Count(string(data), ":"))
		}

		d := Decoder{Reader: bufio.NewReader(bytes.NewReader(data))}
		var e interface{}
		for d.More() {
			if err := d.Decode(&e); err != nil {
				return err
			}
			s = reflect.Append(s, reflect.ValueOf(e))
		}
		rv.Set(s)
	default:
		return TypeMismatch{Type: t, Kind: k}
	}
	return nil
}

func (d *Decoder) More() bool {
	_, err := d.Reader.Peek(1)
	return err == nil
}

func (d *Decoder) size() (uint64, error) {
	var size uint64
	for i := 0; i < limit; i++ {
		if i == limit {
			return 0, SizeLimitExceeded
		}
		b, err := d.ReadByte()
		if err != nil {
			return 0, err
		}
		if ':' == b {
			break
		}
		if '0' <= b && '9' >= b {
			size = 10*size + uint64(b-'0')
		}
		// TODO: check invalid bytes
	}
	return size, nil
}

type pair struct {
	t uint8
	k reflect.Kind
}
