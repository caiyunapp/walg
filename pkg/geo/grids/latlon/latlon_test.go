package latlon_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/scorix/walg/pkg/geo/grids"
	"github.com/scorix/walg/pkg/geo/grids/latlon"
)

func TestNewLatLonGrid(t *testing.T) {
	// Create a simple 10x10 grid
	grid := latlon.NewLatLonGrid(
		30.0,  // minLat
		35.0,  // maxLat
		110.0, // minLon
		115.0, // maxLon
		0.5,   // latStep
		0.5,   // lonStep
	)

	assert.Equal(t, 11*11, grid.Size(), "Grid size should be 11x11=121")

	// Test latitude array
	lats := grid.Latitudes()
	assert.Equal(t, 11, len(lats), "Should have 11 latitude values")
	assert.Equal(t, 30.0, lats[0], "First latitude should be 30.0")
	assert.Equal(t, 35.0, lats[len(lats)-1], "Last latitude should be 35.0")

	// Test longitude array
	lons := grid.Longitudes()
	assert.Equal(t, 11, len(lons), "Should have 11 longitude values")
	assert.Equal(t, 110.0, lons[0], "First longitude should be 110.0")
	assert.Equal(t, 115.0, lons[len(lons)-1], "Last longitude should be 115.0")
}

func TestGridIndex(t *testing.T) {
	grid := latlon.NewLatLonGrid(30.0, 35.0, 110.0, 115.0, 0.5, 0.5)

	tests := []struct {
		name     string
		lat      float64
		lon      float64
		expected int
	}{
		{
			name:     "Bottom Left Corner",
			lat:      30.0,
			lon:      110.0,
			expected: 0,
		},
		{
			name:     "Top Right Corner",
			lat:      35.0,
			lon:      115.0,
			expected: 120,
		},
		{
			name:     "Center Point",
			lat:      32.5,
			lon:      112.5,
			expected: 60,
		},
		{
			name:     "Cyclic Mapping - Below Range",
			lat:      29.0, // 29.0 -> 34.0 (循环映射：29.0 + 5.0 = 34.0)
			lon:      112.0,
			expected: 92,
		},
		{
			name:     "Cyclic Mapping - Above Range",
			lat:      32.0,
			lon:      116.0, // 116.0 -> 111.0 (循环映射：116.0 - 5.0 = 111.0)
			expected: 46,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := grid.GridIndex(tt.lat, tt.lon)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGridIndexFromIndices(t *testing.T) {
	grid := latlon.NewLatLonGrid(30.0, 35.0, 110.0, 115.0, 0.5, 0.5)

	tests := []struct {
		name     string
		latIdx   int
		lonIdx   int
		expected int
	}{
		{
			name:     "Bottom Left Corner",
			latIdx:   0,
			lonIdx:   0,
			expected: 0,
		},
		{
			name:     "Top Right Corner",
			latIdx:   10,
			lonIdx:   10,
			expected: 120,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := grid.GridIndexFromIndices(tt.latIdx, tt.lonIdx)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGridPoint(t *testing.T) {
	grid := latlon.NewLatLonGrid(30.0, 35.0, 110.0, 115.0, 0.5, 0.5)

	tests := []struct {
		name        string
		index       int
		expectedLat float64
		expectedLon float64
	}{
		{
			name:        "Bottom Left Corner",
			index:       0,
			expectedLat: 30.0,
			expectedLon: 110.0,
		},
		{
			name:        "Top Right Corner",
			index:       120,
			expectedLat: 35.0,
			expectedLon: 115.0,
		},
		{
			name:        "Center Point",
			index:       60,
			expectedLat: 32.5,
			expectedLon: 112.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lat, lon := grid.GridPoint(tt.index)
			assert.Equal(t, tt.expectedLat, lat)
			assert.Equal(t, tt.expectedLon, lon)
		})
	}
}

func TestGridRoundTrip(t *testing.T) {
	grid := latlon.NewLatLonGrid(30.0, 35.0, 110.0, 115.0, 0.5, 0.5)

	// Test conversion from lat/lon to index and back
	originalLat, originalLon := 32.5, 112.5
	index := grid.GridIndex(originalLat, originalLon)
	recoveredLat, recoveredLon := grid.GridPoint(index)

	assert.Equal(t, originalLat, recoveredLat, "Latitude should remain unchanged")
	assert.Equal(t, originalLon, recoveredLon, "Longitude should remain unchanged")
}

func TestGetLatitudeIndex(t *testing.T) {
	grid := latlon.NewLatLonGrid(30.0, 35.0, 110.0, 115.0, 0.5, 0.5)

	tests := []struct {
		name     string
		lat      float64
		expected int
	}{
		{
			name:     "Minimum Latitude",
			lat:      30.0,
			expected: 0,
		},
		{
			name:     "Maximum Latitude",
			lat:      35.0,
			expected: 10,
		},
		{
			name:     "Middle Latitude",
			lat:      32.5,
			expected: 5,
		},
		{
			name:     "One Interval Below",
			lat:      29.9,
			expected: 9, // 29.9 -> 34.9 (循环映射：29.9 + 5.0 = 34.9)
		},
		{
			name:     "One Interval Above",
			lat:      35.1,
			expected: 0, // 35.1 -> 30.1 (循环映射：35.1 - 5.0 = 30.1)
		},
		{
			name:     "Multiple Intervals Below",
			lat:      20.0,
			expected: 0, // 20.0 -> 30.0 (循环映射：20.0 + 3*5.0 = 35.0 - 5.0 = 30.0)
		},
		{
			name:     "Multiple Intervals Above",
			lat:      45.0,
			expected: 10, // 45.0 -> 35.0 (循环映射：45.0 - 2*5.0 = 35.0)
		},
		{
			name:     "Exact Multiple Intervals",
			lat:      40.0,
			expected: 10, // 40.0 -> 35.0 (循环映射：40.0 - 1*5.0 = 35.0)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := grid.GetLatitudeIndex(tt.lat)
			assert.Equal(t, tt.expected, result, "For latitude %.2f", tt.lat)
		})
	}
}

func TestGetLongitudeIndex(t *testing.T) {
	grid := latlon.NewLatLonGrid(30.0, 35.0, 110.0, 115.0, 0.5, 0.5)

	tests := []struct {
		name     string
		lat      float64
		lon      float64
		expected int
	}{
		{
			name:     "Minimum Longitude",
			lat:      30.0,
			lon:      110.0,
			expected: 0,
		},
		{
			name:     "Maximum Longitude",
			lat:      30.0,
			lon:      115.0,
			expected: 10,
		},
		{
			name:     "Middle Longitude",
			lat:      30.0,
			lon:      112.5,
			expected: 5,
		},
		{
			name:     "One Interval Below",
			lat:      30.0,
			lon:      109.9,
			expected: 9, // 109.9 -> 114.9 (循环映射：109.9 + 5.0 = 114.9)
		},
		{
			name:     "One Interval Above",
			lat:      30.0,
			lon:      115.1,
			expected: 0, // 115.1 -> 110.1 (循环映射：115.1 - 5.0 = 110.1)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := grid.GetLongitudeIndex(tt.lat, tt.lon)
			assert.Equal(t, tt.expected, result, "For longitude %.2f", tt.lon)
		})
	}
}

func TestGridWithScanMode(t *testing.T) {
	tests := []struct {
		name     string
		grid     grids.Grid
		lat      float64
		lon      float64
		expected int
	}{
		{
			name: "Default Scan Mode",
			grid: latlon.NewLatLonGrid(30.0, 35.0, 110.0, 115.0, 0.5, 0.5),
			lat:  32.0,
			lon:  112.0,
			// latCount = (35.0-30.0)/0.5 + 1 = 11
			// lonCount = (115.0-110.0)/0.5 + 1 = 11
			// latIdx = (32.0-30.0)/0.5 = 4
			// lonIdx = (112.0-110.0)/0.5 = 4
			// index = latIdx*11 + lonIdx = 4*11 + 4 = 44 + 4 = 48
			expected: 48,
		},
		{
			name: "Negative I Direction",
			grid: latlon.NewLatLonGrid(30.0, 35.0, 110.0, 115.0, 0.5, 0.5,
				latlon.WithScanMode(grids.ScanModeNegativeI)),
			lat: 32.0,
			lon: 112.0,
			// latCount = (35.0-30.0)/0.5 + 1 = 11
			// lonCount = (115.0-110.0)/0.5 + 1 = 11
			// latIdx = (32.0-30.0)/0.5 = 4
			// lonIdx = (112.0-110.0)/0.5 = 4
			// lonIdx_reversed = 10 - 4 = 6
			// index = latIdx*11 + lonIdx_reversed = 4*11 + 6 = 44 + 6 = 50
			expected: 50,
		},
		{
			name: "Positive J Direction",
			grid: latlon.NewLatLonGrid(30.0, 35.0, 110.0, 115.0, 0.5, 0.5,
				latlon.WithScanMode(grids.ScanModePositiveJ)),
			lat: 32.0,
			lon: 112.0,
			// latCount = (35.0-30.0)/0.5 + 1 = 11
			// lonCount = (115.0-110.0)/0.5 + 1 = 11
			// latIdx = (32.0-30.0)/0.5 = 4
			// latIdx_reversed = 10 - 4 = 6
			// lonIdx = (112.0-110.0)/0.5 = 4
			// index = latIdx_reversed*11 + lonIdx = 6*11 + 4 = 66 + 4 = 70
			expected: 70,
		},
		{
			name: "Consecutive J",
			grid: latlon.NewLatLonGrid(30.0, 35.0, 110.0, 115.0, 0.5, 0.5,
				latlon.WithScanMode(grids.ScanModeConsecutiveJ)),
			lat: 32.0,
			lon: 112.0,
			// latCount = (35.0-30.0)/0.5 + 1 = 11
			// lonCount = (115.0-110.0)/0.5 + 1 = 11
			// latIdx = (32.0-30.0)/0.5 = 4
			// lonIdx = (112.0-110.0)/0.5 = 4
			// index = lonIdx*11 + latIdx = 4*11 + 4 = 44 + 4 = 48
			expected: 48,
		},
		{
			name: "Opposite Rows",
			grid: latlon.NewLatLonGrid(30.0, 35.0, 110.0, 115.0, 0.5, 0.5,
				latlon.WithScanMode(grids.ScanModeOppositeRows)),
			lat: 32.5,
			lon: 112.0,
			// latCount = (35.0-30.0)/0.5 + 1 = 11
			// lonCount = (115.0-110.0)/0.5 + 1 = 11
			// latIdx = (32.5-30.0)/0.5 = 5 (奇数行)
			// lonIdx = (112.0-110.0)/0.5 = 4
			// lonIdx_reversed = 10 - 4 = 6 (因为是奇数行，所以反转)
			// index = latIdx*11 + lonIdx_reversed = 5*11 + 6 = 55 + 6 = 61
			expected: 61,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试 GridIndex
			idx := tt.grid.GridIndex(tt.lat, tt.lon)
			assert.Equal(t, tt.expected, idx, "GridIndex mismatch")

			// 测试 GridPoint
			lat, lon := tt.grid.GridPoint(tt.expected)
			assert.Equal(t, tt.lat, lat, "Latitude mismatch")
			assert.Equal(t, tt.lon, lon, "Longitude mismatch")

			// 测试往返转换
			roundTripIdx := tt.grid.GridIndex(lat, lon)
			assert.Equal(t, tt.expected, roundTripIdx, "Round trip index mismatch")
		})
	}
}
