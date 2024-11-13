package gps

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	gpx "github.com/twpayne/go-gpx"

	"github.com/ggrrrr/btmt-ui/be/help"
)

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
	var i = 0
	for v := range seg.TrkPt {
		currentPt := seg.TrkPt[v]
		if lastPoint != nil {
			d := Distance(*lastPoint, *currentPt)
			delta = delta + d
			fmt.Printf("// p1: %v, %v p2: %v, %v == %2f ++ %2f\n", lastPoint.Lon, lastPoint.Lat, currentPt.Lon, currentPt.Lat, d, delta)
		}
		lastPoint = currentPt
		if delta > 100 {
			printWp(seg.TrkPt[v])
			delta = 0
			i++
		}
	}
	fmt.Println("]")
	fmt.Printf("// points: %d\n", i)

}

func printWp(wp *gpx.WptType) {
	// if dist > 10 {
	// fmt.Printf("\t[%v, %v],//[%d] %2f\n", wp.Lon, wp.Lat, i, dist)
	fmt.Printf("\t[%v, %v],\n", wp.Lon, wp.Lat)
	// }
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
