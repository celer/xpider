package hdlc

import (
	//"github.com/howeyc/crc16"
	"io"
)

type Reader struct {
	reader   io.Reader
	frame    []byte
	framePos int
	frameLen int
	escape   bool
}

func NewReader(r io.Reader) *Reader {
	reader := &Reader{reader: r}
	reader.frame = make([]byte, 2048)
	return reader
}

func (r *Reader) crcCheck(pos uint16) bool {
	//FIXME Implement this!
	/*
		crcRecv:=binary.LittleEndian.Uint16(r.frameBuffer[pos-2:])
		crcCalc:=crc16.ChecksumCCITT(r.frameBuffer[:pos-2])
		crcCalc^=0xFFFF
		fmt.Printf("%x %x\n",crcRecv,crcCalc)
		return crcRecv==crcCalc
	*/
	return false
}

func (r *Reader) Read(p []byte) (int, error) {
	pl := 0
	var err error
	for true {
		if pl+1 > cap(p) {
			return pl, nil
		}
		if r.frameLen == 0 || r.frameLen == r.framePos {
			r.frameLen, err = r.reader.Read(r.frame)
			if err != nil {
				return 0, err
			}
		}
		d := r.frame[r.framePos]
		r.framePos += 1
		if d == FRAME_BOUNDARY_OCTET {
			if pl == 0 {
				//FRAME START
				continue
			} else {
				//FRAME END
				return pl - 2, nil
			}
		}
		if d == CONTROL_ESCAPE_OCTET {
			r.escape = true
			continue
		}
		if r.escape {
			d ^= INVERT_OCTET
		}
		p[pl] = d
		pl += 1
	}
	return 0, nil
}
