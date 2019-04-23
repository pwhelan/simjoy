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

	quit := make(chan int)
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
