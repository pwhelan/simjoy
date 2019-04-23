package vjoy

//go:generate ../../bin/vjoygenerate

import (
	"fmt"
	"github.com/tajtiattila/vjoy"
)

type vjoyWindowsHandler struct {
	dev *vjoy.Device
}

func OpenVJoy(index uint) (*VJoy, error) {
	var desc vjoyWindowsHandler
	if vjoy.Available() == false {
		return nil, fmt.Errorf("vJoy is not available")
	}

	dev, err := vjoy.Acquire(index)
	if err != nil {
		return nil, err
	}
	desc.dev = dev
	dev.Reset()

	return &VJoy{handler: &desc}, nil
}

func (vj *VJoy) SetButton(btn int, status int) {
	desc := vj.handler.(*vjoyWindowsHandler)
	if b := desc.dev.Button(uint(btn)); b != nil {
		b.Set(status == 1)
	}
}

func (vj *VJoy) SetAxis(axis int, pos int) {
	desc := vj.handler.(*vjoyWindowsHandler)
	if axis := desc.dev.Axis(vjoy.AxisName(axis)); axis != nil {
		axis.Seti(pos)
	}
}

func (vj *VJoy) Tick() {
	desc := vj.handler.(*vjoyWindowsHandler)
	desc.dev.Update()
}

func (vj *VJoy) Close() {
	desc := vj.handler.(*vjoyWindowsHandler)
	desc.dev.Relinquish()
}
