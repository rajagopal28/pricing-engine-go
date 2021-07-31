package util

import (
  "log"
  "testing"
  "reflect"
)


func TestAssert(t *testing.T) {
  log.Println("Some testing here!")
  AssertEqual(1, 1, t)
  AssertEqual([1]int{1}, [1]int{1}, t)
}

func AssertEqual(a interface{}, b interface{}, t *testing.T) {
  if reflect.DeepEqual(a, b) {
    return
  }
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}
