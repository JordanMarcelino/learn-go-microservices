package validationutils

import (
	"reflect"
	"strings"
)

func TagNameFormatter(fld reflect.StructField) string {
	var name string

	jsonTag := fld.Tag.Get("json")
	formTag := fld.Tag.Get("form")

	if jsonTag != "" {
		name = strings.SplitN(jsonTag, ",", 2)[0]
	} else if formTag != "" {
		name = strings.SplitN(formTag, ",", 2)[0]
	}

	if name == "-" {
		return ""
	}
	return name
}
