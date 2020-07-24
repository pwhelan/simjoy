package vjoyengine

import (
	"time"

	"github.com/eapache/channels"
	"github.com/pwhelan/simjoy/vjoy"
	lua "github.com/yuin/gopher-lua"
)

const vJoyName = "vjoy"

// VJoyEngine engine for virtual joysticks
type VJoyEngine struct {
	tick    int64
	vjoys   []*vjoy.VJoy
	channel <-chan time.Time
}

// Register VJoy type
func Register(L *lua.LState) (*VJoyEngine, error) {
	vjoys := make([]*vjoy.VJoy, 1)
	vj, err := vjoy.OpenVJoy(0)
	if err != nil {
		panic(err)
	}

	vjoys[0] = vj

	mt := L.NewTypeMetatable(vJoyName)
	// methods
	L.SetField(mt,
		"__index",
		L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
			"SetButton": func(L *lua.LState) int {
				ud := L.CheckUserData(1)
				if vj, ok := ud.Value.(*vjoy.VJoy); !ok {
					L.ArgError(1, "vjoy expected")
				} else {
					if L.GetTop() == 3 {
						btn := L.CheckInt(2)
						state := L.CheckInt(3)
						vj.SetButton(btn, state)
					} else {
						L.ArgError(1, "two arguments expected")
					}
				}
				return 0
			},
			"SetAxis": func(L *lua.LState) int {
				ud := L.CheckUserData(1)
				if vj, ok := ud.Value.(*vjoy.VJoy); !ok {
					L.ArgError(1, "vjoy expected")
				} else {
					if L.GetTop() == 3 {
						axis := L.CheckInt(2)
						pos := L.CheckInt(3)
						vj.SetAxis(axis, pos)
					} else {
						L.ArgError(1, "two arguments expected")
					}
				}
				return 0
			},
		}))

	vjoy.Lua(L)
	vjoystable := L.NewTable()
	for _, vj := range vjoys {
		vjoystable.Append(userdataVJoy(L, vj))
	}
	L.SetGlobal("vjoys", vjoystable)

	ft := (1000 * time.Millisecond) / 60.0
	ticker := time.NewTicker(ft)

	return &VJoyEngine{tick: 0, vjoys: vjoys, channel: ticker.C}, nil
}

// Channel for receiving ticks
func (engine *VJoyEngine) Channel() channels.SimpleOutChannel {
	return channels.Wrap(engine.channel)
}

// Tick the engine once
func (engine *VJoyEngine) Tick(L *lua.LState, data interface{}) {
	tickfn := L.GetGlobal("tick")
	if tickfn != nil {
		if err := L.CallByParam(lua.P{
			Fn:      tickfn,
			NRet:    0,
			Protect: true,
		}, lua.LNumber(engine.tick)); err != nil {
			panic(err)
		}
	}
	for _, vj := range engine.vjoys {
		vj.Tick()
	}
	engine.tick++
}

// Constructor
func userdataVJoy(L *lua.LState, vj *vjoy.VJoy) lua.LValue {
	ud := L.NewUserData()
	ud.Value = vj
	L.SetMetatable(ud, L.GetTypeMetatable(vJoyName))
	return ud
}
