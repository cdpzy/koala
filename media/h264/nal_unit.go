package h264

import (
	"bytes"
	"errors"
	"fmt"

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
}

// NewNalUnit NAL unit syntax
func NewNalUnit() *NalUnit {
	return new(NalUnit)
}

// ParseBytes parse
func (nalUnit *NalUnit) ParseBytes(reader *bytes.Reader) error {
	header := make([]byte, 4)
	_, err := reader.Read(header)
	if err != nil {
		return err
	}

	test1 := header[0] == 0x00 && header[1] == 0x00 && header[2] == 0x00 && header[3] == 0x01
	test2 := header[0] == 0x00 && header[1] == 0x00 && header[2] == 0x01

	if !test1 && !test2 {
		return errors.New("NotFound")
	}

	if test2 {
		reader.UnreadByte()
	}

	_, err = reader.Read(header)
	if err != nil {
		return err
	}

	eg := helper.NewExpGolombReader(header)
	nalUnit.ForbiddenZeroBit = handleParseError(eg.ReadBit())
	nalUnit.NalRefIdc = handleParseError(eg.ReadBits(2))
	nalUnit.NalUnitType = handleParseError(eg.ReadBits(5))
	return nil
}

// read
func (nalUnit *NalUnit) Read(reader *bytes.Reader, data *[]byte) error {
	byte0, err := reader.ReadByte()
	if err != nil {
		return err
	}

	if byte0 != 0x00 {
		*data = append(*data, byte0)
		return nalUnit.Read(reader, data)
	}

	reader.UnreadByte()
	header := make([]byte, 4)
	_, err = reader.Read(header)
	if err != nil {
		return err
	}
	fmt.Println("header = ", header)
	reader.UnreadByte()
	reader.UnreadByte()
	reader.UnreadByte()

	test1 := header[0] == 0x00 && header[1] == 0x00 && header[2] == 0x00 && header[3] == 0x01
	test2 := header[0] == 0x00 && header[1] == 0x00 && header[2] == 0x01
	if !test1 && !test2 {
		*data = append(*data, header[0])
		return nalUnit.Read(reader, data)
	}

	reader.UnreadByte()
	return nil
}
