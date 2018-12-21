package vkbd

//go:generate ../../bin/vkbdgenerate

import (
	//uinput "github.com/ynsta/uinput"
	"github.com/bendahl/uinput"
)

type vkbdLinuxHandler struct {
	kbd uinput.Keyboard
}

func Open() (*VKbd, error) {
	var desc vkbdLinuxHandler

	kbd, err := uinput.CreateKeyboard("/dev/uinput", []byte("simjoy virtual keyboard"))
	if err != nil {
		return nil, err
	}
	desc.kbd = kbd
	return &VKbd{desc: &desc}, nil
}

func (vkbd *VKbd) KeyPress(keycode int) {
	desc, ok := vkbd.desc.(*vkbdLinuxHandler)
	if !ok {
		return
	}
	desc.kbd.KeyPress(keycode)
}

func (vkbd *VKbd) KeyUp() {
	desc, ok := vkbd.desc.(*vkbdLinuxHandler)
	if !ok {
		return
	}
	desc.kbd.KeyUp(uinput.KeyA)
}

func (vkbd *VKbd) KeyDown() {
	desc, ok := vkbd.desc.(*vkbdLinuxHandler)
	if !ok {
		return
	}
	desc.kbd.KeyDown(uinput.KeyA)
}
