package latlon_test

import (
	"fmt"
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
	assert.Equal(t, []float64{35.0, 34.5, 34.0, 33.5, 33.0, 32.5, 32.0, 31.5, 31.0, 30.5, 30.0}, lats)

	// Test longitude array
	lons := grid.Longitudes()
	assert.Equal(t, []float64{110.0, 110.5, 111.0, 111.5, 112.0, 112.5, 113.0, 113.5, 114.0, 114.5, 115.0}, lons)
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
			name:     "Top Left Corner",
			lat:      35.0,
			lon:      110.0,
			expected: 0,
		},
		{
			name:     "Bottom Right Corner",
			lat:      30.0,
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
			expected: 26,
		},
		{
			name:     "Cyclic Mapping - Above Range",
			lat:      36.0,  // 36.0 -> 31.0 (循环映射：36.0 - 5.0 = 31.0)
			lon:      116.0, // 116.0 -> 111.0 (循环映射：116.0 - 5.0 = 111.0)
			expected: 90,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := grids.GridIndex(grid, tt.lat, tt.lon, 0)
			assert.Equal(t, tt.expected, result, "For point (%.2f, %.2f)", tt.lat, tt.lon)
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
			name:     "Top Left Corner",
			latIdx:   0,
			lonIdx:   0,
			expected: 0,
		},
		{
			name:     "Bottom Right Corner",
			latIdx:   10,
			lonIdx:   10,
			expected: 120,
		},
		{
			name:     "Middle Point",
			latIdx:   5,
			lonIdx:   5,
			expected: 60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := grids.GridIndexFromIndices(grid, tt.latIdx, tt.lonIdx, 0)
			assert.Equal(t, tt.expected, result, "For indices (lat=%d, lon=%d)", tt.latIdx, tt.lonIdx)
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
			name:        "Top Left Corner",
			index:       0,
			expectedLat: 35.0,
			expectedLon: 110.0,
		},
		{
			name:        "Bottom Right Corner",
			index:       120,
			expectedLat: 30.0,
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
			lat, lon := grids.GridPoint(grid, tt.index, 0)
			assert.Equal(t, tt.expectedLat, lat, "Latitude mismatch for index %d", tt.index)
			assert.Equal(t, tt.expectedLon, lon, "Longitude mismatch for index %d", tt.index)
		})
	}
}

func TestGridRoundTrip(t *testing.T) {
	grid := latlon.NewLatLonGrid(30.0, 35.0, 110.0, 115.0, 0.5, 0.5)

	// Test conversion from lat/lon to index and back
	originalLat, originalLon := 32.5, 112.5
	index := grids.GridIndex(grid, originalLat, originalLon, 0)
	recoveredLat, recoveredLon := grids.GridPoint(grid, index, 0)

	assert.Equal(t, originalLat, recoveredLat, "Latitude should remain unchanged")
	assert.Equal(t, originalLon, recoveredLon, "Longitude should remain unchanged")
}

func BenchmarkNewLatLonGrid(b *testing.B) {
	for _, latStep := range []float64{0.5, 0.25, 0.125} {
		b.Run(fmt.Sprintf("latStep=%.3f", latStep), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				latlon.NewLatLonGrid(-90, 90, -180, 180, latStep, latStep)
			}
		})
	}
}
