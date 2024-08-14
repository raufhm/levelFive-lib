package printing

import (
	"errors"
	"fmt"
	"github.com/google/gousb"
	"github.com/raufhm/levelfive-lib/shared"
	"net"
)

type Printer interface {
	Print(formattedStr string, args shared.PrinterArgs) error
}

type PrintType int

const (
	Network PrintType = iota
	USB
)

func NewPrinter(pt PrintType) (Printer, error) {
	switch pt {
	case Network:
		return &OverNetwork{}, nil
	case USB:
		return &OverUSB{}, nil
	default:
		return nil, errors.New("unsupported print type")
	}
}

type OverUSB struct{}

func (p *OverUSB) Print(formattedStr string, _ shared.PrinterArgs) error {
	// Reference is https://github.com/google/gousb/blob/master/example_test.go

	// Initialize a new Context.
	ctx := gousb.NewContext()
	defer func(ctx *gousb.Context) {
		err := ctx.Close()
		if err != nil {
			return
		}
	}(ctx)

	// Open any device with a given VID/PID using a convenience function.
	// VID/PID can be obtained from the terminal. For example, I am using MACOS, I can check it with this command `system_profiler SPUSBDataType`
	dev, err := ctx.OpenDeviceWithVIDPID(0x0781, 0x5567)
	if err != nil {
		return err
	}
	defer func(dev *gousb.Device) {
		err := dev.Close()
		if err != nil {
			return
		}
	}(dev)

	cfg, err := dev.Config(1)
	if err != nil {
		return err
	}

	// check interface for debugging
	i, err := cfg.Interface(0, 0)
	if err != nil {
		return err
	}
	fmt.Println(i)

	// TODO: this claim interface is not working, it keep returning bad access [code -3]
	intf, done, err := dev.DefaultInterface()
	if err != nil {
		return err
	}
	defer done()

	// Open an OUT endpoint.
	ep, err := intf.OutEndpoint(7)
	if err != nil {
		return err
	}

	// Generate some data to write.
	data := []byte(formattedStr)

	// Write data to the USB device.
	_, err = ep.Write(data)
	if err != nil {
		return err
	}
	return nil
}

type OverNetwork struct{}

func (p *OverNetwork) Print(formattedStr string, args shared.PrinterArgs) error {
	var address, port string
	if args.NetworkArgs == nil {
		address = shared.DefaultAddress
		port = shared.DefaultPort
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", address, port))
	if err != nil {
		return fmt.Errorf("could not connect to printer: %v", err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	// Send the data to the printer
	_, err = conn.Write([]byte(formattedStr))
	if err != nil {
		return fmt.Errorf("failed to write to printer: %v", err)
	}
	return nil
}
