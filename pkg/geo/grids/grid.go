package grids

import (
	"math"
)

type Grid interface {
	Size() int
	Latitudes() []float64
	Longitudes() []float64
	ScanMode() ScanMode
	GetNearestIndex(lat, lon float64) (int, int)
}

func GridIndex(g Grid, lat, lon float64) int {
	latIdx, lonIdx := g.GetNearestIndex(lat, lon)

	return GridIndexFromIndices(g, latIdx, lonIdx)
}

func GridPoint(g Grid, index int) (lat, lon float64) {
	if index < 0 || index >= g.Size() {
		return math.NaN(), math.NaN()
	}

	latitudesSize := len(g.Latitudes())
	longitudesSize := len(g.Longitudes())

	var latIdx, lonIdx int
	if g.ScanMode().IsConsecutiveJ() {
		lonIdx = index / latitudesSize
		latIdx = index % latitudesSize
	} else {
		latIdx = index / longitudesSize
		lonIdx = index % longitudesSize
	}

	// 处理负方向扫描
	if g.ScanMode().IsNegativeI() {
		lonIdx = longitudesSize - 1 - lonIdx
	}

	// 处理J方向扫描
	if g.ScanMode().IsPositiveJ() {
		latIdx = latitudesSize - 1 - latIdx
	}

	// 处理交替行
	if g.ScanMode().HasOppositeRows() && latIdx%2 == 1 {
		lonIdx = longitudesSize - 1 - lonIdx
	}

	return g.Latitudes()[latIdx], g.Longitudes()[lonIdx]
}

func GridIndexFromIndices(g Grid, latIdx, lonIdx int) int {
	if latIdx < 0 || latIdx >= len(g.Latitudes()) || lonIdx < 0 || lonIdx >= len(g.Longitudes()) {
		return -1
	}

	// 处理-i方向扫描
	if g.ScanMode().IsNegativeI() {
		lonIdx = len(g.Longitudes()) - 1 - lonIdx
	}

	// 处理+j方向扫描
	// 数组是从北到南排列，所以：
	// - 当是正向J扫描（从南到北）时，需要反转索引
	// - 当是负向J扫描（从北到南）时，不需要反转索引
	if g.ScanMode().IsPositiveJ() {
		latIdx = len(g.Latitudes()) - 1 - latIdx
	}

	// 处理交替行
	if g.ScanMode().HasOppositeRows() {
		switch {
		case g.ScanMode().IsConsecutiveI() && latIdx%2 == 1:
			lonIdx = len(g.Longitudes()) - 1 - lonIdx
		case g.ScanMode().IsConsecutiveJ() && lonIdx%2 == 1:
			latIdx = len(g.Latitudes()) - 1 - latIdx
		}
	}

	// 连续 J 方向
	if g.ScanMode().IsConsecutiveJ() {
		return lonIdx*len(g.Latitudes()) + latIdx
	}

	// 连续 I 方向
	return latIdx*len(g.Longitudes()) + lonIdx
}
