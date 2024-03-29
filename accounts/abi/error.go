















package abi

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	errBadBool = errors.New("abi: improperly encoded boolean value")
)



func formatSliceString(kind reflect.Kind, sliceSize int) string {
	if sliceSize == -1 {
		return fmt.Sprintf("[]%v", kind)
	}
	return fmt.Sprintf("[%d]%v", sliceSize, kind)
}



func sliceTypeCheck(t Type, val reflect.Value) error {
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return typeErr(formatSliceString(t.GetType().Kind(), t.Size), val.Type())
	}

	if t.T == ArrayTy && val.Len() != t.Size {
		return typeErr(formatSliceString(t.Elem.GetType().Kind(), t.Size), formatSliceString(val.Type().Elem().Kind(), val.Len()))
	}

	if t.Elem.T == SliceTy || t.Elem.T == ArrayTy {
		if val.Len() > 0 {
			return sliceTypeCheck(*t.Elem, val.Index(0))
		}
	}

	if val.Type().Elem().Kind() != t.Elem.GetType().Kind() {
		return typeErr(formatSliceString(t.Elem.GetType().Kind(), t.Size), val.Type())
	}
	return nil
}



func typeCheck(t Type, value reflect.Value) error {
	if t.T == SliceTy || t.T == ArrayTy {
		return sliceTypeCheck(t, value)
	}

	
	if t.GetType().Kind() != value.Kind() {
		return typeErr(t.GetType().Kind(), value.Kind())
	} else if t.T == FixedBytesTy && t.Size != value.Len() {
		return typeErr(t.GetType(), value.Type())
	} else {
		return nil
	}

}


func typeErr(expected, got interface{}) error {
	return fmt.Errorf("abi: cannot use %v as type %v as argument", got, expected)
}
