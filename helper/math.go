package helper

import (
	"math"
)

const (
	PI          float64 = 3.1416
	EarthRadius float64 = 6378.137
)

// MaxInt 选择最大的数
func MaxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// GPSDistance GPS定位距离计算
func GPSDistance(lat1, lon1, lat2, lon2 float64) float64 {
	radlat1 := lat1 * PI / 180.0
	radlon1 := lon1 * PI / 180.0
	radlat2 := lat2 * PI / 180.0
	radlon2 := lon2 * PI / 180.0

	if radlat1 < 0 {
		radlat1 = PI/2 + math.Abs(radlat1) // south
	}

	if radlat1 > 0 {
		radlat1 = PI/2 - math.Abs(radlat1) // north
	}

	if radlon1 < 0 {
		radlon1 = PI*2 - math.Abs(radlon1) // west
	}

	if radlat2 < 0 {
		radlat2 = PI/2 + math.Abs(radlat2) // south
	}

	if radlat2 > 0 {
		radlat2 = PI/2 - math.Abs(radlat2) // north
	}

	if radlon2 < 0 {
		radlon2 = PI*2 - math.Abs(radlon2) // west
	}

	x1 := EarthRadius * math.Cos(radlon1) * math.Sin(radlat1)
	y1 := EarthRadius * math.Sin(radlon1) * math.Sin(radlat1)
	z1 := EarthRadius * math.Cos(radlat1)

	x2 := EarthRadius * math.Cos(radlon2) * math.Sin(radlat2)
	y2 := EarthRadius * math.Sin(radlon2) * math.Sin(radlat2)
	z2 := EarthRadius * math.Cos(radlat2)

	d := math.Sqrt((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2) + (z1-z2)*(z1-z2))
	theta := math.Acos((EarthRadius*EarthRadius + EarthRadius*EarthRadius - d*d) / (2 * EarthRadius * EarthRadius))
	return theta * EarthRadius
}
