package util

import (
	"fmt"
	"reflect"
	"strings"
)

// DiffStruct utility function, used to print failing assertions in tests
func DiffStruct(got, want any) string {
	rv := reflect.ValueOf(got)
	rw := reflect.ValueOf(want)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rw.Kind() == reflect.Ptr {
		rw = rw.Elem()
	}
	rt := rv.Type()
	var sb strings.Builder
	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		fw := rw.Field(i)
		if !reflect.DeepEqual(fv.Interface(), fw.Interface()) {
			fmt.Fprintf(&sb, "  %-26s got %-30s want %s\n",
				rt.Field(i).Name+":",
				derefField(fv),
				derefField(fw),
			)
		}
	}
	return sb.String()
}

func derefField(v reflect.Value) string {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "nil"
		}
		return fmt.Sprintf("%v", v.Elem().Interface())
	}
	return fmt.Sprintf("%v", v.Interface())
}
