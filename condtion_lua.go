package cond

import (
	"github.com/vela-security/vela-public/lua"
	"github.com/vela-security/vela-public/pipe"
)

type CondL struct {
	onMatch *pipe.Px
	noMatch *pipe.Px
	co      *lua.LState
	cnd     *Cond
}

func (c *CondL) String() string                         { return "vela.condition" }
func (c *CondL) Type() lua.LValueType                   { return lua.LTObject }
func (c *CondL) AssertFloat64() (float64, bool)         { return 0, false }
func (c *CondL) AssertString() (string, bool)           { return "", false }
func (c *CondL) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (c *CondL) Peek() lua.LValue                       { return c }

func (c *CondL) okL(L *lua.LState) int {
	c.onMatch = pipe.NewByLua(L, pipe.Seek(0), pipe.Env(xEnv))
	L.Push(c)
	return 1
}

func (c *CondL) noL(L *lua.LState) int {
	c.noMatch = pipe.NewByLua(L, pipe.Seek(0), pipe.Env(xEnv))
	L.Push(c)
	return 1
}

func (c *CondL) match(lv lua.LValue, L *lua.LState) bool {
	if c.cnd.Match(lv) {
		c.onMatch.Do(lv, L, func(err error) {
			xEnv.Errorf("condition match function pipe fail %v", err)
		})

		return true
	}

	c.noMatch.Do(lv, L, func(err error) {
		xEnv.Errorf("condition not match function pipe fail %v", err)
	})
	return false
}

func (c *CondL) matchL(L *lua.LState) int {
	n := L.GetTop()
	if n == 0 {
		L.Push(lua.LFalse)
		return 1
	}

	for i := 1; i <= n; i++ {
		if c.match(L.Get(i), L) {
			L.Push(lua.LTrue)
			return 1
		}
	}

	L.Push(lua.LFalse)
	return 1
}

func (c *CondL) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "ok":
		return lua.NewFunction(c.okL)
	case "no":
		return lua.NewFunction(c.noL)
	case "match":
		return lua.NewFunction(c.matchL)
	}

	return lua.LNil
}

func (c *CondL) NewIndex(L *lua.LState, key string, val lua.LValue) {
	switch key {
	case "ok":
		c.onMatch = pipe.New(pipe.Env(xEnv))
		c.onMatch.LValue(val)
	case "no":
		c.onMatch = pipe.New(pipe.Env(xEnv))
		c.onMatch.LValue(val)
	}

}

func newCondL(L *lua.LState) *CondL {
	return &CondL{
		onMatch: pipe.New(pipe.Env(xEnv)),
		noMatch: pipe.New(pipe.Env(xEnv)),
		co:      xEnv.Clone(L),
		cnd:     CheckMany(L, Seek(0)),
	}
}
