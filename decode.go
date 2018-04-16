package tnetstrings

import (
	"bufio"
	"bytes"
	"io"
	"reflect"
	"strconv"
	"strings"
)

const limit = 10

// Decoder is a streaming tnetstrings decoder.
type Decoder struct {
	*bufio.Reader
}

// NewDecoder returns a new Decoder instance.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{Reader: bufio.NewReader(r)}
}

// Decode decodes a tnetstring from the stream.
func (d *Decoder) Decode(val interface{}) error {
	size, err := d.size()
	if err != nil {
		return err
	}
	data := make([]uint8, size+1)
	if _, err = io.ReadFull(d, data[:]); err != nil {
		return err
	}

	t := data[len(data)-1]
	data = data[:len(data)-1]
	rv := reflect.Indirect(reflect.ValueOf(val))
	switch t {
	case ',':
		return decodeString(data, rv)
	case '#':
		return decodeInteger(data, rv)
	case '^':
		return decodeFloat(data, rv)
	case '!':
		return decodeBool(data, rv)
	case '~':
		return decodeNull(data, rv)
	case '}':
		return decodeDictionary(data, rv)
	case ']':
		return decodeList(data, rv)
	}
	return ErrInvalidTypeChar(t)
}

// More returns true iff the underlying stream can return more than 1 byte.
func (d *Decoder) More() bool {
	_, err := d.Reader.Peek(1)
	return err == nil
}

func (d *Decoder) size() (uint64, error) {
	var size uint64
	for i := 0; i < limit; i++ {
		b, err := d.ReadByte()
		if err != nil {
			return 0, err
		}
		switch b {
		case ':':
			return size, nil
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			size = 10*size + uint64(b-'0')
		default:
			return 0, ErrInvalidSizeChar(b)
		}
	}
	return 0, ErrSizeLimitExceeded
}

func decodeString(data []byte, rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Interface:
		if rv.Type().NumMethod() != 0 {
			return ErrUnsupportedType{Type: rv.Type()}
		}
		rv.Set(reflect.ValueOf(string(data)))
		return nil
	case reflect.String:
		rv.SetString(string(data))
		return nil
	default:
		return ErrUnsupportedType{Type: rv.Type()}
	}
}

func decodeInteger(data []byte, rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Interface:
		if rv.Type().NumMethod() != 0 {
			return ErrUnsupportedType{Type: rv.Type()}
		}
		i, err := strconv.ParseInt(string(data), 0, 64)
		if err != nil {
			return err
		}
		rv.Set(reflect.ValueOf(i))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(string(data), 0, int(rv.Type().Size()))
		if err != nil {
			return err
		}
		rv.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(string(data), 0, int(rv.Type().Size()))
		if err != nil {
			return err
		}
		rv.SetUint(i)
	default:
		return ErrUnsupportedType{Type: rv.Type()}
	}
	return nil
}

func decodeFloat(data []byte, rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Interface:
		if rv.Type().NumMethod() != 0 {
			return ErrUnsupportedType{Type: rv.Type()}
		}
		f, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return err
		}
		rv.Set(reflect.ValueOf(f))
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(string(data), int(rv.Type().Size()))
		if err != nil {
			return err
		}
		rv.SetFloat(f)
	default:
		return ErrUnsupportedType{Type: rv.Type()}
	}
	return nil
}

func decodeBool(data []byte, rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Interface:
		if rv.Type().NumMethod() != 0 {
			return ErrUnsupportedType{Type: rv.Type()}
		}
		b, err := strconv.ParseBool(string(data))
		if err != nil {
			return err
		}
		rv.Set(reflect.ValueOf(b))
	case reflect.Bool:
		b, err := strconv.ParseBool(string(data))
		if err != nil {
			return err
		}
		rv.SetBool(b)
	default:
		return ErrUnsupportedType{Type: rv.Type()}
	}
	return nil
}

func decodeNull(_ []byte, rv reflect.Value) error {
	rv.Set(reflect.Zero(rv.Type()))
	return nil
}

func decodeDictionary(data []byte, rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Interface:
		if rv.Type().NumMethod() != 0 {
			return ErrUnsupportedType{Type: rv.Type()}
		}
		return decodeDictionaryInterface(data, rv)
	case reflect.Map:
		return decodeDictionaryMap(data, rv)
	case reflect.Struct:
		return decodeDictionaryStruct(data, rv)
	default:
		return ErrUnsupportedType{Type: rv.Type()}
	}
}

func decodeDictionaryInterface(data []byte, rv reflect.Value) error {
	m := reflect.ValueOf(map[string]interface{}{})
	d := NewDecoder(bytes.NewReader(data))
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
	return nil
}

func decodeDictionaryMap(data []byte, rv reflect.Value) error {
	m := reflect.MakeMap(rv.Type())
	d := NewDecoder(bytes.NewReader(data))
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
	return nil
}

func decodeDictionaryStruct(data []byte, rv reflect.Value) error {
	d := NewDecoder(bytes.NewReader(data))
	tags := make(map[string]*tag, rv.NumField())
	for i := 0; i < rv.NumField(); i++ {
		tag := parseTag(rv.Type().Field(i))
		if tag == nil {
			continue
		}
		tags[tag.displayName] = tag
	}
	for d.More() {
		var key string
		if err := d.Decode(&key); err != nil {
			return err
		}

		var val interface{}
		if err := d.Decode(&val); err != nil {
			return err
		}

		rv.FieldByName(tags[key].name).Set(reflect.ValueOf(val))
	}
	return nil
}

func decodeList(data []byte, rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Array:
		return decodeListArray(data, rv)
	case reflect.Interface:
		if rv.Type().NumMethod() != 0 {
			return ErrUnsupportedType{Type: rv.Type()}
		}
		return decodeListInterface(data, rv)
	case reflect.Slice:
		return decodeListSlice(data, rv)
	default:
		return ErrUnsupportedType{Type: rv.Type()}
	}
}

func decodeListArray(data []byte, rv reflect.Value) error {
	d := NewDecoder(bytes.NewReader(data))
	for i := 0; i < rv.Len(); i++ {
		if !d.More() {
			rv.Index(i).Set(reflect.Zero(rv.Type().Elem()))
			continue
		}
		if err := d.Decode(rv.Index(i).Addr().Interface()); err != nil {
			return err
		}
	}
	return nil
}

func decodeListInterface(data []byte, rv reflect.Value) error {
	s := reflect.MakeSlice(reflect.TypeOf([]interface{}{}), 0, strings.Count(string(data), ":"))
	d := NewDecoder(bytes.NewReader(data))
	var e interface{}
	for d.More() {
		if err := d.Decode(&e); err != nil {
			return err
		}
		s = reflect.Append(s, reflect.ValueOf(e))
	}
	rv.Set(s)
	return nil
}

func decodeListSlice(data []byte, rv reflect.Value) error {
	s := reflect.MakeSlice(rv.Type(), 0, strings.Count(string(data), ":"))
	d := NewDecoder(bytes.NewReader(data))
	var e interface{}
	for d.More() {
		if err := d.Decode(&e); err != nil {
			return err
		}
		s = reflect.Append(s, reflect.ValueOf(e))
	}
	rv.Set(s)
	return nil
}
