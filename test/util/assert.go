package util

import (
  "testing"
  "reflect"
)


func AssertEqual(a interface{}, b interface{}, t *testing.T) {
  if reflect.DeepEqual(a, b) {
    return
  }
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}


func AssertTrue(a bool, t *testing.T) {
  if a {
    return
  }
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), !a, reflect.TypeOf(a))
}

func AssertFalse(a bool, t *testing.T) {
  AssertTrue(!a, t)
}
