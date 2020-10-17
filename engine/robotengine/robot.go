package robotengine

import (
	"fmt"

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
				L.CheckUserData(1)
				fmt.Printf("ARGS=%d\n", L.GetTop())
				key := L.CheckString(2)

				if (L.GetTop() - 2) > 0 {
					modnum := L.GetTop() - 2
					mods := make([]interface{}, modnum)

					for i := 3; i <= L.GetTop(); i++ {
						mods[i-3] = L.CheckString(i)
					}
					robotgo.KeyTap(key, mods...)
				} else {
					robotgo.KeyTap(key)
				}

				return 0
			},
		}))
	tkbd := L.NewUserData()
	L.SetMetatable(tkbd, L.GetTypeMetatable("keyboard"))
	L.SetGlobal("keyboard", tkbd)
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
