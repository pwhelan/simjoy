package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	lua "github.com/yuin/gopher-lua"

	"github.com/pwhelan/simjoy/joystick"
	"github.com/pwhelan/simjoy/midi"
	"github.com/pwhelan/simjoy/vjoy"
)

func run(ctxt context.Context, vjoys []*vjoy.VJoy, midis *midi.MIDIS, joysticks chan joystick.Event) {
	tick := 0
	ft := (1000 * time.Millisecond) / 60.0

	L := lua.NewState()
	defer L.Close()

	registerMIDI(L)
	registerVJoy(L)
	vjoy.Lua(L)
	vjoystable := L.NewTable()
	for _, vj := range vjoys {
		vjoystable.Append(userdataVJoy(L, vj))
	}

	joystick.Lua(L)

	miditable := L.NewTable()
	for idx := range midis.Devices {
		m := midis.Devices[idx]
		func(m *midi.MIDI) {
			fmt.Printf("Added MIDI device: %d:%s\n", m.ID, m.Info.Name)
			miditable.RawSetInt(int(m.ID), userdataMIDI(L, m))
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

	tickfn := L.GetGlobal("tick")
	for {
		select {
		case <-ctxt.Done():
			fmt.Println("Restart")
			return
		case ev := <-midis.Channel:
			fmt.Printf("DEVICE-ID: %d\n", ev.DeviceID)
			m := miditable.RawGetInt(ev.DeviceID)
			if m == lua.LNil {
				fmt.Println("BAD DEVICE")
				continue
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
		case ev := <-joysticks:
			switch ev.(type) {
			case joystick.AxisEvent:
				ca := ev.(joystick.AxisEvent)

				axisevent := L.NewTable()
				axisevent.RawSetString("type", lua.LNumber(ca.Type))
				axisevent.RawSetString("joystickid", lua.LNumber(ca.JoystickID))
				axisevent.RawSetString("axis", lua.LNumber(ca.Axis))
				axisevent.RawSetString("value", lua.LNumber(ca.Value))

				if err := L.CallByParam(lua.P{
					Fn:      L.GetGlobal("joystickrecv"),
					NRet:    0,
					Protect: true,
				}, axisevent); err != nil {
					panic(err)
				}
			case joystick.ButtonEvent:
				//fmt.Println("PUSH THE BUTTON")
				btn := ev.(joystick.ButtonEvent)

				btnevent := L.NewTable()
				btnevent.RawSetString("type", lua.LNumber(btn.Type))
				btnevent.RawSetString("joystickid", lua.LNumber(btn.JoystickID))
				btnevent.RawSetString("button", lua.LNumber(btn.Button))
				btnevent.RawSetString("value", lua.LBool(btn.Value))

				if err := L.CallByParam(lua.P{
					Fn:      L.GetGlobal("joystickrecv"),
					NRet:    0,
					Protect: true,
				}, btnevent); err != nil {
					panic(err)
				}
			}
		case <-ticker.C:
			if tickfn != nil {
				if err := L.CallByParam(lua.P{
					Fn:      tickfn,
					NRet:    0,
					Protect: true,
				}, lua.LNumber(tick)); err != nil {
					panic(err)
				}
			}
			for _, vj := range vjoys {
				vj.Tick()
			}
			tick++
		}
	}
}

func main() {
	ctxtmain, finish := context.WithCancel(context.Background())

	joysticks, err := joystick.OpenDevices()
	if err != nil {
		panic(err)
	}

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

	ctxteng, restart := context.WithCancel(ctxtmain)
	go run(ctxteng, vjoys, midis, joysticks)

	csig := make(chan os.Signal, 1)
	signal.Notify(csig, syscall.SIGTERM, syscall.SIGHUP)

	for {
		select {
		case sig := <-csig:
			fmt.Println("SIGNAL")
			switch sig {
			case syscall.SIGTERM:
				fmt.Println("TERM")
				finish()
				vj.Close()
				os.Exit(0)
			case syscall.SIGHUP:
				fmt.Println("HUP")
				restart()
				ctxteng, restart = context.WithCancel(ctxtmain)
				go run(ctxteng, vjoys, midis, joysticks)
				break
			}
		}
	}
}
