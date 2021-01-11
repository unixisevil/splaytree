package splaytree

import (
	"testing"
)

type Int int

func (i Int) Compare(other Comparable) int {
	switch ov := other.(type) {
	case PositiveInfinity:
		return -1
	case Int:
		if int(i) < int(ov) {
			return -1
		} else if int(i) == int(ov) {
			return 0
		} else {
			return 1
		}
	}
	return 1
}

func TestInsert(t *testing.T) {
	st := New()
	for i := 1; i <= 10; i++ {
		st.Insert(Int(i))
	}
	if expect, ret := Int(10), st.root.key.(Int); expect != ret {
		t.Errorf("expect root elem is 10, but got %v", ret)
	}
	t.Logf("splaytree:\n%s\n", st.String())

	for w := st.root; w != nil; w = w.links[Left] {
		t.Logf("data:  %v", w.key)
	}

	if expect, ret := true, st.Exist(Int(5)); expect != ret {
		t.Error("elem 5 should exist")
	}
	t.Logf("splaytree:\n%s\n", st.String())
}

func TestDelete(t *testing.T) {
	st := New()
	for i := 1; i <= 20; i += 2 {
		st.Insert(Int(i))
	}
	t.Logf("before delete:\n%s\n", st.String())

	st.Delete(Int(9))

	if expect, ret := false, st.Exist(Int(9)); expect != ret {
		t.Error("elem 9 should not exist")
	}
	t.Logf("after delete:\n%s\n", st.String())
}
