package ubx

import (
	"encoding/binary"
	"fmt"

	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

type I2CDevice struct {
	bus  i2c.BusCloser
	addr uint16
}

func NewI2CDevice(busNumber int, addr uint16) (*I2CDevice, error) {
	// Initialize the host, which is required for I2C communication
	if _, err := host.Init(); err != nil {
		return nil, err
	}

	// Open the I2C bus
	bus, err := i2creg.Open(fmt.Sprintf("/dev/i2c-%d", busNumber))
	if err != nil {
		return nil, err
	}

	return &I2CDevice{
		bus:  bus,
		addr: addr,
	}, nil
}

func (d *I2CDevice) Write(data []byte) error {
	return d.i2cWrite(0xFF, data)
}

func (d *I2CDevice) Read(count int) ([]byte, error) {
	return d.i2cRead(0xFF, count)
}

func (d *I2CDevice) GetReceivedCount() (int, error) {
	countReg, err := d.i2cRead(0xFD, 2)
	if err != nil {
		return 0, err
	}
	count := int(binary.BigEndian.Uint16(countReg))
	if count > 0xFFF {
		count = 0
	}
	return count, nil
}

// i2cRead reads data from a specific register on the I2C device
func (d *I2CDevice) i2cRead(reg byte, count int) ([]byte, error) {
	data := make([]byte, 0, count)
	for count > 0 {
		readSize := min(count, 32)
		chunk := make([]byte, readSize)
		if err := d.bus.Tx(d.addr, []byte{reg}, chunk); err != nil {
			return nil, err
		}
		data = append(data, chunk...)
		count -= readSize
	}
	return data, nil
}

// i2cWrite writes data to a specific register on the I2C device
func (d *I2CDevice) i2cWrite(reg byte, data []byte) error {
	for len(data) > 0 {
		writeSize := min(len(data), 31)
		chunk := append([]byte{reg}, data[:writeSize]...)
		if err := d.bus.Tx(d.addr, chunk, nil); err != nil {
			return err
		}
		data = data[writeSize:]
	}
	return nil
}
