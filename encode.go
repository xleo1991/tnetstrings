package tnetstrings

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
)

type Encoder struct {
	io.Writer
}

func (e *Encoder) Encode(val interface{}) error {
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String:
		s := val.(string)
		fmt.Fprintf(e, "%d:%s,", len(s), s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s := fmt.Sprintf("%d", val)
		fmt.Fprintf(e, "%d:%s#", len(s), s)
	case reflect.Float32, reflect.Float64:
		s := fmt.Sprintf("%f", val)
		fmt.Fprintf(e, "%d:%s^", len(s), s)
	case reflect.Bool:
		if val.(bool) {
			fmt.Fprint(e, "4:true!")
		} else {
			fmt.Fprint(e, "5:false!")
		}
	case reflect.Invalid:
		fmt.Fprint(e, "0:~")
	case reflect.Ptr:
		v = v.Elem()
		return e.Encode(v.Interface())
	case reflect.Map:
		var buf bytes.Buffer
		f := Encoder{
			Writer: &buf,
		}
		ks := v.MapKeys()
		sort.Slice(ks, func(i, j int) bool {
			return ks[i].String() < ks[j].String()
		})
		for _, k := range ks {
			f.Encode(k.Interface())
			f.Encode(v.MapIndex(k).Interface())
		}
		fmt.Fprintf(e, "%d:%s}", buf.Len(), buf.Bytes())
	case reflect.Struct:
		var buf bytes.Buffer
		f := Encoder{
			Writer: &buf,
		}
		for i := 0; i < v.NumField(); i++ {
			t := v.Type().Field(i)
			v := v.Field(i)

			if !v.CanInterface() {
				continue
			}

			name := t.Name
			if tag, ok := t.Tag.Lookup("tnetstrings"); ok {
				if tag == "-" {
					continue
				}
				ts := strings.Split(tag, ",")
				if len(ts) > 0 && ts[0] != "" {
					name = ts[0]
				}
				if len(ts) > 1 && ts[1] == "omitempty" && v == reflect.Zero(t.Type) {
					continue
				}
			}

			f.Encode(name)
			f.Encode(v.Interface())
		}
		fmt.Fprintf(e, "%d:%s}", buf.Len(), buf.Bytes())
	case reflect.Array, reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			s := fmt.Sprintf("%s", val)
			fmt.Fprintf(e, "%d:%s,", len(s), s)
			return nil
		}
		var buf bytes.Buffer
		f := Encoder{
			Writer: &buf,
		}
		for i := 0; i < v.Len(); i++ {
			f.Encode(v.Index(i).Interface())
		}
		fmt.Fprintf(e, "%d:%s]", buf.Len(), buf.Bytes())
	default:
		return UnsupportedType{Kind: v.Kind()}
	}
	return nil
}

type UnsupportedType struct {
	reflect.Kind
}

func (e UnsupportedType) Error() string {
	return fmt.Sprintf("unsupported type: %v", e.Kind)
}

var NonStringKey = errors.New("non string key")
