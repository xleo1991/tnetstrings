package tnetstrings

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
)

// Encoder is a streaming tnetstrings encoder.
type Encoder struct {
	io.Writer
}

// NewEncoder returns a new Encoder instance.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{Writer: w}
}

// Encode encodes a value into tnetstring.
func (e *Encoder) Encode(val interface{}) error {
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String:
		s := val.(string)
		// _, err := fmt.Fprintf(e, "%d:%s;", len(d), d)
		d := []byte(s)
		_, err := fmt.Fprintf(e, "%d:", len(d))
		if err != nil {
			return err
		}
		_, err = e.Write(d)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(e, ";")
		return err
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s := fmt.Sprintf("%d", val)
		_, err := fmt.Fprintf(e, "%d:%s#", len(s), s)
		return err
	case reflect.Float32, reflect.Float64:
		s := fmt.Sprintf("%f", val)
		_, err := fmt.Fprintf(e, "%d:%s^", len(s), s)
		return err
	case reflect.Bool:
		s := strconv.FormatBool(val.(bool))
		_, err := fmt.Fprintf(e, "%d:%s!", len(s), s)
		return err
	case reflect.Invalid:
		_, err := fmt.Fprint(e, "0:~")
		return err
	case reflect.Ptr:
		v = v.Elem()
		return e.Encode(v.Interface())
	case reflect.Map:
		return e.encodeMap(v)
	case reflect.Struct:
		return e.encodeStruct(v)
	case reflect.Array, reflect.Slice:
		return e.encodeSlice(v)
	}
	return ErrUnsupportedType{Type: v.Type()}
}

func (e *Encoder) encodeMap(v reflect.Value) error {
	var buf bytes.Buffer
	f := NewEncoder(&buf)
	ks := v.MapKeys()
	sort.Slice(ks, func(i, j int) bool {
		return ks[i].String() < ks[j].String()
	})
	for _, k := range ks {
		if err := f.Encode(k.Interface()); err != nil {
			return err
		}
		if err := f.Encode(v.MapIndex(k).Interface()); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(e, "%d:%s}", buf.Len(), buf.Bytes())
	return err
}

func (e *Encoder) encodeStruct(v reflect.Value) error {
	var buf bytes.Buffer
	f := NewEncoder(&buf)
	for i := 0; i < v.NumField(); i++ {
		ft := v.Type().Field(i)
		fv := v.Field(i)
		if !fv.CanInterface() {
			continue
		}

		tag := parseTag(ft)
		if tag == nil {
			continue
		}
		if tag.omitEmpty && fv == reflect.Zero(ft.Type) {
			continue
		}

		if err := f.Encode(tag.displayName); err != nil {
			return err
		}
		if err := f.Encode(fv.Interface()); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(e, "%d:%s}", buf.Len(), buf.Bytes())
	return err
}

func (e *Encoder) encodeSlice(v reflect.Value) error {
	if v.Type().Elem().Kind() == reflect.Uint8 {
		s := fmt.Sprintf("%s", v.Interface())
		_, err := fmt.Fprintf(e, "%d:%s,", len(s), s)
		return err
	}
	var buf bytes.Buffer
	f := NewEncoder(&buf)
	for i := 0; i < v.Len(); i++ {
		if err := f.Encode(v.Index(i).Interface()); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(e, "%d:%s]", buf.Len(), buf.Bytes())
	return err
}
