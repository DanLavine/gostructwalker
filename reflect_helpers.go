package gostructwalker

import "reflect"

func pointerDereference(val reflect.Value) reflect.Value {
	for {
		switch val.Kind() {
		case reflect.Pointer:
			return pointerDereference(reflect.Indirect(val))
		default:
			return val
		}
	}
}
