package gpsd

import (
	"log"
	"net"

	"example.com/ubx-gpsd-bridge/pkg/tcp"
	"example.com/ubx-gpsd-bridge/pkg/ubx"
)

type Feeder struct {
	listener net.Listener
	ubx      *ubx.I2CDevice
}

func NewFeeder(ubx *ubx.I2CDevice, address string) (*Feeder, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	return &Feeder{
		listener: listener,
		ubx:      ubx,
	}, nil
}

func (f *Feeder) Serve() error {
	for {
		conn, err := f.listener.Accept()
		if err != nil {
			return err
		}
		log.Printf("New connection from %s", conn.RemoteAddr())
		go tcp.NewHandler(conn, f.ubx).Handle()
	}
}
