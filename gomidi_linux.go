package gomidi

import (
	"sync"

	log "github.com/schollz/logger"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

var mutex sync.Mutex
var devicesOpen map[string]drivers.Out

func init() {
	devicesOpen = make(map[string]drivers.Out)
}

type Device struct {
	name    string
	num     int
	notesOn map[uint8]bool
}

func New(name string) (d Device, err error) {
	d.name, d.num, err = filterName(name)
	d.notesOn = make(map[uint8]bool)
	return
}

func Close() {
	mutex.Lock()
	defer mutex.Unlock()
	for _, out := range devicesOpen {
		out.Close()
	}
}

func (d Device) Open() (err error) {
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := devicesOpen[d.name]; ok {
		return
	}
	out, err := midi.FindOutPort(d.name)
	if err == nil {
		devicesOpen[d.name] = out
		err = out.Open()
	}
	if err == nil {
		log.Tracef("opened %s", d.name)
	} else {
		log.Error(err)
	}
	return
}

func (d Device) Close() (err error) {
	// send note off to every note
	for note := range d.notesOn {
		d.NoteOff(0, note)
	}
	mutex.Lock()
	defer mutex.Unlock()
	if out, ok := devicesOpen[d.name]; ok {
		err = out.Close()
		delete(devicesOpen, d.name)
	}
	return
}

func (d Device) NoteOn(channel, note, velocity uint8) (err error) {
	mutex.Lock()
	defer mutex.Unlock()
	if out, ok := devicesOpen[d.name]; ok {
		err = out.Send([]byte{0x90 | channel, note, velocity})
		if err != nil {
			log.Error(err)
		} else {
			d.notesOn[note] = true
			log.Tracef("[%s] note on %d %d %d", d.name, channel, note, velocity)
		}
	}
	return
}

func (d Device) NoteOff(channel, note uint8) (err error) {
	mutex.Lock()
	defer mutex.Unlock()
	if out, ok := devicesOpen[d.name]; ok {
		err = out.Send([]byte{0x80 | channel, note, 0})
		if err != nil {
			log.Error(err)
		} else {
			delete(d.notesOn, note)
			log.Tracef("[%s] note off %d %d", d.name, channel, note)
		}
	}
	return
}

func Devices() (devices []string) {
	outs := midi.GetOutPorts()
	for _, out := range outs {
		devices = append(devices, out.String())
	}
	return
}
