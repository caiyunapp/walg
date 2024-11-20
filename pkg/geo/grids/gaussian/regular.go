package gaussian

import (
	"cmp"
	"math"

	"github.com/scorix/walg/pkg/geo/distance"
	"github.com/scorix/walg/pkg/geo/grids"
)

type regular struct {
	n          int
	latitudes  []float64
	longitudes []float64
	scanMode   grids.ScanMode
}

type RegularOption func(*regular)

func WithScanMode(scanMode grids.ScanMode) RegularOption {
	return func(r *regular) {
		r.scanMode = scanMode
	}
}

func NewRegular(n int, opts ...RegularOption) *regular {
	r := &regular{
		n: n,
	}

	for _, opt := range opts {
		opt(r)
	}

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

func (g *regular) ScanMode() grids.ScanMode {
	return g.scanMode
}

func (g *regular) GetNearestIndex(lat, lon float64) (int, int) {
	latitudes := g.Latitudes()
	longitudes := g.Longitudes()

	indicesLat := grids.FindNearestIndices(lat, latitudes)
	indicesLon := grids.FindNearestIndices(lon, longitudes)

	latIdx := indicesLat[0]
	lonIdx := indicesLon[0]

	// iterations 3 is enough for comparing the minimum distance
	const iterations = 3

	dist := distance.VincentyIterations(lat, lon, latitudes[latIdx], longitudes[lonIdx], iterations)
	for _, i := range indicesLat {
		for _, j := range indicesLon {
			d := distance.VincentyIterations(lat, lon, latitudes[i], longitudes[j], iterations)
			if cmp.Compare(d, dist) < 0 {
				dist = d
				latIdx = i
				lonIdx = j
			}
		}
	}

	return latIdx, lonIdx
}

func (g *regular) GuessNearestIndex(lat, lon float64) (int, int) {
	latitudes := g.Latitudes()
	longitudes := g.Longitudes()

	indicesLat := grids.FindNearestIndices(lat, latitudes)
	indicesLon := grids.FindNearestIndices(lon, longitudes)

	latIdx := indicesLat[0]
	lonIdx := indicesLon[0]

	if math.Abs(latitudes[latIdx]-lat) > math.Abs(latitudes[indicesLat[1]]-lat) {
		latIdx = indicesLat[1]
	}

	if math.Abs(longitudes[lonIdx]-lon) > math.Abs(longitudes[indicesLon[1]]-lon) {
		lonIdx = indicesLon[1]
	}

	return latIdx, lonIdx
}
