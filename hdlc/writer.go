package hdlc

import (
	"github.com/howeyc/crc16"
	"io"
)

func low(x uint16) byte {
	return byte(((x) & 0xFF))
}

func high(x uint16) byte {
	return byte((((x) >> 8) & 0xFF))
}

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	h := &Writer{writer: w}
	return h
}

func (h *Writer) Write(data []byte) (int, error) {
	var fcs uint16
	var d byte
	fcs = crc16.ChecksumCCITT(data)
	packet := make([]byte, 0)
	packet = append(packet, FRAME_BOUNDARY_OCTET)

	for i := 0; i < len(data); i++ {
		d = data[i]
		if (d == CONTROL_ESCAPE_OCTET) ||
			(d == FRAME_BOUNDARY_OCTET) {
			packet = append(packet, CONTROL_ESCAPE_OCTET)
			d ^= INVERT_OCTET
		}
		packet = append(packet, d)
	}
	fcs ^= 0xFFFF
	d = low(fcs)
	if (d == CONTROL_ESCAPE_OCTET) ||
		(d == FRAME_BOUNDARY_OCTET) {
		packet = append(packet, CONTROL_ESCAPE_OCTET)
		d ^= INVERT_OCTET
	}
	packet = append(packet, d)
	d = high(fcs)
	if (d == CONTROL_ESCAPE_OCTET) ||
		(d == FRAME_BOUNDARY_OCTET) {
		packet = append(packet, CONTROL_ESCAPE_OCTET)
		d ^= INVERT_OCTET
	}
	packet = append(packet, d)
	packet = append(packet, FRAME_BOUNDARY_OCTET)
	_, err := h.writer.Write(packet)
	if err != nil {
		return 0, err
	}
	return len(data), nil
}
