package gps

import (
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	gpx "github.com/twpayne/go-gpx"

	"github.com/ggrrrr/btmt-ui/be/help"
)

const earthR = 6371e3 // metres

func calcDist(wp1 gpx.WptType, wp2 gpx.WptType) float64 {
	if wp1.Lat == 0 || wp1.Lon == 0 ||
		wp2.Lat == 0 || wp2.Lon == 0 {
		return 0
	}
	f1 := wp1.Lat * math.Pi / 180
	f2 := wp2.Lat * math.Pi / 180

	deltaF := (wp2.Lat - wp1.Lat) * math.Pi / 180
	deltaG := (wp2.Lon - wp1.Lon) * math.Pi / 180

	// Math.sin(Δφ/2) * Math.sin(Δφ/2)
	a := (math.Sin(deltaF/2) * math.Sin(deltaF/2)) +
		(math.Cos(f1)*math.Cos(f2))*(math.Sin(deltaG/2)*math.Sin(deltaG/2))
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthR * c
}

func printTrk(trk *gpx.TrkType) {
	fmt.Printf("\n %v %+v \n", trk.Name, trk.Desc)
	fmt.Printf("\tTrkSeg:%v \n", len(trk.TrkSeg))

	for v := range trk.TrkSeg {
		printSeg(trk.TrkSeg[v])
	}
}

func printSeg(seg *gpx.TrkSegType) {
	fmt.Printf("\t\tTrkPt:%v \n", len(seg.TrkPt))
	fmt.Println("[")
	var lastPoint *gpx.WptType
	var delta float64 = 0
	for v := range seg.TrkPt {
		currentPt := seg.TrkPt[v]
		if lastPoint != nil {
			delta = calcDist(*lastPoint, *currentPt)
		}
		lastPoint = currentPt
		printWp(seg.TrkPt[v], delta)
	}
	fmt.Println("]")
}

func printWp(wp *gpx.WptType, dist float64) {
	if dist > 10 {
		fmt.Printf("\t[%v, %v],\n", wp.Lon, wp.Lat)
	}
}

func Test_Load(t *testing.T) {

	pwd := help.RepoDir()

	file, err := os.Open(fmt.Sprintf("%s/gara-dolene-ravno-bore.gpx", pwd))
	require.NoError(t, err)

	gpxReader, err := gpx.Read(file)
	require.NoError(t, err)

	fmt.Printf("t.Wpt[0] == %v", gpxReader.Trk[0].TrkSeg[0].TrkPt[0])
	printTrk(gpxReader.Trk[0])
}
