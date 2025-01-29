package validationutils

import (
	"reflect"

	"github.com/shopspring/decimal"
)

func DecimalType(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(decimal.Decimal); ok {
		return valuer.String()
	}
	return nil
}
