package interpolators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAverageInterpolator_Interpolate(t *testing.T) {
	avg := &AverageInterpolator{}

	tests := []struct {
		name    string
		points  []float64
		weights []float64 // weights 在这个实现中不影响结果
		want    float64
	}{
		{
			name:    "same values",
			points:  []float64{10, 10, 10, 10},
			weights: []float64{0.5, 0.5},
			want:    10,
		},
		{
			name:    "different values",
			points:  []float64{10, 20, 30, 40},
			weights: []float64{0.5, 0.5},
			want:    25, // (10 + 20 + 30 + 40) / 4
		},
		{
			name:    "zero values",
			points:  []float64{0, 0, 0, 0},
			weights: []float64{0.5, 0.5},
			want:    0,
		},
		{
			name:    "negative values",
			points:  []float64{-10, -20, -30, -40},
			weights: []float64{0.5, 0.5},
			want:    -25, // (-10 + -20 + -30 + -40) / 4
		},
		{
			name:    "mixed values",
			points:  []float64{-10, 10, -20, 20},
			weights: []float64{0.5, 0.5},
			want:    0, // (-10 + 10 + -20 + 20) / 4
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := avg.Interpolate(tt.points, tt.weights)
			assert.InDelta(t, tt.want, got, 1e-10, "interpolation result mismatch")
		})
	}
}

// 测试边界情况
func TestAverageInterpolator_Interpolate_EdgeCases(t *testing.T) {
	avg := &AverageInterpolator{}

	// 测试不同权重值不会影响结果
	t.Run("different weights same result", func(t *testing.T) {
		points := []float64{10, 20, 30, 40}
		weights1 := []float64{0.1, 0.1}
		weights2 := []float64{0.9, 0.9}

		result1 := avg.Interpolate(points, weights1)
		result2 := avg.Interpolate(points, weights2)

		assert.Equal(t, result1, result2, "results should be the same regardless of weights")
	})

	// 测试极大值
	t.Run("large values", func(t *testing.T) {
		points := []float64{1e6, 2e6, 3e6, 4e6}
		weights := []float64{0.5, 0.5}
		got := avg.Interpolate(points, weights)
		want := 2.5e6 // (1e6 + 2e6 + 3e6 + 4e6) / 4
		assert.InDelta(t, want, got, 1e-10)
	})
}
