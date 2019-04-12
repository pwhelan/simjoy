package midi

import (
	"github.com/rakyll/portmidi"
)

type MIDI struct {
	in  *portmidi.Stream
	out *portmidi.Stream
}

func (m *MIDI) Send(channel int64, status int64, hb int64, lb int64) {
	m.out.WriteShort((channel&0x0f)|(status&0xf0), hb, lb)
}

type Event struct {
	portmidi.Event
	DeviceID int
}

type MIDIS struct {
	Devices []*MIDI
	Channel chan Event
}

func init() {
	portmidi.Initialize()
}

func OpenDevices() (*MIDIS, error) {
	devs := make([]*MIDI, 0)
	for i := 0; i < portmidi.CountDevices(); i += 2 {
		midi, err := OpenDevice(i)
		if err != nil {
			return nil, err
		}
		if midi != nil {
			devs = append(devs, midi)
		}
	}
	midis := &MIDIS{Devices: devs, Channel: make(chan Event)}
	for id, dev := range devs {
		go func(deviceID int, c <-chan portmidi.Event) {
			for {
				ev := <-c
				midis.Channel <- Event{ev, id}
			}
		}(id, dev.in.Listen())
	}
	return midis, nil
}

func OpenDevice(deviceID int) (*MIDI, error) {
	out, err := portmidi.NewOutputStream(portmidi.DeviceID(deviceID), 1024, 0)
	if err != nil {
		return nil, err
	}
	in, err := portmidi.NewInputStream(portmidi.DeviceID(deviceID+1), 1024)
	if err != nil {
		out.Close()
		return nil, err
	}
	return &MIDI{in, out}, nil
}
