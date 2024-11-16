package grids

type Grid interface {
	Size() int
	GridIndex(lat, lon float64) int
	GridPoint(index int) (lat, lon float64)
}
