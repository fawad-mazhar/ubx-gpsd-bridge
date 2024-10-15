package tcp

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"example.com/ubx-gpsd-bridge/internal/utils"
	"example.com/ubx-gpsd-bridge/pkg/ubx"
)

type Handler struct {
	conn net.Conn
	ubx  *ubx.I2CDevice
}

func NewHandler(conn net.Conn, ubx *ubx.I2CDevice) *Handler {
	return &Handler{conn: conn, ubx: ubx}
}

func (h *Handler) Handle() {
	defer h.conn.Close()

	for {
		if err := h.processGPSDCmd(); err != nil {
			if err == io.EOF {
				log.Println("Client disconnected")
				return
			}
			log.Printf("Error processing GPSD command: %v", err)
			return
		}

		ubxProcessed, err := h.processUBXData()
		if err != nil {
			log.Printf("Error processing UBX data: %v", err)
			return
		}
		fmt.Printf("ubx processed : %d\n", ubxProcessed)

		time.Sleep(time.Second)
	}
}

func (h *Handler) processGPSDCmd() error {
	buffer := make([]byte, 1024)
	n, err := h.conn.Read(buffer)
	if err != nil {
		if err == io.EOF {
			return err
		}
		return fmt.Errorf("error reading from socket: %v", err)
	}

	fromGPSD := buffer[:n]
	for len(fromGPSD) >= 8 {
		preamble := binary.BigEndian.Uint16(fromGPSD[:2])
		if preamble == 0xb562 {
			var packetHead ubx.PacketHead
			if err := packetHead.Unpack(fromGPSD[:6]); err != nil {
				return err
			}
			packetLength := packetHead.GetPacketLength()
			if len(fromGPSD) < packetLength {
				break
			}
			packet := fromGPSD[:packetLength]
			fromGPSD = fromGPSD[packetLength:]
			fmt.Printf("gpsd =UBX=> M10S : %x\n", packet)
			if err := h.ubx.Write(packet); err != nil {
				return err
			}
			time.Sleep(10 * time.Millisecond)
		} else {
			fmt.Printf("dump gpsd cmd : %x\n", fromGPSD)
			fromGPSD = nil
		}
	}
	return nil
}

func (h *Handler) processUBXData() (int, error) {
	processed := 0
	for {
		receivedCount, err := h.ubx.GetReceivedCount()
		if err != nil {
			return processed, err
		}
		if receivedCount == 0 {
			break
		}

		fromDev, err := h.ubx.Read(receivedCount)
		if err != nil {
			return processed, err
		}
		fmt.Printf("data from i2c : %d\n", len(fromDev))

		for len(fromDev) > 2 {
			preamble := binary.BigEndian.Uint16(fromDev[:2])
			if preamble == 0xFFFF {
				fmt.Printf("M10S => : throwing out FFFF... %d bytes\n", receivedCount)
				break
			}

			if preamble == 0xb562 && len(fromDev) > 6 {
				var packetHead ubx.PacketHead
				if err := packetHead.Unpack(fromDev[:6]); err != nil {
					return processed, err
				}
				packetLength := packetHead.GetPacketLength()
				if len(fromDev) >= packetLength {
					packet := fromDev[:packetLength]
					fromDev = fromDev[packetLength:]
					fmt.Printf("M10S =UBX=> gpsd : %x\n", packet)
					if _, err := h.conn.Write(packet); err != nil {
						return processed, err
					}
					processed++
					time.Sleep(10 * time.Millisecond)
					continue
				} else {
					fmt.Printf("packet is not complete : %d > %d\n", packetLength, len(fromDev))
				}
			} else if preamble == 0x2447 || preamble == 0xa447 {
				// Process NMEA sentences
				lineEnd := -1
				for i := 2; i < len(fromDev)-1; i++ {
					if fromDev[i] == 0x0d && fromDev[i+1] == 0x0a {
						lineEnd = i + 2
						break
					}
				}
				if lineEnd != -1 {
					line := append([]byte{0x24, 0x47}, fromDev[2:lineEnd]...)
					if !utils.Contains(line, 0xff) {
						fmt.Printf("M10S =NMEA=> gpsd : %s", string(line))
						if _, err := h.conn.Write(line); err != nil {
							return processed, err
						}
						fromDev = fromDev[lineEnd:]
						processed++
						time.Sleep(10 * time.Millisecond)
						continue
					}
				}
			}

			fmt.Printf("dump i2c data : %x\n", fromDev)
			break
		}
	}
	return processed, nil
}
