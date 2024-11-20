package grids

import (
	"math"
)

type Grid interface {
	Size() int
	Latitudes() []float64
	GetNearestIndex(lat, lon float64) (int, int)
	GuessNearestIndex(lat, lon float64) (int, int)
	LonPointsAt(latIdx int) int
}

func GuessGridIndex(g Grid, lat, lon float64, mode ScanMode) int {
	latIdx, lonIdx := g.GuessNearestIndex(lat, lon)

	return GridIndexFromIndices(g, latIdx, lonIdx, mode)
}

func GridIndex(g Grid, lat, lon float64, mode ScanMode) int {
	latIdx, lonIdx := g.GetNearestIndex(lat, lon)

	return GridIndexFromIndices(g, latIdx, lonIdx, mode)
}

func GridPoint(g Grid, index int, mode ScanMode) (lat, lon float64) {
	if index < 0 || index >= g.Size() {
		return math.NaN(), math.NaN()
	}

	latitudesSize := len(g.Latitudes())
	var latIdx, lonIdx int

	if mode.IsConsecutiveJ() {
		// 连续J模式：先计算经度索引，再计算纬度索引
		lonIdx = index / latitudesSize
		latIdx = index % latitudesSize
	} else {
		// 连续I模式：先计算纬度索引，再计算经度索引
		totalIndex := 0
		for i := 0; i < latitudesSize; i++ {
			pointsInLat := g.LonPointsAt(i)
			nextTotal := totalIndex + pointsInLat
			if index < nextTotal {
				latIdx = i
				lonIdx = index - totalIndex
				break
			}
			totalIndex = nextTotal
		}
	}

	// 处理J方向扫描
	if mode.IsPositiveJ() {
		latIdx = latitudesSize - 1 - latIdx
	}

	// 处理负方向扫描
	if mode.IsNegativeI() {
		lonIdx = g.LonPointsAt(latIdx) - 1 - lonIdx
	}

	// 处理交替行
	if mode.HasOppositeRows() {
		if mode.IsConsecutiveJ() {
			// 连续J模式：奇数列反转纬度
			if lonIdx%2 == 1 {
				latIdx = latitudesSize - 1 - latIdx
			}
		} else {
			// 连续I模式：奇数行反转经度
			if latIdx%2 == 1 {
				lonIdx = g.LonPointsAt(latIdx) - 1 - lonIdx
			}
		}
	}

	// 获取纬度值
	lat = g.Latitudes()[latIdx]

	// 计算经度值
	lonStep := 360.0 / float64(g.LonPointsAt(latIdx))
	lon = float64(lonIdx) * lonStep

	return lat, lon
}

func GridIndexFromIndices(g Grid, latIdx, lonIdx int, mode ScanMode) int {
	if latIdx < 0 || latIdx >= len(g.Latitudes()) {
		return -1
	}

	lonPoints := g.LonPointsAt(latIdx)
	if lonIdx < 0 || lonIdx >= lonPoints {
		return -1
	}

	// 处理J方向扫描
	if mode.IsPositiveJ() {
		latIdx = len(g.Latitudes()) - 1 - latIdx
	}

	// 处理负方向扫描
	if mode.IsNegativeI() {
		lonIdx = lonPoints - 1 - lonIdx
	}

	// 处理交替行
	if mode.HasOppositeRows() {
		if mode.IsConsecutiveJ() {
			// 连续J模式：奇数列反转纬度
			if lonIdx%2 == 1 {
				latIdx = len(g.Latitudes()) - 1 - latIdx
			}
		} else {
			// 连续I模式：奇数行反转经度
			if latIdx%2 == 1 {
				lonIdx = lonPoints - 1 - lonIdx
			}
		}
	}

	if mode.IsConsecutiveJ() {
		// 连续J模式：计算每列的起始索引
		baseIndex := 0
		for i := 0; i < lonIdx; i++ {
			baseIndex += len(g.Latitudes())
		}
		return baseIndex + latIdx
	}

	// 连续I模式：计算每行的起始索引
	baseIndex := 0
	for i := 0; i < latIdx; i++ {
		baseIndex += g.LonPointsAt(i)
	}
	return baseIndex + lonIdx
}
