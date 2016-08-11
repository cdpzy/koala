package h264

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
