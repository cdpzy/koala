package rtp

import (
	"testing"
)

func TestParseRangeHeader(t *testing.T) {
	testData := []string{
		"npt = 1.023 - 1.2355",
		"npt = 1.023 -",
		"npt = 1 - 2",
		"npt = 1.023 -;dsadsa=dsadsaa",
		"npt =-",
		"NPT =-",
		"npt =now",
		"npt =now-",
		"npt =now-1.231;444=ss",
		"clock=19960213T143205Z-",
		"clock=19960213T143205Z-19960213T143205Z",
		"clock=19960213T143205Z-;time=19970123T143720Z",
		"clock=19960213T143205Z-19960213T143205Z;time=19970123T143720Z",
	}

	for _, s := range testData {
		n, e := ParseRangeHeader(s)
		t.Log(s, "=", n, e)
	}

}
