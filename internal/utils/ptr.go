package utils

import "reflect"

func Ref[T any](val T) *T {
	return &val
}

func NilRefIfZero[T any](val T) *T {
	if reflect.ValueOf(val).IsZero() {
		return nil
	}

	return &val
}
