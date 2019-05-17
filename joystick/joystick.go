package joystick

import (
	"fmt"
	"time"

	"github.com/simulatedsimian/joystick"
	lua "github.com/yuin/gopher-lua"
)

const (
	// EventButton is a button event
	EventButton = iota
	// EventAxis is an axis event
	EventAxis
)

// Event represents a joystick Event, either button or axis.
type Event interface{}

// EventHeader is the common header for joystick events.
type EventHeader struct {
	Type       int
	JoystickID int
}

// AxisEvent holds the axis event data.
type AxisEvent struct {
	EventHeader
	Axis  int
	Value int
}

// ButtonEvent holds the button event data.
type ButtonEvent struct {
	EventHeader
	Button int
	Value  bool
}

// OpenDevices opens all joysticks and returns either a channel for all joystick events
// or nil and an error
func OpenDevices() (chan Event, error) {
	ch := make(chan Event)
	for jsid := 0; jsid < 4; jsid++ {
		js, err := joystick.Open(jsid)
		if err != nil {
			continue
		}
		fmt.Printf("Name: %s\n", js.Name())
		fmt.Printf("\tAxis Count: %d\n", js.AxisCount())
		fmt.Printf("\tButton Count: %d\n", js.ButtonCount())

		go func(jsid int, js joystick.Joystick) {
			axis := make([]int, js.AxisCount())
			buttons := uint32(0)

			tick := time.NewTicker(16 * time.Millisecond)
			for {
				<-tick.C
				state, err := js.Read()
				if err != nil {
					panic(err)
				}
				for i := 0; i < js.AxisCount(); i++ {
					if state.AxisData[i] != axis[i] {
						ch <- AxisEvent{
							EventHeader{
								Type:       EventAxis,
								JoystickID: jsid,
							},
							i,
							state.AxisData[i],
						}
						axis[i] = state.AxisData[i]
					}
				}
				for i := 0; i < js.ButtonCount(); i++ {
					btnbit := uint32(1 << uint(i))
					if (btnbit & state.Buttons) != (btnbit & buttons) {
						ch <- ButtonEvent{
							EventHeader{
								Type:       EventButton,
								JoystickID: jsid,
							},
							i,
							(btnbit & state.Buttons) != 0,
						}
					}
				}
				buttons = state.Buttons
			}
		}(jsid, js)
	}
	return ch, nil
}

// Lua exposes the joystick constants
func Lua(l *lua.LState) {
	t := l.GetGlobal("joystick").(*lua.LTable)
	t.RawSetString("EVENT_BUTTON", lua.LNumber(EventButton))
	t.RawSetString("EVENT_AXIS", lua.LNumber(EventAxis))
}
