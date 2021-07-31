package util

import (
  "testing"
  "reflect"
)


func AssertEqual(a interface{}, b interface{}, t *testing.T) {
  if a == b || reflect.DeepEqual(a, b) {
    return
  }
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}
