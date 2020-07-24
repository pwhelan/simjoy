package midi

import (
	"fmt"

	"github.com/rakyll/portmidi"
)

type MIDI struct {
	ID   portmidi.DeviceID
	Info *portmidi.DeviceInfo
	in   *portmidi.Stream
	out  *portmidi.Stream
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

func (midis *MIDIS) listenDevice(midi *MIDI, c <-chan portmidi.Event) {
	fmt.Printf("Listening on %s[%d]\n", midi.Info.Name, midi.ID)
	for {
		fmt.Printf("Recv'd on %s[%d]\n", midi.Info.Name, midi.ID)
		ev := <-c
		midis.Channel <- Event{ev, int(midi.ID)}
	}
}

func OpenDevices() (*MIDIS, error) {
	devs := make([]*MIDI, 0)
	for i := 0; i < portmidi.CountDevices(); i += 2 {
		midi, err := openDevice(portmidi.DeviceID(i))
		if err != nil {
			return nil, err
		}
		if midi != nil {
			devs = append(devs, midi)
		}
	}
	midis := &MIDIS{Devices: devs, Channel: make(chan Event)}
	for _, dev := range devs {
		fmt.Printf("Listen on %d\n", dev.ID)
		go midis.listenDevice(dev, dev.in.Listen())
	}
	return midis, nil
}

func openDevice(deviceID portmidi.DeviceID) (*MIDI, error) {
	out, err := portmidi.NewOutputStream(deviceID, 1024, 0)
	if err != nil {
		return nil, err
	}
	in, err := portmidi.NewInputStream(deviceID+1, 1024)
	if err != nil {
		out.Close()
		return nil, err
	}
	return &MIDI{
		deviceID,
		portmidi.Info(deviceID),
		in, out}, nil
}
