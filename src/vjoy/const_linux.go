package vjoy
import (
	lua "github.com/yuin/gopher-lua"
)
const (
	BTN_TRIGGER = 288
	BTN_THUMB = 289
	BTN_THUMB2 = 290
	BTN_TOP = 291
	BTN_TOP2 = 292
	BTN_PINKIE = 293
	BTN_BASE = 294
	BTN_BASE2 = 295
	BTN_BASE3 = 296
	BTN_BASE4 = 297
	BTN_BASE5 = 298
	BTN_BASE6 = 299
	BTN_DEAD = 303
	BTN_GAMEPAD = 304
	BTN_SOUTH = 304
	BTN_A = 304
	BTN_EAST = 305
	BTN_B = 305
	BTN_C = 306
	BTN_NORTH = 307
	BTN_X = 307
	BTN_WEST = 308
	BTN_Y = 308
	BTN_Z = 309
	BTN_TL = 310
	BTN_TR = 311
	BTN_TL2 = 312
	BTN_TR2 = 313
	BTN_SELECT = 314
	BTN_START = 315
	BTN_MODE = 316
	BTN_THUMBL = 317
	BTN_THUMBR = 318
)
func Lua(L *lua.LState) {
	t := L.NewTable()
	t.RawSetString("BTN_TRIGGER", lua.LNumber(288))
	t.RawSetString("BTN_THUMB", lua.LNumber(289))
	t.RawSetString("BTN_THUMB2", lua.LNumber(290))
	t.RawSetString("BTN_TOP", lua.LNumber(291))
	t.RawSetString("BTN_TOP2", lua.LNumber(292))
	t.RawSetString("BTN_PINKIE", lua.LNumber(293))
	t.RawSetString("BTN_BASE", lua.LNumber(294))
	t.RawSetString("BTN_BASE2", lua.LNumber(295))
	t.RawSetString("BTN_BASE3", lua.LNumber(296))
	t.RawSetString("BTN_BASE4", lua.LNumber(297))
	t.RawSetString("BTN_BASE5", lua.LNumber(298))
	t.RawSetString("BTN_BASE6", lua.LNumber(299))
	t.RawSetString("BTN_DEAD", lua.LNumber(303))
	t.RawSetString("BTN_GAMEPAD", lua.LNumber(304))
	t.RawSetString("BTN_SOUTH", lua.LNumber(304))
	t.RawSetString("BTN_A", lua.LNumber(304))
	t.RawSetString("BTN_EAST", lua.LNumber(305))
	t.RawSetString("BTN_B", lua.LNumber(305))
	t.RawSetString("BTN_C", lua.LNumber(306))
	t.RawSetString("BTN_NORTH", lua.LNumber(307))
	t.RawSetString("BTN_X", lua.LNumber(307))
	t.RawSetString("BTN_WEST", lua.LNumber(308))
	t.RawSetString("BTN_Y", lua.LNumber(308))
	t.RawSetString("BTN_Z", lua.LNumber(309))
	t.RawSetString("BTN_TL", lua.LNumber(310))
	t.RawSetString("BTN_TR", lua.LNumber(311))
	t.RawSetString("BTN_TL2", lua.LNumber(312))
	t.RawSetString("BTN_TR2", lua.LNumber(313))
	t.RawSetString("BTN_SELECT", lua.LNumber(314))
	t.RawSetString("BTN_START", lua.LNumber(315))
	t.RawSetString("BTN_MODE", lua.LNumber(316))
	t.RawSetString("BTN_THUMBL", lua.LNumber(317))
	t.RawSetString("BTN_THUMBR", lua.LNumber(318))
	L.SetGlobal("joystick", t)
}
