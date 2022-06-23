package cond

import (
	"fmt"
)

type Peek func(string) string

func String(raw string) Peek {
	return func(key string) string {
		return raw
	}
}

func NewPeek(v interface{}) Peek {
	switch item := v.(type) {
	case interface{ Field(string) string }:
		return item.Field
	case string:
		return String(item)
	case []byte:
		return String(string(item))
	case fmt.Stringer:
		return String(item.String())

	case func() string:
		return func(string) string {
			return item()
		}

	case Peek:
		return item
	}

	return nil
}
