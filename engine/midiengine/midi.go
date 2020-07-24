package midiengine

import (
	"fmt"

	"github.com/eapache/channels"
	"github.com/pwhelan/simjoy/midi"
	lua "github.com/yuin/gopher-lua"
)

const midiTypeName = "midi"

// MIDIEngine for MIDI!
type MIDIEngine struct {
	devices   *midi.MIDIS
	miditable *lua.LTable
}

// Register MIDI Engine
func Register(L *lua.LState) (*MIDIEngine, error) {
	midis, err := midi.OpenDevices()
	if err != nil {
		panic(err)
	}

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
	miditable := L.NewTable()
	for idx := range midis.Devices {
		m := midis.Devices[idx]
		func(m *midi.MIDI) {
			fmt.Printf("Added MIDI device: %d:%s\n", m.ID, m.Info.Name)
			miditable.RawSetInt(int(m.ID), userdataMIDI(L, m))
		}(m)
	}
	L.SetGlobal("midi", miditable)

	return &MIDIEngine{devices: midis, miditable: miditable}, nil
}

// Channel for receiving incoming events
func (engine *MIDIEngine) Channel() channels.SimpleOutChannel {
	return channels.Wrap(engine.devices.Channel)
}

// Tick the MIDI engine off...
func (engine *MIDIEngine) Tick(L *lua.LState, data interface{}) {
	ev := data.(midi.Event)
	fmt.Printf("DEVICE-ID: %d\n", ev.DeviceID)
	m := engine.miditable.RawGetInt(ev.DeviceID)
	if m == lua.LNil {
		fmt.Println("BAD DEVICE")
		return
	}
	mt := m.(*lua.LTable)
	t := L.NewTable()
	dt := L.NewTable()
	t.RawSetString("channel", lua.LNumber(ev.Status&0x0f))
	t.RawSetString("status", lua.LNumber(ev.Status&0xf0))
	dt.RawSetInt(1, lua.LNumber(ev.Data1))
	dt.RawSetInt(2, lua.LNumber(ev.Data2))
	t.RawSetString("data", dt)
	mt.RawSet(lua.LString("in"), t)
	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("midirecv"),
		NRet:    0,
		Protect: true,
	}, mt, t); err != nil {
		panic(err)
	}
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
