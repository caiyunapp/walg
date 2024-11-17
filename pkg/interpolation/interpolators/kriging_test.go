package interpolators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKrigingInterpolator_Interpolate(t *testing.T) {
	ki := NewKrigingInterpolator(1.0, 1.0, 0.1)

	tests := []struct {
		name    string
		points  []float64
		weights []float64
		want    float64
	}{
		{
			// weights [0, 0] 表示目标点在左下角位置
			// 此时距离为0，应直接返回左下角点的值
			name:    "exact corner - bottom left",
			points:  []float64{10, 20, 30, 40},
			weights: []float64{0, 0},
			want:    10,
		},
		{
			// weights [1, 1] 表示目标点在右上角位置
			// 此时距离为0，应直接返回右上角点的值
			name:    "exact corner - top right",
			points:  []float64{10, 20, 30, 40},
			weights: []float64{1, 1},
			want:    40,
		},
		{
			// weights [0.5, 0.5] 表示目标点在正中心位置
			// 计算过程：
			// 1. 到四个角点的距离相等：sqrt(0.5² + 0.5²) ≈ 0.707
			// 2. 变异函数值相等：γ(h) = 0.1 + 1.0 * (1.5*0.707 - 0.5*0.707³) ≈ 0.85
			// 3. 权重相等：1/0.85 ≈ 1.176
			// 4. 结果为算术平均值：(10 + 20 + 30 + 40)/4 = 25
			name:    "center point",
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
			points:  []float64{10, 10, 10, 10},
			weights: []float64{0.3, 0.7},
			want:    10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ki.Interpolate(tt.points, tt.weights)
			assert.InDelta(t, tt.want, got, 1e-10)
		})
	}
}

func TestKrigingInterpolator_Parameters(t *testing.T) {
	// weights [0.2, 0.2] 表示目标点位于：
	// - 距离底部 20% 的高度
	// - 距离左侧 20% 的宽度
	// 此位置到四个角点的欧几里得距离：
	// - 左下角: sqrt(0.2² + 0.2²) ≈ 0.283
	// - 右下角: sqrt(0.8² + 0.2²) ≈ 0.825
	// - 左上角: sqrt(0.2² + 0.8²) ≈ 0.825
	// - 右上角: sqrt(0.8² + 0.8²) ≈ 1.131
	//
	// k1 (Range=1.0, Nugget=0.1): 所有距离都在变程内，变异函数完整计算
	// k2 (Range=0.5, Nugget=0.1): 部分距离超出变程，使用基台值
	// k3 (Range=1.0, Nugget=0.5): 较大的块金值降低了空间相关性
	points := []float64{10, 20, 30, 40}
	weights := []float64{0.2, 0.2}

	k1 := NewKrigingInterpolator(1.0, 1.0, 0.1)
	k2 := NewKrigingInterpolator(1.0, 0.5, 0.1)
	k3 := NewKrigingInterpolator(1.0, 1.0, 0.5)

	v1 := k1.Interpolate(points, weights)
	v2 := k2.Interpolate(points, weights)
	v3 := k3.Interpolate(points, weights)

	// Check if different parameters produce different results
	assert.NotEqual(t, v1, v2, "range parameter should affect the interpolation result")
	assert.NotEqual(t, v1, v3, "nugget parameter should affect the interpolation result")

	// Results should be closer to the nearest point (10)
	assert.Greater(t, v1, 10.0, "interpolation result should be greater than nearest point")
	assert.Less(t, v1, 25.0, "interpolation result should be less than center value")
}
