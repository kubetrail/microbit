# microbit
microbit i/o using tinygo

[Micro:Bit v2](https://microbit.org/new-microbit/) is a fun microcontroller for
interesting educational projects. This repo aims to build simple I/O using
[tinygo](https://tinygo.org/) as a backend.

## installation
You will need `tinygo` to build and flash this code. Connect `micro:bit` and identify
the port:
```bash
tinygo ports
```
```text
Port                 ID        Boards
/dev/ttyACM0         ****:****
```

Flash the code specifying target and the port. Of course, change the
port to what you see on your computer
```bash
tinygo flash -target microbit-v2 -port /dev/ttyACMO
```

## usage
Example code below displays three different patterns on LED matrix along
with buzzer buzzing at different tones that change on button press.
```go
package main

import (
	"context"
	"time"

	"github.com/kubetrail/microbit"
)

func main() {
	// create a device handle
	device := microbit.NewDevice()

	// get context that will cancel on either button A or button B press
	ctx := device.OnButtonPress()

	// set LED matrix, non-zero values are true and will light up LEDs,
	// then display this on LED matrix, start a buzzer simultaneously,
	// and wait for button press.
	device.SetMatrix([microbit.NumRows][microbit.NumCols]uint8{
		{0, 1, 0, 1, 0},
		{1, 0, 1, 0, 1},
		{0, 1, 0, 1, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 0, 0, 0},
	},
	).Display(ctx).Buzz(ctx, 400).Wait(ctx)

	time.Sleep(time.Millisecond * 250)
	ctx = device.OnButtonPress()
	device.SetMatrix([microbit.NumRows][microbit.NumCols]uint8{
		{1, 0, 0, 0, 0},
		{0, 1, 0, 0, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 0, 1, 0},
		{0, 0, 0, 0, 1},
	},
	).Display(ctx).Buzz(ctx, 200).Wait(ctx)

	time.Sleep(time.Millisecond * 250)
	ctx = device.OnButtonPress()
	device.SetMatrix([microbit.NumRows][microbit.NumCols]uint8{
		{1, 1, 0, 1, 1},
		{1, 1, 0, 1, 1},
		{1, 1, 1, 1, 1},
		{1, 1, 0, 1, 1},
		{1, 1, 0, 1, 1},
	},
	).Display(ctx).Wait(ctx)

	time.Sleep(time.Millisecond * 250)
	ctx = context.Background()
	device.SetMatrix([microbit.NumRows][microbit.NumCols]uint8{
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
	},
	).Display(ctx)

	time.Sleep(time.Second)
}
```
