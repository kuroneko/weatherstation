package probe

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Probe struct {
	devicePath string
	dev        *os.File
	err        error
	lineBuffer string
	Updates    chan *Update
}

type Update struct {
	When        time.Time
	Temperature *float64
	Humidity    *float64
	DoorStatus  bool
}

func Open(path string) (p *Probe, err error) {
	p = &Probe{
		devicePath: path,
	}
	p.Updates = make(chan *Update, 1)

	p.dev, err = os.OpenFile(p.devicePath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Probe) Start() {
	go p.poll()
}

func (p *Probe) Stop() {
	p.dev.Close()
	close(p.Updates)
}

func (p *Probe) poll() {
	for {
		var nextB [1]byte

		c, err := p.dev.Read(nextB[:])
		if err != nil {
			p.err = err
			break
		}
		if c == 1 {
			if nextB[0] == '\n' {
				if len(p.lineBuffer) > 0 {
					p.process()
				}
				p.lineBuffer = ""
			} else {
				p.lineBuffer += string(nextB[:])
			}
		} else {
			break
		}
	}
}

func (p *Probe) process() {
	if p.lineBuffer[0] != '^' {
		// string has wrong preamble - probably out of sync.
		return
	}
	probeData := p.lineBuffer[1:]
	// split the line into it's three components
	dataparts := strings.Split(probeData, ":")

	result := &Update{
		When: time.Now(),
	}

	tempFloat, err := strconv.ParseFloat(dataparts[0], 64)
	if err == nil {
		result.Temperature = &tempFloat
	}
	humidityFloat, err := strconv.ParseFloat(dataparts[1], 64)
	if err == nil {
		result.Humidity = &humidityFloat
	}
	result.DoorStatus, _ = strconv.ParseBool(dataparts[2])

	select {
	case p.Updates <- result:
	default:
	}
}
