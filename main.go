package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	lua "github.com/yuin/gopher-lua"

	"github.com/pwhelan/simjoy/joystick"

	"github.com/pwhelan/simjoy/engine/midiengine"
	"github.com/pwhelan/simjoy/engine/robotengine"
	"github.com/pwhelan/simjoy/engine/vjoyengine"
)

func start(ctxt context.Context, joysticks chan joystick.Event) {
	L := lua.NewState()
	defer L.Close()

	midi, _ := midiengine.Register(L)
	vj, _ := vjoyengine.Register(L)
	robot, _ := robotengine.Register(L)
	joystick.Lua(L)

	fmt.Printf("loading: %s\n", os.Args[1])
	err := L.DoFile(os.Args[1])
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

	vjchan := vj.Channel().Out()
	midichan := midi.Channel().Out()
	robotchan := robot.Channel().Out()

	for {
		select {
		case <-ctxt.Done():
			fmt.Println("Restart")
			return
		case data := <-vjchan:
			vj.Tick(L, data)
		case data := <-midichan:
			midi.Tick(L, data)
		case data := <-robotchan:
			robot.Tick(L, data)
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
		}
	}
}

func main() {
	ctxtmain, finish := context.WithCancel(context.Background())

	joysticks, err := joystick.OpenDevices()
	if err != nil {
		panic(err)
	}

	ctxteng, restart := context.WithCancel(ctxtmain)
	go start(ctxteng, joysticks)

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
				//vj.Close()
				os.Exit(0)
			case syscall.SIGHUP:
				fmt.Println("HUP")
				restart()
				ctxteng, restart = context.WithCancel(ctxtmain)
				go start(ctxteng, joysticks)
				break
			}
		}
	}
}
