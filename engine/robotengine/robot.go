package robotengine

import (
	"github.com/eapache/channels"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	lua "github.com/yuin/gopher-lua"
)

// Engine for automating mouse/keyboard
type Engine struct {
	channel chan hook.Event
}

// Register the Robot Engine
func Register(L *lua.LState) (*Engine, error) {
	c := hook.Start()

	lkbd := L.NewTypeMetatable("keyboard")
	L.SetField(lkbd,
		"__index",
		L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
			"Send": func(L *lua.LState) int {
				//ud := L.CheckUserData(1)
				keys := make([]string, L.GetTop()-1)

				for i := 1; i < L.GetTop()-1; i++ {
					keys[i-1] = L.CheckString(i)
				}
				if len(keys) == 1 {
					robotgo.KeyTap(keys[0])
				} else {
					args := make([]interface{}, len(keys)-1)
					for i := 1; i < len(keys); i++ {
						args[i] = keys[i]
					}
					robotgo.KeyTap(keys[0], args...)
				}
				return 0
			},
		}))
	L.SetGlobal("keyboard", lkbd)
	//lmouse := L.NewTypeMetatable("mouse")

	return &Engine{channel: c}, nil
}

// Channel for Robot Events
func (engine *Engine) Channel() channels.SimpleOutChannel {
	return channels.Wrap(engine.channel)
}

// Tick the Robot off!
func (engine *Engine) Tick(L *lua.LState, data interface{}) {

	ev := data.(hook.Event)
	t := L.NewTable()

	switch ev.Kind {
	case hook.MouseDown:
		t.RawSetString("event", lua.LString("mousedown"))
		t.RawSetString("button", lua.LNumber(ev.Button))
	case hook.MouseUp:
		t.RawSetString("event", lua.LString("mouseup"))
		t.RawSetString("button", lua.LNumber(ev.Button))
	default:
		return
	}

	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("mouse"),
		NRet:    0,
		Protect: true,
	}, t); err != nil {
		panic(err)
	}
}
