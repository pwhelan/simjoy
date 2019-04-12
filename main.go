package main

import (
	"fmt"
	"time"

	lua "github.com/yuin/gopher-lua"

	"midi"
	"vjoy"
)

const luaVJOyTypeName = "vjoy"
const luaMIDITypeName = "midi"

// Registers VJoy type
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

// Registers MIDI Type
func registerMIDI(L *lua.LState) {
	mt := L.NewTypeMetatable(luaMIDITypeName)
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"Send": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			if m, ok := ud.Value.(*midi.MIDI); !ok {
				L.ArgError(1, "MIDI expected")
				return 0
			} else {
				if L.GetTop() == 5 {
					channel := L.CheckInt64(2)
					status := L.CheckInt64(3)
					hb := L.CheckInt64(4)
					lb := L.CheckInt64(5)
					m.Send(channel, status, hb, lb)
					return 0
				} else {
					L.ArgError(1, "4 arguments expected")
				}
			}
			return 0
		},
	}))
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
	mud := L.NewUserData()
	mud.Value = m
	L.SetMetatable(mud, L.GetTypeMetatable(luaMIDITypeName))
	t.RawSet(lua.LString("out"), mud)
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
		L := lua.NewState()
		defer L.Close()

		tick := 0

		registerMIDI(L)

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

		script, err := L.LoadFile("hello.lua")
		if err != nil {
			panic(err)
		}

		L.Push(script)
		if err := L.PCall(0, lua.MultRet, nil); err != nil {
			panic(err)
		}
		L.SetGlobal("starting", lua.LBool(false))

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
				L.Push(script)
				if err := L.PCall(0, lua.MultRet, nil); err != nil {
					panic(err)
				}
				vj.Tick()
				tick++
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
