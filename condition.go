package cond

import "github.com/vela-security/vela-public/lua"

const (
	or uint8 = iota + 1
	and
)

type Cond struct {
	logic uint8
	data  []*Section
}

func New(c ...string) *Cond {
	n := len(c)
	cond := &Cond{
		logic: and,
		data:  make([]*Section, len(c)),
	}

	if n == 0 {
		return cond
	}

	for i := 0; i < n; i++ {
		cond.data[i] = Compile(c[i])
	}
	return cond
}

func CheckMany(L *lua.LState, opt ...func(*option)) *Cond {
	cnd := &Cond{logic: and}
	cnd.CheckMany(L, opt...)
	return cnd
}

func Check(L *lua.LState, idx int) *LCond {
	ov := L.CheckObject(idx)

	lc, ok := ov.(*LCond)
	if ok {
		lc.co = xEnv.Clone(L)
		return lc
	}

	L.RaiseError("invalid condition object , got %p", &ov)
	return nil
}

func LValue(L *lua.LState, val lua.LValue) *LCond {
	if val.Type() != lua.LTObject {
		L.RaiseError("invalid condition type , got %v", val.Type().String())

	}

	lc, ok := val.(*LCond)
	if ok {
		lc.co = xEnv.Clone(L)
		return lc
	}

	L.RaiseError("invalid condition object , got %p", &val)
	return nil
}
