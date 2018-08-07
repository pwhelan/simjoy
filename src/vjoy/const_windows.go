package vjoy

import (
	lua "github.com/yuin/gopher-lua"
)

const (
	BTN_TRIGGER = 0
	BTN_THUMB   = 1
	BTN_THUMB2  = 2
	BTN_TOP     = 3
	BTN_TOP2    = 4
	BTN_PINKIE  = 5
	BTN_BASE    = 6
	BTN_BASE2   = 7
	BTN_BASE3   = 8
	BTN_BASE4   = 9
	BTN_BASE5   = 10
	BTN_BASE6   = 11
	BTN_DEAD    = 12
	BTN_A       = 13
	BTN_B       = 14
	BTN_C       = 15
	BTN_X       = 16
	BTN_Y       = 17
	BTN_Z       = 18
	BTN_TL      = 19
	BTN_TR      = 20
	BTN_TL2     = 21
	BTN_TR2     = 22
	BTN_SELECT  = 23
	BTN_START   = 24
	BTN_MODE    = 25
	BTN_THUMBL  = 26
	BTN_THUMBR  = 27
)

func Lua(L *lua.LState) {
	t := L.NewTable()
	t.RawSetString("BTN_TRIGGER", lua.LNumber(0))
	t.RawSetString("BTN_THUMB", lua.LNumber(1))
	t.RawSetString("BTN_THUMB2", lua.LNumber(2))
	t.RawSetString("BTN_TOP", lua.LNumber(3))
	t.RawSetString("BTN_TOP2", lua.LNumber(4))
	t.RawSetString("BTN_PINKIE", lua.LNumber(5))
	t.RawSetString("BTN_BASE", lua.LNumber(6))
	t.RawSetString("BTN_BASE2", lua.LNumber(7))
	t.RawSetString("BTN_BASE3", lua.LNumber(8))
	t.RawSetString("BTN_BASE4", lua.LNumber(9))
	t.RawSetString("BTN_BASE5", lua.LNumber(10))
	t.RawSetString("BTN_BASE6", lua.LNumber(11))
	t.RawSetString("BTN_DEAD", lua.LNumber(12))
	t.RawSetString("BTN_A", lua.LNumber(13))
	t.RawSetString("BTN_B", lua.LNumber(14))
	t.RawSetString("BTN_C", lua.LNumber(15))
	t.RawSetString("BTN_X", lua.LNumber(16))
	t.RawSetString("BTN_Y", lua.LNumber(17))
	t.RawSetString("BTN_Z", lua.LNumber(18))
	t.RawSetString("BTN_TL", lua.LNumber(19))
	t.RawSetString("BTN_TR", lua.LNumber(20))
	t.RawSetString("BTN_TL2", lua.LNumber(21))
	t.RawSetString("BTN_TR2", lua.LNumber(22))
	t.RawSetString("BTN_SELECT", lua.LNumber(23))
	t.RawSetString("BTN_START", lua.LNumber(24))
	t.RawSetString("BTN_MODE", lua.LNumber(25))
	t.RawSetString("BTN_THUMBL", lua.LNumber(26))
	t.RawSetString("BTN_THUMBR", lua.LNumber(27))
	L.SetGlobal("joystick", t)
}
