package latlon_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
			lat, lon, ok := grids.GridPoint(grid, tt.index, 0)
			require.True(t, ok)

			assert.Equal(t, tt.expectedLat, lat, "Latitude mismatch for index %d", tt.index)
			assert.Equal(t, tt.expectedLon, lon, "Longitude mismatch for index %d", tt.index)
		})
	}
}

func TestGridPoint_Wave0p16(t *testing.T) {
	grid := latlon.NewLatLonGrid(-15.0, 52.5, 0.0, 359.833, 0.166667, 0.166667)
	assert.Equal(t, 406, len(grid.Latitudes()))
	assert.Equal(t, 2160, len(grid.Longitudes()))
	assert.Equal(t, 876960, grid.Size())

	tests := []struct {
		lat, lon float64
		index    int
	}{
		{lat: 52.5, lon: 0.0, index: 0},
		{lat: 52.5, lon: 0.166667, index: 1},
		{lat: 52.5, lon: 0.333333, index: 2},
		{lat: 52.5, lon: 0.5, index: 3},
		{lat: 52.5, lon: 0.666667, index: 4},
		{lat: 52.5, lon: 0.833333, index: 5},
		{lat: 52.5, lon: 1, index: 6},
		{lat: 52.5, lon: 2, index: 12},
		{lat: 52.5, lon: 3, index: 18},
		{lat: 52.5, lon: 359, index: 359 * 6},
		{lat: 52.5, lon: 359.166667, index: 359*6 + 1},
		{lat: 52.5, lon: 359.333333, index: 359*6 + 2},
		{lat: 52.5, lon: 359.5, index: 359*6 + 3},
		{lat: 52.5, lon: 359.666667, index: 359*6 + 4},
		{lat: 52.5, lon: 359.833333, index: 359*6 + 5},
		{lat: -15.0, lon: 0, index: 874800},
		{lat: -15.0, lon: 359.833, index: 876959},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("lat=%.2f, lon=%.2f", tt.lat, tt.lon), func(t *testing.T) {
			lat, lon, ok := grids.GridPoint(grid, tt.index, 0)
			assert.InDelta(t, tt.lat, lat, 1e-2)
			assert.InDelta(t, tt.lon, lon, 1e-2)
			assert.True(t, ok)
		})
	}
}

func TestGridRoundTrip(t *testing.T) {
	grid := latlon.NewLatLonGrid(30.0, 35.0, 110.0, 115.0, 0.5, 0.5)

	// Test conversion from lat/lon to index and back
	originalLat, originalLon := 32.5, 112.5
	index := grids.GridIndex(grid, originalLat, originalLon, 0)
	recoveredLat, recoveredLon, ok := grids.GridPoint(grid, index, 0)
	require.True(t, ok)

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

func TestLatLonIsSphere(t *testing.T) {
	t.Run("Not Sphere", func(t *testing.T) {
		grid := latlon.NewLatLonGrid(-90, 90, 0, 180, 0.5, 0.5)
		assert.False(t, grid.IsSphere())
	})

	t.Run("Sphere", func(t *testing.T) {
		grid := latlon.NewLatLonGrid(-90, 90, 0, 359.75, 0.25, 0.25)
		assert.True(t, grid.IsSphere())
	})
}

func TestLatLonGetNearestIndex(t *testing.T) {
	grid := latlon.NewLatLonGrid(-90, 90, 0, 359.75, 0.25, 0.25)

	tests := []struct {
		name   string
		lat    float64
		lon    float64
		latIdx int
		lonIdx int
	}{
		{name: "Center", lat: 32.5, lon: 112.5, latIdx: 230, lonIdx: 450},
		{name: "Top", lat: 90, lon: 112.5, latIdx: 0, lonIdx: 450},
		{name: "Bottom", lat: -90, lon: 112.5, latIdx: 720, lonIdx: 450},
		{name: "Left", lat: 32.5, lon: 0, latIdx: 230, lonIdx: 0},
		{name: "Right", lat: 32.5, lon: 359.75, latIdx: 230, lonIdx: 1439},
		{name: "Top Left", lat: 90, lon: 0, latIdx: 0, lonIdx: 0},
		{name: "90W_Meridian", lat: 0, lon: -90, latIdx: 360, lonIdx: 1080},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			latIdx, lonIdx := grid.GetNearestIndex(tt.lat, tt.lon)

			assert.Equal(t, tt.latIdx, latIdx)
			assert.Equal(t, tt.lonIdx, lonIdx)

			latIdx, lonIdx = grid.GuessNearestIndex(tt.lat, tt.lon)
			assert.Equal(t, tt.latIdx, latIdx)
			assert.Equal(t, tt.lonIdx, lonIdx)
		})
	}
}

func BenchmarkGuessNearestIndex(b *testing.B) {
	grid := latlon.NewLatLonGrid(-90, 90, 0, 359.75, 0.25, 0.25)

	for i := 0; i < b.N; i++ {
		grid.GuessNearestIndex(32.5, 112.5)
	}
}
