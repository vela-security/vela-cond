package cond

import (
	"github.com/vela-security/vela-public/assert"
	"github.com/vela-security/vela-public/lua"
)

var xEnv assert.Environment

/*
	cnd = vela.cond("name eq zhangsan,lisi,wangwu").ok(a).no(b)

	hash := "lisi"
	cnd.match("lisi") @end(true)
*/

func newLuaCondition(L *lua.LState) int {
	L.Push(newLCond(L))
	return 1
}

func WithEnv(env assert.Environment) {
	xEnv = env
	xEnv.Set("cond", lua.NewFunction(newLuaCondition))
}
