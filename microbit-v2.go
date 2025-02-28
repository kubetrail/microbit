//go:build microbit_v2

package microbit

import (
	"context"
	"machine"
	"sync"
	"time"
)

const (
	NumRows = 5
	NumCols = 5
)

type Device struct {
	mu sync.Mutex

	colPins [NumCols]machine.Pin
	rowPins [NumRows]machine.Pin
	buffer  [NumRows][NumCols]bool

	buzzerPin machine.Pin

	buttonA machine.Pin
	buttonB machine.Pin
}

// NewDevice returns a new device
func NewDevice() *Device {
	device := new(Device)

	device.colPins[0] = machine.LED_COL_1
	device.colPins[1] = machine.LED_COL_2
	device.colPins[2] = machine.LED_COL_3
	device.colPins[3] = machine.LED_COL_4
	device.colPins[4] = machine.LED_COL_5

	device.rowPins[0] = machine.LED_ROW_1
	device.rowPins[1] = machine.LED_ROW_2
	device.rowPins[2] = machine.LED_ROW_3
	device.rowPins[3] = machine.LED_ROW_4
	device.rowPins[4] = machine.LED_ROW_5

	for _, pin := range device.colPins {
		pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	}

	for _, pin := range device.rowPins {
		pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	}

	device.buzzerPin = machine.BUZZER
	device.buzzerPin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	device.buttonA = machine.BUTTONA
	device.buttonB = machine.BUTTONB

	device.buttonA.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	device.buttonB.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	return device
}

// Wait is a blocking call that will block until input context is done
func (device *Device) Wait(ctx context.Context) *Device {
	select {
	case <-ctx.Done():
	}

	return device
}

// OnButtonAPress returns a context that is done when button A is pressed
func (device *Device) OnButtonAPress() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()
		for {
			if !device.buttonA.Get() {
				return
			}
			time.Sleep(time.Millisecond * 250)
		}
	}()

	return ctx
}

// OnButtonBPress returns a context that is done when the button B is pressed
func (device *Device) OnButtonBPress() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()
		for {
			if !device.buttonB.Get() {
				return
			}
			time.Sleep(time.Millisecond * 250)
		}
	}()

	return ctx
}

// OnButtonPress returns a context that is done when the button B is pressed
func (device *Device) OnButtonPress() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()
		for {
			if !device.buttonA.Get() || !device.buttonB.Get() {
				return
			}
			time.Sleep(time.Millisecond * 250)
		}
	}()

	return ctx
}

// clear clears the internal state
func (device *Device) clear() *Device {
	device.mu.Lock()
	defer device.mu.Unlock()
	for i, row := range device.buffer {
		for j := range row {
			device.buffer[i][j] = false
		}
	}

	return device
}

// SetMatrix configures the internal state to display on LED matrix later
func (device *Device) SetMatrix(image [NumRows][NumCols]uint8) *Device {
	device.mu.Lock()
	defer device.mu.Unlock()

	for i := 0; i < NumRows; i++ {
		for j := 0; j < NumCols; j++ {
			if image[i][j] != 0 {
				device.buffer[i][j] = true
			} else {
				device.buffer[i][j] = false
			}
		}
	}
	return device
}

// Buzz makes a sound at a given frequency and stops when the input context is done
func (device *Device) Buzz(ctx context.Context, frequency uint32) *Device {
	go func() {
		pin := device.buzzerPin
		periodNs := uint64(1000000000) / uint64(frequency) // Calculate the period in nanoseconds

		for {
			select {
			case <-ctx.Done():
				return
			default:
				pin.High()
				time.Sleep(time.Duration(periodNs / 2)) // High for half the period
				pin.Low()
				time.Sleep(time.Duration(periodNs / 2)) // Low for half the period
			}
		}

		pin.Low()
	}()

	return device
}

// Display enables LED's on the 5x5 matrix according to the image set
func (device *Device) Display(ctx context.Context) *Device {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				for col := 0; col < NumCols; col++ {
					device.colPins[col].High()
					for row := 0; row < NumRows; row++ {
						if device.buffer[row][col] {
							device.rowPins[row].High()
							device.colPins[col].Low()
						}
						time.Sleep(time.Microsecond * 500)
						device.rowPins[row].Low()
						device.colPins[col].High()
					}
				}
			}
		}
	}()

	return device
}
