package tnetstrings

import (
	"reflect"
	"strings"
)

type tag struct {
	name        string
	displayName string
	omitEmpty   bool
}

func parseTag(f reflect.StructField) *tag {
	t := tag{name: f.Name, displayName: f.Name}
	if tnetstrings, ok := f.Tag.Lookup("tnetstrings"); ok {
		if tnetstrings == "-" {
			return nil
		}
		ts := strings.Split(tnetstrings, ",")
		if len(ts) > 0 && ts[0] != "" {
			t.displayName = ts[0]
		}
		t.omitEmpty = len(ts) > 1 && ts[1] == "omitempty"
	}
	return &t
}
