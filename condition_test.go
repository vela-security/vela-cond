package cond

import (
	"strconv"
	"testing"
)

type Event struct {
	Type  string
	Value int
}

func (ev *Event) Field(key string) string {
	switch key {
	case "type":
		return ev.Type
	case "value":
		return strconv.Itoa(ev.Value)
	}

	return ""
}

func TestExp(t *testing.T) {
	cnd := New("value eq 123,456,789")
	ev := &Event{
		Type:  "typeof",
		Value: 456,
	}

	t.Log(cnd.Match(ev))
}
