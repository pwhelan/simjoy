package vjoy

//go:generate ../../bin/vjoygenerate

import (
	uinput "github.com/ynsta/uinput"
)

type vjoyLinuxHandler struct {
	wd uinput.WriteDevice
	ui uinput.UInput
}

func OpenVJoy(index uint) (*VJoy, error) {

	var desc vjoyLinuxHandler

	desc.wd.Open()

	if err := desc.ui.Init(
		&desc.wd,
		"Microsoft X-Box 360 pad",
		0x045e,
		0x028e,
		0x110,
		[]uinput.EventCode{uinput.BTN_START,
			uinput.BTN_MODE,
			uinput.BTN_SELECT,
			uinput.BTN_A,
			uinput.BTN_B,
			uinput.BTN_X,
			uinput.BTN_Y,
			uinput.BTN_TL,
			uinput.BTN_TR,
			uinput.BTN_THUMBL,
			uinput.BTN_THUMBR,
		},
		[]uinput.EventCode{},
		[]uinput.AxisSetup{
			{Code: uinput.ABS_X, Min: -32768, Max: 32767, Fuzz: 16, Flat: 128},
			{Code: uinput.ABS_Y, Min: -32768, Max: 32767, Fuzz: 16, Flat: 128},
			{Code: uinput.ABS_RX, Min: -32768, Max: 32767, Fuzz: 16, Flat: 128},
			{Code: uinput.ABS_RY, Min: -32768, Max: 32767, Fuzz: 16, Flat: 128},
			{Code: uinput.ABS_Z, Min: 0, Max: 255, Fuzz: 0, Flat: 0},
			{Code: uinput.ABS_RZ, Min: 0, Max: 255, Fuzz: 0, Flat: 0},
			{Code: uinput.ABS_HAT0X, Min: -1, Max: 1, Fuzz: 0, Flat: 0},
			{Code: uinput.ABS_HAT0Y, Min: -1, Max: 1, Fuzz: 0, Flat: 0},
		},
		false,
	); err != nil {
		return nil, err
	}

	return &VJoy{handler: &desc}, nil
}

func (vj *VJoy) SetButton(btn int, status int) {
	desc := vj.handler.(*vjoyLinuxHandler)
	desc.ui.KeyEvent(uinput.EventCode(btn), uinput.EventValue(status))
}

func (vj *VJoy) Tick() {
	desc := vj.handler.(*vjoyLinuxHandler)
	desc.ui.SynEvent()
}

func (vj *VJoy) Close() {
	desc := vj.handler.(*vjoyLinuxHandler)
	desc.wd.Close()
}
