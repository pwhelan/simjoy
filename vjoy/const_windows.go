package vjoy

import (
	"github.com/tajtiattila/vjoy"
	lua "github.com/yuin/gopher-lua"
)

const (
	ABS_X       int = vjoy.AxisX
	ABS_Y           = AxisY
	ABS_Z           = AxisZ
	ABS_RX          = AxisRX
	ABS_RY          = AxisRY
	ABS_RZ          = AxisRZ
	BTN_TRIGGER     = 0
	BTN_THUMB       = 1
	BTN_THUMB2      = 2
	BTN_TOP         = 3
	BTN_TOP2        = 4
	BTN_PINKIE      = 5
	BTN_BASE        = 6
	BTN_BASE2       = 7
	BTN_BASE3       = 8
	BTN_BASE4       = 9
	BTN_BASE5       = 10
	BTN_BASE6       = 11
	BTN_DEAD        = 12
	BTN_A           = 13
	BTN_B           = 14
	BTN_C           = 15
	BTN_X           = 16
	BTN_Y           = 17
	BTN_Z           = 18
	BTN_TL          = 19
	BTN_TR          = 20
	BTN_TL2         = 21
	BTN_TR2         = 22
	BTN_SELECT      = 23
	BTN_START       = 24
	BTN_MODE        = 25
	BTN_THUMBL      = 26
	BTN_THUMBR      = 27
)

func Lua(L *lua.LState) {
	t := L.NewTable()
	t.RawSetString("ABS_X", lua.LNumber(ABS_X))
	t.RawSetString("ABS_Y", lua.LNumber(ABS_Y))
	t.RawSetString("ABS_Z", lua.LNumber(ABS_Z))
	t.RawSetString("ABS_RX", lua.LNumber(ABS_RX))
	t.RawSetString("ABS_RY", lua.LNumber(ABS_RY))
	t.RawSetString("ABS_RZ", lua.LNumber(ABS_RZ))
	t.RawSetString("BTN_TRIGGER", lua.LNumber(BTN_TRIGGER))
	t.RawSetString("BTN_THUMB", lua.LNumber(BTN_THUMB))
	t.RawSetString("BTN_THUMB2", lua.LNumber(BTN_THUMB2))
	t.RawSetString("BTN_TOP", lua.LNumber(BTN_TOP))
	t.RawSetString("BTN_TOP2", lua.LNumber(BTN_TOP2))
	t.RawSetString("BTN_PINKIE", lua.LNumber(BTN_PINKIE))
	t.RawSetString("BTN_BASE", lua.LNumber(BTN_BASE))
	t.RawSetString("BTN_BASE2", lua.LNumber(BTN_BASE2))
	t.RawSetString("BTN_BASE3", lua.LNumber(BTN_BASE3))
	t.RawSetString("BTN_BASE4", lua.LNumber(BTN_BASE4))
	t.RawSetString("BTN_BASE5", lua.LNumber(BTN_BASE5))
	t.RawSetString("BTN_BASE6", lua.LNumber(BTN_BASE6))
	t.RawSetString("BTN_DEAD", lua.LNumber(BTN_DEAD))
	t.RawSetString("BTN_A", lua.LNumber(BTN_A))
	t.RawSetString("BTN_B", lua.LNumber(BTN_B))
	t.RawSetString("BTN_C", lua.LNumber(BTN_C))
	t.RawSetString("BTN_X", lua.LNumber(BTN_X))
	t.RawSetString("BTN_Y", lua.LNumber(BTN_Y))
	t.RawSetString("BTN_Z", lua.LNumber(BTN_Z))
	t.RawSetString("BTN_TL", lua.LNumber(BTN_TL))
	t.RawSetString("BTN_TR", lua.LNumber(BTN_TR))
	t.RawSetString("BTN_TL2", lua.LNumber(BTN_TL2))
	t.RawSetString("BTN_TR2", lua.LNumber(BTN_TR2))
	t.RawSetString("BTN_SELECT", lua.LNumber(BTN_SELECT))
	t.RawSetString("BTN_START", lua.LNumber(BTN_START))
	t.RawSetString("BTN_MODE", lua.LNumber(BTN_MODE))
	t.RawSetString("BTN_THUMBL", lua.LNumber(BTN_THUMBL))
	t.RawSetString("BTN_THUMBR", lua.LNumber(BTN_THUMBR))
	L.SetGlobal("joystick", t)
}
