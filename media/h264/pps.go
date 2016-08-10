package h264

import (
	"errors"

	"github.com/doublemo/koala/helper"
)

// PicParameterSetRBSP Picture parameter set RBSP syntax
type PicParameterSetRBSP struct {
	PicParameterSetID                     uint
	SeqParameterSetID                     uint
	EntropyCodingModeFlag                 uint
	BottomFieldPicOrderInFramePresentFlag uint
	NumSliceGroupsMinus1                  uint
	SliceGroupMapType                     uint
	RunLengthMinus1                       []uint
	TopLeft                               []uint
	BottomRight                           []uint
	SliceGroupChangeDirectionFlag         uint
	SliceGroupChangeRateMinus1            uint
	PicSizeInMapUnitsMinus1               uint
	SliceGroupID                          []uint
	NumRefIdxL0DefaultActiveMinus1        uint
	NumRefIdxL1DefaultActiveMinus1        uint
	WeightedPredFlag                      uint
	WeightedBipredIdc                     uint
	PicInitQpMinus26                      int
	PicInitQsMinus26                      int
	ChromaQpIndexOffset                   int
	DeblockingFilterControlPresentFlag    uint
	ConstrainedIntraPredFlag              uint
	RedundantPicCntPresentFlag            uint
	Transform8x8ModeFlag                  uint
	PicScalingMatrixPresentFlag           uint
	PicScalingListPresentFlag             []uint
	ScalingList4X4                        []*SequenceParameterSetScaling
	ScalingList8X8                        []*SequenceParameterSetScaling
	SecondChromaQpIndexOffset             int
	ChromaFormatIdc                       uint
}

// PicParameterSetRBSP Picture parameter set RBSP syntax
func NewPicParameterSetRBSP() *PicParameterSetRBSP {
	return new(PicParameterSetRBSP)
}

func (picParameterSetRBSP *PicParameterSetRBSP) ParseBytes(b []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else if e, ok := r.(string); ok {
				err = errors.New(e)
			} else {
				panic(r)
			}
		}
	}()

	var (
		iGroup                        uint
		i                             uint
		pic_scaling_list_present_flag uint
	)

	eg := helper.NewExpGolombReader(b)
	picParameterSetRBSP.PicParameterSetID = handleParseError(eg.ReadUE())
	picParameterSetRBSP.SeqParameterSetID = handleParseError(eg.ReadUE())
	picParameterSetRBSP.EntropyCodingModeFlag = handleParseError(eg.ReadBit())
	picParameterSetRBSP.BottomFieldPicOrderInFramePresentFlag = handleParseError(eg.ReadBit())
	picParameterSetRBSP.NumSliceGroupsMinus1 = handleParseError(eg.ReadUE())

	if picParameterSetRBSP.NumSliceGroupsMinus1 > 0 {
		picParameterSetRBSP.SliceGroupMapType = handleParseError(eg.ReadUE())
		if picParameterSetRBSP.SliceGroupMapType == 0 {
			for iGroup = 0; iGroup < picParameterSetRBSP.NumSliceGroupsMinus1; iGroup++ {
				picParameterSetRBSP.RunLengthMinus1 = append(picParameterSetRBSP.RunLengthMinus1, handleParseError(eg.ReadUE()))
			}
		} else if picParameterSetRBSP.SliceGroupMapType == 2 {
			for iGroup = 0; iGroup < picParameterSetRBSP.NumSliceGroupsMinus1; iGroup++ {
				picParameterSetRBSP.TopLeft = append(picParameterSetRBSP.TopLeft, handleParseError(eg.ReadUE()))
				picParameterSetRBSP.BottomRight = append(picParameterSetRBSP.BottomRight, handleParseError(eg.ReadUE()))
			}
		} else if helper.InUint(picParameterSetRBSP.SliceGroupMapType, []uint{3, 4, 5}) {
			picParameterSetRBSP.SliceGroupChangeDirectionFlag = handleParseError(eg.ReadBit())
			picParameterSetRBSP.SliceGroupChangeRateMinus1 = handleParseError(eg.ReadUE())
		} else if picParameterSetRBSP.SliceGroupMapType == 6 {
			picParameterSetRBSP.PicSizeInMapUnitsMinus1 = handleParseError(eg.ReadUE())
			for i = 0; i < picParameterSetRBSP.PicSizeInMapUnitsMinus1; i++ {
				picParameterSetRBSP.SliceGroupID = append(picParameterSetRBSP.SliceGroupID, handleParseError(eg.ReadUV())) // u(v)
			}
		}
	}

	picParameterSetRBSP.NumRefIdxL0DefaultActiveMinus1 = handleParseError(eg.ReadUE())
	picParameterSetRBSP.NumRefIdxL1DefaultActiveMinus1 = handleParseError(eg.ReadUE())
	picParameterSetRBSP.WeightedPredFlag = handleParseError(eg.ReadBit())
	picParameterSetRBSP.WeightedBipredIdc = handleParseError(eg.ReadBits(2))
	picParameterSetRBSP.PicInitQpMinus26 = handleParseSEError(eg.ReadSE())
	picParameterSetRBSP.PicInitQsMinus26 = handleParseSEError(eg.ReadSE())
	picParameterSetRBSP.ChromaQpIndexOffset = handleParseSEError(eg.ReadSE())
	picParameterSetRBSP.DeblockingFilterControlPresentFlag = handleParseError(eg.ReadBit())
	picParameterSetRBSP.ConstrainedIntraPredFlag = handleParseError(eg.ReadBit())
	picParameterSetRBSP.RedundantPicCntPresentFlag = handleParseError(eg.ReadBit())

	if picParameterSetRBSP.more_rbsp_data(eg) {
		picParameterSetRBSP.Transform8x8ModeFlag = handleParseError(eg.ReadBit())
		picParameterSetRBSP.PicScalingMatrixPresentFlag = handleParseError(eg.ReadBit())

		if picParameterSetRBSP.PicScalingMatrixPresentFlag > 0 {
			n := 2
			if picParameterSetRBSP.ChromaFormatIdc == 3 {
				n = 6
			}

			n = 6 + n*int(picParameterSetRBSP.Transform8x8ModeFlag)
			for j := 0; j < n; j++ {
				pic_scaling_list_present_flag = handleParseError(eg.ReadBit())
				picParameterSetRBSP.PicScalingListPresentFlag = append(picParameterSetRBSP.PicScalingListPresentFlag, pic_scaling_list_present_flag)

				if pic_scaling_list_present_flag > 0 {
					if i < 6 {
						scalingList, useDefaultScalingMatrixFlag := scalingList(eg, 16)
						picParameterSetRBSP.ScalingList4X4 = append(picParameterSetRBSP.ScalingList4X4, &SequenceParameterSetScaling{scalingList, useDefaultScalingMatrixFlag})
					} else {
						scalingList, useDefaultScalingMatrixFlag := scalingList(eg, 16)
						picParameterSetRBSP.ScalingList8X8 = append(picParameterSetRBSP.ScalingList8X8, &SequenceParameterSetScaling{scalingList, useDefaultScalingMatrixFlag})
					}
				}
			}

			picParameterSetRBSP.SecondChromaQpIndexOffset = handleParseSEError(eg.ReadSE())
		}
	}
	return
}

// more_rbsp_data is eof
func (picParameterSetRBSP *PicParameterSetRBSP) more_rbsp_data(eg *helper.ExpGolombReader) bool {
	_, err := eg.ReadAtBits(1, eg.GetStartBit())
	if err != nil {
		return false
	}
	return true
}
