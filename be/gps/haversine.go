package gps

import (
	"math"

	"github.com/twpayne/go-gpx"
)

const earthRKm float64 = 6371 // 6371
// const earthR = 6378.137e.

func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}

// https://en.wikipedia.org/wiki/Earth_radius
// https://github.com/umahmood/haversine/blob/master/haversine.go
func Distance(p1, p2 gpx.WptType) float64 {
	lat1 := degreesToRadians(p1.Lat)
	lon1 := degreesToRadians(p1.Lon)
	lat2 := degreesToRadians(p2.Lat)
	lon2 := degreesToRadians(p2.Lon)

	diffLat := lat2 - lat1
	diffLon := lon2 - lon1

	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*
		math.Pow(math.Sin(diffLon/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return c * earthRKm * 1000

}
