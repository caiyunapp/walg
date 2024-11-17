package interpolators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNearestInterpolator_Interpolate(t *testing.T) {
	ni := &NearestInterpolator{}

	tests := []struct {
		name    string
		points  []float64
		weights []float64
		want    float64
	}{
		{
			name:    "bottom left quadrant",
			points:  []float64{10, 20, 30, 40},
			weights: []float64{0.2, 0.2},
			want:    10, // 左下角最近
		},
		{
			name:    "bottom right quadrant",
			points:  []float64{10, 20, 30, 40},
			weights: []float64{0.2, 0.8},
			want:    20, // 右下角最近
		},
		{
			name:    "top left quadrant",
			points:  []float64{10, 20, 30, 40},
			weights: []float64{0.8, 0.2},
			want:    30, // 左上角最近
		},
		{
			name:    "top right quadrant",
			points:  []float64{10, 20, 30, 40},
			weights: []float64{0.8, 0.8},
			want:    40, // 右上角最近
		},
		{
			name:    "exact middle",
			points:  []float64{10, 20, 30, 40},
			weights: []float64{0.5, 0.5},
			want:    40, // 边界情况，按代码实现选择右上角
		},
		{
			name:    "same values",
			points:  []float64{10, 10, 10, 10},
			weights: []float64{0.3, 0.7},
			want:    10, // 所有点值相同
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ni.Interpolate(tt.points, tt.weights)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNearestInterpolator_EdgeCases(t *testing.T) {
	ni := &NearestInterpolator{}

	t.Run("boundary cases", func(t *testing.T) {
		tests := []struct {
			weights []float64
			want    int // 期望返回的点的索引
		}{
			{[]float64{0.5, 0.0}, 2}, // 左边界中点
			{[]float64{0.5, 1.0}, 3}, // 右边界中点
			{[]float64{0.0, 0.5}, 1}, // 下边界中点
			{[]float64{1.0, 0.5}, 3}, // 上边界中点
			{[]float64{0.5, 0.5}, 3}, // 正中心点
		}

		points := []float64{10, 20, 30, 40}
		for _, tt := range tests {
			got := ni.Interpolate(points, tt.weights)
			assert.Equal(t, points[tt.want], got)
		}
	})
}
