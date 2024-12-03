package grids

import (
	"math"
)

type Grid interface {
	Size() int
	Latitudes() []float64
	Longitudes() []float64
	GetNearestIndex(lat, lon float64) (int, int)
	GuessNearestIndex(lat, lon float64) (int, int)
}

func GuessGridIndex(g Grid, lat, lon float64, mode ScanMode) int {
	latIdx, lonIdx := g.GuessNearestIndex(lat, lon)

	return GridIndexFromIndices(g, latIdx, lonIdx, mode)
}

func GridIndex(g Grid, lat, lon float64, mode ScanMode) int {
	latIdx, lonIdx := g.GetNearestIndex(lat, lon)

	return GridIndexFromIndices(g, latIdx, lonIdx, mode)
}

func GridPoint(g Grid, index int, mode ScanMode) (lat, lon float64, ok bool) {
	if index < 0 || index >= g.Size() {
		return math.NaN(), math.NaN(), false
	}

	latitudesSize := len(g.Latitudes())
	longitudesSize := len(g.Longitudes())

	var latIdx, lonIdx int
	if mode.IsConsecutiveJ() {
		lonIdx = index / latitudesSize
		latIdx = index % latitudesSize
	} else {
		latIdx = index / longitudesSize
		lonIdx = index % longitudesSize
	}

	// 处理负方向扫描
	if mode.IsNegativeI() {
		lonIdx = longitudesSize - 1 - lonIdx
	}

	// 处理J方向扫描
	if mode.IsPositiveJ() {
		latIdx = latitudesSize - 1 - latIdx
	}

	// 处理交替行
	if mode.HasOppositeRows() && latIdx%2 == 1 {
		lonIdx = longitudesSize - 1 - lonIdx
	}

	return g.Latitudes()[latIdx], g.Longitudes()[lonIdx], true
}

func GridIndexFromIndices(g Grid, latIdx, lonIdx int, mode ScanMode) int {
	if latIdx < 0 || latIdx >= len(g.Latitudes()) || lonIdx < 0 || lonIdx >= len(g.Longitudes()) {
		return -1
	}

	// 处理-i方向扫描
	if mode.IsNegativeI() {
		lonIdx = len(g.Longitudes()) - 1 - lonIdx
	}

	// 处理+j方向扫描
	// 数组是从北到南排列，所以：
	// - 当是正向J扫描（从南到北）时，需要反转索引
	// - 当是负向J扫描（从北到南）时，不需要反转索引
	if mode.IsPositiveJ() {
		latIdx = len(g.Latitudes()) - 1 - latIdx
	}

	// 处理交替行
	if mode.HasOppositeRows() {
		switch {
		case mode.IsConsecutiveI() && latIdx%2 == 1:
			lonIdx = len(g.Longitudes()) - 1 - lonIdx
		case mode.IsConsecutiveJ() && lonIdx%2 == 1:
			latIdx = len(g.Latitudes()) - 1 - latIdx
		}
	}

	// 连续 J 方向
	if mode.IsConsecutiveJ() {
		return lonIdx*len(g.Latitudes()) + latIdx
	}

	// 连续 I 方向
	return latIdx*len(g.Longitudes()) + lonIdx
}
