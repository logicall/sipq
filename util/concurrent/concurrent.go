package concurrent

import (
	"sync"
)

type Pair struct {
	First, Second interface{}
}

type List struct {
	sync.Mutex
	items []interface{}
}

func NewList() *List {
	return &List{}
}

func (l *List) Add(item interface{}) {
	l.Lock()
	defer l.Unlock()

	l.items = append(l.items, item)
}

func (l *List) Len() int {
	l.Lock()
	defer l.Unlock()

	return len(l.items)
}

//0 <= ith < Len
func (l *List) Remove(ith int) {
	l.Lock()
	defer l.Unlock()

	l.remove(ith)

}

func (l *List) remove(ith int) {
	if ith < 0 || ith >= len(l.items) {
		return
	}
	l.items = append(l.items[:ith], l.items[ith+1:]...)

}

func (l *List) RemoveItem(item interface{}) {
	l.Lock()
	defer l.Unlock()

	i, ok := l.find(item)
	if !ok {
		return
	}
	l.remove(i)
}

func (l *List) find(item interface{}) (index int, ok bool) {
	for i := 0; i < len(l.items); i++ {
		if item == l.items[i] {
			index = i
			ok = true
			return
		}
	}
	return -1, false
}

func (l *List) Find(item interface{}) (index int, ok bool) {
	l.Lock()
	defer l.Unlock()
	return l.find(item)
}

func (l *List) FindItemBy(predicate func(item interface{}) bool) (result interface{}, ok bool) {
	idx, found := l.findBy(predicate)
	if !found {
		return nil, false
	}
	return l.items[idx], true
}

func (l *List) FindBy(predicate func(item interface{}) bool) (index int, ok bool) {
	l.Lock()
	defer l.Unlock()

	return l.findBy(predicate)
}

func (l *List) findBy(predicate func(item interface{}) bool) (index int, ok bool) {
	for i := 0; i < len(l.items); i++ {
		if predicate(l.items[i]) {
			index = i
			ok = true
			return
		}
	}
	return -1, false
}

//The user uses range to iter.
func (l *List) Iter() <-chan interface{} {
	l.Lock()

	c := make(chan interface{})

	go func() {
		defer l.Unlock()
		defer close(c)

		for _, v := range l.items {
			c <- v
		}

	}()

	return c
}

//The user uses range to iter.
func (l *List) IterPair() <-chan Pair {
	l.Lock()

	c := make(chan Pair)

	go func() {
		defer l.Unlock()
		defer close(c)

		for i, v := range l.items {
			c <- Pair{First: i, Second: v}
		}

	}()

	return c
}

func (l *List) Get(ith int) interface{} {
	l.Lock()
	defer l.Unlock()

	return l.items[ith]
}

type Map struct {
	sync.Mutex
	items map[interface{}]interface{}
}
