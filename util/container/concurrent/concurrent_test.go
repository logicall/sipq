package concurrent

import (
	"fmt"
	"testing"
)

func TestList(t *testing.T) {
	l := NewList()
	l.Add(0)
	l.Add(1)
	l.Add(2)
	l.Add(3)
	if l.Len() != 4 {
		t.Error("unexpected", l.Len())
	}

	l.Remove(3)
	if l.Len() != 3 {
		t.Error("unexpected", l.Len())
	}

	var values1 []int
	for value := range l.Iter() {
		values1 = append(values1, value.(int))
	}

	var values2 []int
	for value := range l.IterPair() {
		values2 = append(values2, value.Second.(int))
		if values2[value.First.(int)] != values1[value.First.(int)] {
			t.Error("unexpected")
		}
	}
}

func TestMap(t *testing.T) {
	m := NewMap()
	m.Put("0", 0)
	m.Put("1", 1)

	_, ok := m.Get("0")
	if !ok {
		t.Error("not expected")
	}

	_, ok = m.Get("1")
	if !ok {
		t.Error("not expected")
	}

	for p := range m.IterPair() {
		k := p.First.(string)
		v := p.Second.(int)
		if fmt.Sprint(v) != k {
			t.Error("not expected")
		}
	}
}
