package grids

type NearestGrids struct {
	g          Grid
	latitudes  []float64
	longitudes []float64
}

func NewNearestGrids(g Grid) *NearestGrids {
	return &NearestGrids{
		g:          g,
		latitudes:  g.Latitudes(),
		longitudes: g.Longitudes(),
	}
}

func (ng *NearestGrids) NearestGrids(lat, lon float64) []int {
	// Find nearest latitude indices
	latIdx := FindNearestIndices(lat, ng.latitudes)

	// Find nearest longitude indices
	lonIdx := FindNearestIndices(lon, ng.longitudes)

	// Calculate all possible combinations of nearby grid points
	result := make([]int, 0, len(latIdx)*len(lonIdx))
	for _, i := range latIdx {
		for _, j := range lonIdx {
			gridIndex := GridIndexFromIndices(ng.g, i, j)
			result = append(result, gridIndex)
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

func (ng *NearestGrids) NearestGrid(lat, lon float64) int {
	return GridIndex(ng.g, lat, lon)
}
