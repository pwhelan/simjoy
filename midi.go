package main

import (
	"github.com/pwhelan/simjoy/midi"
	lua "github.com/yuin/gopher-lua"
)

const midiTypeName = "midi"

// Registers MIDI Type
func registerMIDI(L *lua.LState) {
	miditype := L.NewTypeMetatable(midiTypeName)
	// methods
	L.SetField(miditype,
		"__index",
		L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
			"Send": func(L *lua.LState) int {
				ud := L.CheckUserData(1)
				if m, ok := ud.Value.(*midi.MIDI); !ok {
					L.ArgError(1, "MIDI expected")
				} else {
					if L.GetTop() == 5 {
						channel := L.CheckInt64(2)
						status := L.CheckInt64(3)
						hb := L.CheckInt64(4)
						lb := L.CheckInt64(5)
						m.Send(channel, status, hb, lb)
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
	L.SetMetatable(mud, L.GetTypeMetatable(midiTypeName))
	t.RawSet(lua.LString("out"), mud)
	t.RawSet(lua.LString("id"), lua.LNumber(m.ID))
	t.RawSet(lua.LString("name"), lua.LString(m.Info.Name))
	return t
}
