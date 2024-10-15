package ubx

import (
	"encoding/binary"
	"fmt"
)

type PacketHead struct {
	Preamble uint16
	MsgClass uint8
	MsgID    uint8
	Length   uint16
}

func (h *PacketHead) Unpack(data []byte) error {
	if len(data) < 6 {
		return fmt.Errorf("not enough data to unpack UBXPacketHead")
	}
	h.Preamble = binary.BigEndian.Uint16(data[0:2])
	if h.Preamble != 0xb562 {
		return fmt.Errorf("invalid preamble: 0x%04x", h.Preamble)
	}
	h.MsgClass = data[2]
	h.MsgID = data[3]
	h.Length = binary.LittleEndian.Uint16(data[4:6])
	return nil
}

func (h *PacketHead) Pack() []byte {
	data := make([]byte, 6)
	binary.BigEndian.PutUint16(data[0:2], h.Preamble)
	data[2] = h.MsgClass
	data[3] = h.MsgID
	binary.LittleEndian.PutUint16(data[4:6], h.Length)
	return data
}

func (h *PacketHead) GetPacketLength() int {
	return 6 + int(h.Length) + 2 // head, payload, checksum
}
