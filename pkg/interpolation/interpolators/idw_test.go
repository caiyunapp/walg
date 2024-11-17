package interpolators

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIDWInterpolator_Interpolate(t *testing.T) {
	// 基本功能测试用例
	tests := []struct {
		name    string
		power   float64
		points  []float64
		weights []float64
		want    float64
	}{
		{
			// weights [0, 0] 表示目标点在左下角位置
			// 此时距离为0，应直接返回左下角点的值
			name:    "exact corner - bottom left",
			power:   2.0,
			points:  []float64{10, 20, 30, 40},
			weights: []float64{0, 0},
			want:    10,
		},
		{
			// weights [1, 1] 表示目标点在右上角位置
			// 此时距离为0，应直接返回右上角点的值
			name:    "exact corner - top right",
			power:   2.0,
			points:  []float64{10, 20, 30, 40},
			weights: []float64{1, 1},
			want:    40,
		},
		{
			// weights [0.5, 0.5] 表示目标点在正中心位置
			// 计算过程：
			// 1. 到四个角点的距离相等：sqrt(0.5^2 + 0.5^2) ≈ 0.707
			// 2. 距离相等导致权重相等，结果为算术平均值
			name:    "center point with power 2",
			power:   2.0,
			points:  []float64{10, 20, 30, 40},
			weights: []float64{0.5, 0.5},
			want:    25,
		},
		{
			// weights [0.3, 0.7] 表示目标点位于：
			// - 距离底部 30% 的高度
			// - 距离左侧 70% 的宽度
			// 当所有点值相同时，无论位置在哪里结果都应该相同
			name:    "same values",
			power:   2.0,
			points:  []float64{10, 10, 10, 10},
			weights: []float64{0.3, 0.7},
			want:    10,
		},
		{
			// weights [0.5, 0.5] 表示目标点在正中心位置
			// power=1 时的计算过程：
			// 1. 权重是距离的倒数（不是平方的倒数）
			// 2. 由于距离相等，权重也相等
			// 3. 结果是算术平均值
			name:    "linear weight (power=1)",
			power:   1.0,
			points:  []float64{10, 20, 30, 40},
			weights: []float64{0.5, 0.5},
			want:    25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idw := NewIDWInterpolator(tt.power)
			got := idw.Interpolate(tt.points, tt.weights)
			assert.InDelta(t, tt.want, got, 1e-10)
		})
	}
}

func TestIDWInterpolator_PowerEffect(t *testing.T) {
	// weights [0.1, 0.1] 表示目标点位于：
	// - 距离底部 10% 的高度
	// - 距离左侧 10% 的宽度
	// 此位置到四个角点的欧几里得距离：
	// 左下角: sqrt(0.1^2 + 0.1^2) ≈ 0.141
	// 右下角: sqrt(0.1^2 + 0.9^2) ≈ 0.906
	// 左上角: sqrt(0.9^2 + 0.1^2) ≈ 0.906
	// 右上角: sqrt(0.9^2 + 0.9^2) ≈ 1.273
	points := []float64{10, 20, 30, 40}
	weights := []float64{0.1, 0.1}

	powers := []float64{1.0, 2.0, 4.0, 8.0}
	var results []float64

	for _, power := range powers {
		idw := NewIDWInterpolator(power)
		result := idw.Interpolate(points, weights)
		results = append(results, result)
	}

	for i := 1; i < len(results); i++ {
		assert.Less(t,
			math.Abs(results[i]-10),
			math.Abs(results[i-1]-10),
			"higher power should result in closer value to nearest point")
	}
}

func TestIDWInterpolator_EdgeCases(t *testing.T) {
	idw := NewIDWInterpolator(2.0)

	t.Run("very close to point", func(t *testing.T) {
		points := []float64{10, 20, 30, 40}
		weights := []float64{0.0000001, 0.0000001}
		got := idw.Interpolate(points, weights)
		assert.InDelta(t, points[0], got, 1e-10)
	})

	t.Run("extreme power values", func(t *testing.T) {
		points := []float64{10, 20, 30, 40}
		weights := []float64{0.1, 0.1}

		idw100 := NewIDWInterpolator(100.0)
		result := idw100.Interpolate(points, weights)
		assert.InDelta(t, 10.0, result, 0.1)
	})
}

func TestIDWInterpolator_Constructor(t *testing.T) {
	tests := []struct {
		name  string
		power float64
	}{
		{"zero power", 0.0},
		{"negative power", -1.0},
		{"very large power", 100.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idw := NewIDWInterpolator(tt.power)
			assert.Equal(t, tt.power, idw.Power)
		})
	}
}
