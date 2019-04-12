package joystick

import (
	"fmt"
	"github.com/simulatedsimian/joystick"
)

func init() {
	js, err := joystick.Open(0)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Joystick Name: %s", js.Name())
	fmt.Printf("   Axis Count: %d", js.AxisCount())
	fmt.Printf(" Button Count: %d", js.ButtonCount())

	state, err := js.Read()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Axis Data: %v", state.AxisData)
	js.Close()
	panic("DONE")
}
