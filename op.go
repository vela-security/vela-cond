package cond

import (
	"github.com/vela-security/vela-public/auxlib"
	"github.com/vela-security/vela-public/grep"
)

const (
	eq op = iota + 10
	re
	cn
	in
	lt
	le
	ge
	gt
	oop
	pass
)

var (
	opTab = []string{"equal", "grep", "contain", "include", "less", "less or equal", "greater or equal", "greater", "oop"}
)

type op uint8

func (o op) String() string {
	return opTab[(int(o) - 10)]
}

func (o op) Do(a, b string) bool {
	switch o {
	case eq:
		return a == b
	case re:
		return grep.New(b)(a)
	case cn:
		return grep.New(b)(a)
	case in:
		return a == b
	case lt:
		return auxlib.ToFloat64(a) < auxlib.ToFloat64(b)
	case le:
		return auxlib.ToFloat64(a) <= auxlib.ToFloat64(b)
	case ge:
		return auxlib.ToFloat64(a) >= auxlib.ToFloat64(b)
	case gt:
		return auxlib.ToFloat64(a) > auxlib.ToFloat64(b)

	}

	return false
}
