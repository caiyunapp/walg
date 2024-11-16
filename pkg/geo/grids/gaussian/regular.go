package gaussian

import (
	"math"

	"github.com/scorix/walg/pkg/geo/grids"
)

type regular struct {
	n          int
	latitudes  []float64
	longitudes []float64
}

func NewRegular(n int) *regular {
	r := &regular{n: n}
	r.latitudes = r.calcLatitudes()
	r.longitudes = r.calcLongitudes()

	return r
}

func (g *regular) Size() int {
	return g.longitudesSize() * g.latitudesSize()
}

func (g *regular) latitudesSize() int {
	return g.n * 2
}

func (g *regular) longitudesSize() int {
	return g.n * 4
}

func (g *regular) GridIndex(lat, lon float64) int {
	latIdx := g.GetLatitudeIndex(lat)
	lonIdx := g.GetLongitudeIndex(lon)

	return g.GridIndexFromIndices(latIdx, lonIdx)
}

func (g *regular) GetLatitudeIndex(lat float64) int {
	latIdxArr := grids.FindNearestIndices(lat, g.latitudes)
	a, b := latIdxArr[0], latIdxArr[1]

	if math.Abs(g.latitudes[a]-lat) <= math.Abs(g.latitudes[b]-lat) {
		return a
	}
	return b
}

func (g *regular) GetLongitudeIndex(lon float64) int {
	lonIdxArr := grids.FindNearestIndices(lon, g.longitudes)
	a, b := lonIdxArr[0], lonIdxArr[1]

	if math.Abs(g.longitudes[a]-lon) <= math.Abs(g.longitudes[b]-lon) {
		return a
	}
	return b
}

func (g *regular) GridPoint(index int) (lat, lon float64) {
	latIdx := index / g.longitudesSize()
	lonIdx := index % g.longitudesSize()

	return g.latitudes[latIdx], g.longitudes[lonIdx]
}

func (g *regular) GridIndexFromIndices(latIdx, lonIdx int) int {
	return latIdx*g.longitudesSize() + lonIdx
}

func (g *regular) Latitudes() []float64 {
	return g.latitudes
}

func (g *regular) calcLatitudes() []float64 {
	latitudes := gaussLegendreZeros(g.latitudesSize())

	return latitudes
}

func (g *regular) Longitudes() []float64 {
	return g.longitudes
}

func (g *regular) calcLongitudes() []float64 {
	length := g.longitudesSize()
	longitudes := make([]float64, length)

	for i := 0; i < length; i++ {
		longitudes[i] = 360.0 * float64(i) / float64(length)
	}

	return longitudes
}
