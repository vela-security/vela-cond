package cond

import (
	"fmt"
	"github.com/vela-security/vela-public/lua"
)

type option struct {
	seek int
	peek Peek
	co   *lua.LState
}

func Seek(i int) func(*option) {
	return func(o *option) {
		o.seek = i
	}
}

func WithCo(co *lua.LState) func(*option) {
	return func(ov *option) {
		ov.co = co
	}
}

func (opt *option) NewPeek(v interface{}) bool {
	switch item := v.(type) {
	case Peek:
		opt.peek = item
		return true

	case interface{ Field(string) string }:
		opt.peek = item.Field
		return true

	case string:
		opt.peek = String(item)
		return true

	case []byte:
		opt.peek = String(string(item))
		return true

	case func() string:
		opt.peek = func(string) string {
			return item()
		}
		return true
	case lua.IndexEx:
		opt.peek = func(key string) string {
			return item.Index(opt.co, key).String()
		}
		return true
	case *lua.LTable:
		opt.peek = func(key string) string {
			return item.RawGetString(key).String()
		}

	case fmt.Stringer:
		opt.peek = String(item.String())
		return true
	}

	return false
}
