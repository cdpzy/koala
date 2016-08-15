package h264

import (
	"errors"

	"github.com/doublemo/koala/helper"
)

const (
	NALU_TYPE_SLICE    uint = iota + 1 // NALU_TYPE_SLICE Coded slice of a non-IDR picture
	NALU_TYPE_DPA                      // Coded slice data partition A
	NALU_TYPE_DPB                      // Coded slice data partition B
	NALU_TYPE_DPC                      // Coded slice data partition C
	NALU_TYPE_IDR                      // Coded slice of an IDR picture
	NALU_TYPE_SEI                      // Supplemental enhancement information (SEI)
	NALU_TYPE_SPS                      // Sequence parameter set
	NALU_TYPE_PPS                      // Picture parameter set
	NALU_TYPE_AUD                      // Access unit delimiter
	NALU_TYPE_EOSEQ                    // End of sequence
	NALU_TYPE_EOSTREAM                 // End of stream
	NALU_TYPE_FILL                     // Filler data
	NALU_TYPE_SPSE                     // Sequence parameter set extension
	NALU_TYPE_PNU                      // Prefix NAL unit
	NALU_TYPE_SSPS                     // Subset sequence parameter set
	NALU_TYPE_DPS                      // Depth parameter set
)

// NalUnit NAL unit syntax
type NalUnit struct {
	ForbiddenZeroBit   uint
	NalRefIdc          uint
	NalUnitType        uint
	SvcExtensionFlag   uint
	Avc3dExtensionFlag uint
	ParameterBytes     []byte
}

// NewNalUnit NAL unit syntax
func NewNalUnit() *NalUnit {
	nal := new(NalUnit)
	nal.ParameterBytes = make([]byte, 0)
	return nal
}

// ParseBytes parse
func (nalUnit *NalUnit) ParseBytes(b []byte) (index int, err error) {
	l, n := nalUnit.IndexByte(b, 0)
	if n == -1 {
		err = errors.New("nomatch")
		return
	}

	idx := n

	var data []byte
	_, next := nalUnit.IndexByte(b, n+l)
	if next == -1 {
		data = b[n:]
		idx = len(b)
	} else {
		data = b[n:next]
		idx = next
	}

	eg := helper.NewExpGolombReader(data[l:])
	nalUnit.ForbiddenZeroBit, err = eg.ReadBit()
	if err != nil {
		return
	}

	nalUnit.NalRefIdc, err = eg.ReadBits(2)
	if err != nil {
		return
	}

	nalUnit.NalUnitType, err = eg.ReadBits(5)
	if err != nil {
		return
	}

	nalUnit.ParameterBytes = data[l+1:]
	index = idx
	return
}

func (nalUnit *NalUnit) IndexByte(b []byte, s int) (int, int) {
	if s > len(b) || len(b) < 3 {
		return 0, -1
	}
	m := false
	x := 0
	for s < len(b) {
		if s+3 <= len(b) {
			if b[s] == 0x00 && b[s+1] == 0x00 && b[s+2] == 0x01 {
				m = true
				x = 3
				break
			}
		}

		if s+4 <= len(b) {
			if b[s] == 0x00 && b[s+1] == 0x00 && b[s+2] == 0x00 && b[s+3] == 0x1 {
				m = true
				x = 4
				break
			}
		}
		s++
	}

	if !m {
		return 0, -1
	}

	return x, s
}
