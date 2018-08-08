package main

import (
	"fmt"
	"time"

	lua "github.com/yuin/gopher-lua"

	"midi"
	"vjoy"
)

const luaVJOyTypeName = "person"

// Registers my person type to given L.
func registerVJoy(L *lua.LState) {
	mt := L.NewTypeMetatable(luaVJOyTypeName)
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"SetButton": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			if vj, ok := ud.Value.(*vjoy.VJoy); !ok {
				L.ArgError(1, "vjoy expected")
				return 0
			} else {
				if L.GetTop() == 3 {
					btn := L.CheckInt(2)
					state := L.CheckInt(3)
					vj.SetButton(btn, state)
					return 0
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
				return 0
			} else {
				if L.GetTop() == 3 {
					axis := L.CheckInt(2)
					pos := L.CheckInt(3)
					vj.SetAxis(axis, pos)
					return 0
				} else {
					L.ArgError(1, "two arguments expected")
				}
			}
			return 0
		},
	}))
}

// Constructor
func userdataVJoy(L *lua.LState, vj *vjoy.VJoy) lua.LValue {
	ud := L.NewUserData()
	ud.Value = vj
	L.SetMetatable(ud, L.GetTypeMetatable(luaVJOyTypeName))
	return ud
}

// Constructor
// midi.in = {
//	channel
//	status
//	data = [
//		0,
//		1
//	]
// }
// midi.out = {
// 	send(channel, status, hb, lb)
// }
func userdataMIDI(L *lua.LState, m *midi.MIDI) lua.LValue {
	t := L.NewTable()
	t.RawSet(lua.LString("in"), lua.LNil)
	return t
}

func main() {
	midis, err := midi.OpenDevices()
	if err != nil {
		panic(err)
	}

	vjoys := make([]*vjoy.VJoy, 1)
	vj, err := vjoy.OpenVJoy(0)
	if err != nil {
		panic(err)
	}

	vjoys[0] = vj

	quit := make(chan int)
	go func() {
		ft := 1000 * time.Millisecond / 60.0
		toggle := 0
		L := lua.NewState()
		defer L.Close()

		tick := 0

		registerVJoy(L)
		vjoy.Lua(L)
		vjoystable := L.NewTable()
		for _, vj := range vjoys {
			vjoystable.Append(userdataVJoy(L, vj))
		}

		L.SetGlobal("starting", lua.LBool(true))

		miditable := L.NewTable()
		for _, m := range midis.Devices {
			miditable.Append(userdataMIDI(L, m))
		}
		L.SetGlobal("midi", miditable)
		L.SetGlobal("vjoys", vjoystable)
		L.SetGlobal("tick", lua.LNumber(float64(tick)))

		if err := L.DoFile("hello.lua"); err != nil {
			panic(err)
		}

		for {
			select {
			case <-quit:
				fmt.Println("quit")
				return
			case ev := <-midis.Channel:
				m := miditable.RawGetInt(ev.DeviceID)
				if m == lua.LNil {
					continue
				}
				mt := m.(*lua.LTable)
				t := L.NewTable()
				dt := L.NewTable()
				dt.RawSetInt(1, lua.LNumber(ev.Data1))
				dt.RawSetInt(2, lua.LNumber(ev.Data2))
				t.RawSetString("channel", lua.LNumber(ev.Status&0x0f))
				t.RawSetString("status", lua.LNumber(ev.Status&0xf0))
				t.RawSetString("data", dt)
				mt.RawSet(lua.LString("in"), t)
			default:
				L.SetGlobal("tick", lua.LNumber(float64(tick)))
				//vj.SetButton(vjoy.BTN_A, toggle)
				if toggle == 1 {
					toggle = 0
				} else {
					toggle = 1
				}
				if err := L.DoFile("hello.lua"); err != nil {
					panic(err)
				}
				vj.Tick()
				tick++
				if tick == 1 {
					L.SetGlobal("starting", lua.LBool(false))
				}
				for i := 0; i < len(midis.Devices); i++ {
					if m, ok := miditable.RawGetInt(i).(*lua.LTable); ok == true {
						m.RawSet(lua.LString("in"), lua.LNil)
					}
				}
				time.Sleep(ft)
			}
		}
	}()

	fmt.Println("Press enter to exit")
	var input string
	fmt.Scanln(&input)
	fmt.Print(input)
	quit <- 0

	vj.Close()
}
