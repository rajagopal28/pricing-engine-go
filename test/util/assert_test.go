package util


import (
  "log"
  "testing"
)


func TestAssert(t *testing.T) {
  log.Println("Some testing here!")
  AssertEqual(1, 1, t)
  AssertEqual([1]int{1}, [1]int{1}, t)
  AssertTrue(true, t)
  AssertFalse(false, t)
}
