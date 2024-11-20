package grids

import (
	"math"
)

type NearestGrids struct {
	g         Grid
	latitudes []float64
}

func NewNearestGrids(g Grid) *NearestGrids {
	return &NearestGrids{
		g:         g,
		latitudes: g.Latitudes(),
	}
}

func (ng *NearestGrids) NearestGrids(lat, lon float64, mode ScanMode) []int {
	// Find nearest latitude indices
	latIdx := FindNearestIndices(lat, ng.latitudes)

	// 标准化经度到 [0, 360)
	lon = math.Mod(lon+360.0, 360.0)

	result := make([]int, 0, 4) // 通常每个点最多有4个最近点
	for _, i := range latIdx {
		// 计算该纬度的经度步长
		lonStep := 360.0 / float64(ng.g.LonPointsAt(i))

		// 找到最近的经度索引
		lonIdxFloat := lon / lonStep
		lonIdx := int(math.Round(lonIdxFloat))

		// 检查前后两个经度点
		lonIndices := []int{
			(lonIdx - 1 + ng.g.LonPointsAt(i)) % ng.g.LonPointsAt(i),
			lonIdx % ng.g.LonPointsAt(i),
			(lonIdx + 1) % ng.g.LonPointsAt(i),
		}

		// 将所有可能的组合添加到结果中
		for _, j := range lonIndices {
			gridIndex := GridIndexFromIndices(ng.g, i, j, mode)
			if gridIndex >= 0 { // 只添加有效的索引
				result = append(result, gridIndex)
			}
		}
	}

	return result
}

// findNearestIndices returns the indices of the two closest values in a sorted slice
func FindNearestIndices(value float64, sorted []float64) [2]int {
	if len(sorted) == 0 {
		return [2]int{}
	}

	// Define comparison function based on array direction
	isOnLeft := func(a, b float64) bool {
		if sorted[0] > sorted[len(sorted)-1] {
			return a > b // descending order
		}
		return a < b // ascending order
	}

	// Handle edge cases
	if !isOnLeft(sorted[0], value) {
		return [2]int{0, 0}
	}
	if !isOnLeft(value, sorted[len(sorted)-1]) {
		return [2]int{len(sorted) - 1, len(sorted) - 1}
	}

	// Binary search to find the closest index
	left, right := 0, len(sorted)-1
	for right-left > 1 {
		mid := (left + right) / 2
		if sorted[mid] == value {
			return [2]int{mid, mid}
		}
		if isOnLeft(sorted[mid], value) {
			left = mid
		} else {
			right = mid
		}
	}

	// Return both surrounding indices
	return [2]int{left, right}
}

func (ng *NearestGrids) NearestGrid(lat, lon float64, mode ScanMode) int {
	return GridIndex(ng.g, lat, lon, mode)
}
