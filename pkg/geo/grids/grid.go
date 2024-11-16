package grids

type Grid interface {
	Size() int
	Latitudes() []float64
	Longitudes() []float64
	GridIndex(lat, lon float64) int
	GridIndexFromIndices(latIdx, lonIdx int) int
	GridPoint(index int) (lat, lon float64)
}
