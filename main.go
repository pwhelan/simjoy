package main

import (
	"fmt"
	"time"

	lua "github.com/yuin/gopher-lua"

	"github.com/pwhelan/simjoy/midi"
	"github.com/pwhelan/simjoy/vjoy"
)

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

	go func() {
		ft := (1000 * time.Millisecond) / 60.0
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

		miditable := L.NewTable()
		for idx := range midis.Devices {
			m := midis.Devices[idx]
			fmt.Printf("Added MIDI device: %s\n", m.Info.Name)
			func(m *midi.MIDI) {
				miditable.Append(userdataMIDI(L, m))
			}(m)
		}
		L.SetGlobal("midi", miditable)
		L.SetGlobal("vjoys", vjoystable)

		fmt.Println("loading ...")
		err := L.DoFile("hello.lua")
		if err != nil {
			panic(err)
		}
		fmt.Println("loaded")

		if err := L.CallByParam(lua.P{
			Fn:      L.GetGlobal("init"),
			NRet:    0,
			Protect: true,
		}); err != nil {
			panic(err)
		}
		fmt.Println("called")

		ticker := time.NewTicker(ft)
		t := L.NewTable()
		dt := L.NewTable()

		for {
			select {
			case ev := <-midis.Channel:
				m := miditable.RawGetInt(ev.DeviceID)
				if m == lua.LNil {
					fmt.Println("BAD DEVICE")
					continue
				}
				mt := m.(*lua.LTable)
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
			case <-ticker.C:
				if err := L.CallByParam(lua.P{
					Fn:      L.GetGlobal("tick"),
					NRet:    0,
					Protect: true,
				}, lua.LNumber(tick)); err != nil {
					panic(err)
				}
				vj.Tick()
				tick++
			}
		}
	}()

	for {
		time.Sleep(10 * time.Second)
	}
	vj.Close()
}
