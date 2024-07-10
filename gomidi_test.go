package gomidi

import (
	"testing"
	"time"

	log "github.com/schollz/logger"
)

func TestRun(t *testing.T) {
	log.SetLevel("trace")
	names := Devices()
	log.Debugf("found %d output devices: ", len(names))
	for i, name := range names {
		log.Debugf("device %d: %s", i, name)
	}
	if len(names) == 0 {
		t.Skip("no devices found")
	}

	device, err := New("USB MIDI")
	if err != nil {
		t.Skip("no USB MIDI device found")
	}
	err = device.Open()
	if err != nil {
		t.Error(err)
	}

	err = device.NoteOn(0, 72, 100)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(1 * time.Second)

	err = device.NoteOff(0, 72)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(1 * time.Second)

	err = device.Close()
	if err != nil {
		t.Error(err)
	}

	Close()
}
