package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v any) *ListItem
	PushBack(v any) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value any
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	front *ListItem
	back  *ListItem
	len   int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}
func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v any) *ListItem {
	item := &ListItem{
		Value: v,
	}
	if l.len == 0 {
		l.front = item
		l.back = item
		item.Next = nil
		item.Prev = nil
	} else {
		l.front.Prev = item
		item.Next = l.front
		item.Prev = nil
		l.front = item
	}
	l.len++
	return item
}

func (l *list) PushBack(v any) *ListItem {
	item := &ListItem{
		Value: v,
	}
	if l.len == 0 {
		l.front = item
		l.back = item
		item.Next = nil
		item.Prev = nil
	} else {
		l.back.Next = item
		item.Prev = l.back
		item.Next = nil
		l.back = item
	}
	l.len++
	return item
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.front = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}
	l.len--
	if l.len == 0 {
		l.front = nil
		l.back = nil
	}

	i.Next = nil
	i.Prev = nil

	if l.len < 0 {
		panic("index out of range")
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil {
		panic("MoveToFront called with nil")
	}
	if i == l.front {
		return
	}
	// отвязать i от текущего места
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if i == l.back {
		l.back = i.Prev
	}
	// вставить i в начало
	i.Prev = nil
	i.Next = l.front
	l.front.Prev = i
	l.front = i
}
