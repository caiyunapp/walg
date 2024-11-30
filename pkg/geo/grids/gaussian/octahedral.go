package gaussian

import (
	"cmp"
	"fmt"
	"math"

	"github.com/scorix/walg/pkg/geo/distance"
	"github.com/scorix/walg/pkg/geo/grids"
	"golang.org/x/sync/singleflight"
)

var octahedralCache = make(map[int]*octahedral)
var octahedralCacheGroup singleflight.Group

type octahedral struct {
	n         int
	latitudes []float64
	lonPoints []int
}

func NewOctahedral(n int) *octahedral {
	cacheKey := fmt.Sprintf("O%d", n)

	res, _, _ := octahedralCacheGroup.Do(cacheKey, func() (any, error) {
		if grid, ok := octahedralCache[n]; ok {
			return grid, nil
		}

		grid := newOctahedral(n)
		octahedralCache[n] = grid
		return grid, nil
	})

	return res.(*octahedral)
}

func newOctahedral(n int) *octahedral {
	o := &octahedral{
		n: n,
	}

	o.latitudes = o.calcLatitudes()
	o.lonPoints = o.calcLonPoints()

	return o
}

// Grid interface implementation
func (g *octahedral) Size() int {
	total := 0
	for i := 0; i < len(g.latitudes); i++ {
		total += g.LonPointsAt(i) * len(g.latitudes)
	}
	return total
}

func (g *octahedral) Latitudes() []float64 {
	return g.latitudes
}

func (g *octahedral) Longitudes() []float64 {
	// 返回赤道处的最大经度点数
	maxPoints := 0
	for _, points := range g.lonPoints {
		if points > maxPoints {
			maxPoints = points
		}
	}

	lons := make([]float64, maxPoints)
	dlon := 360.0 / float64(maxPoints)
	for i := 0; i < maxPoints; i++ {
		lons[i] = float64(i) * dlon
	}
	return lons
}

func (g *octahedral) LongitudesOnLat(lat float64) []float64 {
	indicesLat := grids.FindNearestIndices(lat, g.latitudes)
	lat0, lat1 := g.latitudes[indicesLat[0]], g.latitudes[indicesLat[1]]

	var nearestLatIdx int
	if math.Abs(lat0-lat) <= math.Abs(lat1-lat) {
		nearestLatIdx = indicesLat[0]
	} else {
		nearestLatIdx = indicesLat[1]
	}

	nlon := g.LonPointsAt(nearestLatIdx)
	lons := make([]float64, nlon)
	dlon := 360.0 / float64(nlon)
	for i := 0; i < nlon; i++ {
		lons[i] = float64(i) * dlon
	}
	return lons
}

func (g *octahedral) LonPoints() []int {
	return g.lonPoints
}

func (g *octahedral) GuessNearestIndex(lat, lon float64) (int, int) {
	latitudes := g.Latitudes()
	indicesLat := grids.FindNearestIndices(lat, latitudes)
	latIdx := indicesLat[0]

	if math.Abs(latitudes[latIdx]-lat) > math.Abs(latitudes[indicesLat[1]]-lat) {
		latIdx = indicesLat[1]
	}

	lon = math.Mod(lon+360.0, 360.0)
	longitudes := g.LongitudesOnLat(latitudes[latIdx])
	indicesLon := grids.FindNearestIndices(lon, longitudes)
	lonIdx := indicesLon[0]

	if math.Abs(longitudes[lonIdx]-lon) > math.Abs(longitudes[indicesLon[1]]-lon) {
		lonIdx = indicesLon[1]
	}

	return latIdx, lonIdx
}

func (g *octahedral) GetNearestIndex(lat, lon float64) (int, int) {
	latitudes := g.Latitudes()
	indicesLat := grids.FindNearestIndices(lat, latitudes)

	latIdx := indicesLat[0]
	lonIdx := 0
	dist := math.MaxFloat64

	// iterations 3 is enough for comparing the minimum distance
	const iterations = 3

	// 标准化经度到 [0, 360)
	lon = math.Mod(lon+360.0, 360.0)

	for _, i := range indicesLat {
		longitudes := g.LongitudesOnLat(latitudes[i])
		indicesLon := grids.FindNearestIndices(lon, longitudes)

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

func (g *octahedral) calcLatitudes() []float64 {
	return gaussLegendreZeros(g.n * 2)
}

func (g *octahedral) calcLonPoints() []int {
	points := make([]int, len(g.latitudes))

	for i, lat := range g.latitudes {
		// Convert to colatitude (0 at North Pole)
		colat := 90.0 - lat
		colatRad := colat * math.Pi / 180.0

		// Calculate number of longitude points using sine of colatitude
		nlon := 4 * g.n * int(math.Round(math.Sin(colatRad)))

		// Ensure minimum number of points
		if nlon < 4 {
			nlon = 4
		}

		points[i] = nlon
	}

	return points
}

func (g *octahedral) LonPointsAt(latIdx int) int {
	if latIdx < 0 || latIdx >= len(g.lonPoints) {
		return 0
	}
	return g.lonPoints[latIdx]
}
