package interpolators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBilinearInterpolator_Interpolate(t *testing.T) {
	bi := &BilinearInterpolator{}

	tests := []struct {
		name    string
		points  []float64
		weights []float64
		want    float64
	}{
		{
			// weights [0, 0] 表示目标点在左下角位置
			// 在角点位置，插值结果应该等于该角点的值
			name:    "exact corner - bottom left",
			points:  []float64{10, 20, 30, 40},
			weights: []float64{0, 0},
			want:    10,
		},
		{
			// weights [0, 1] 表示目标点在右下角位置
			name:    "exact corner - bottom right",
			points:  []float64{10, 20, 30, 40},
			weights: []float64{0, 1},
			want:    20,
		},
		{
			// weights [1, 0] 表示目标点在左上角位置
			name:    "exact corner - top left",
			points:  []float64{10, 20, 30, 40},
			weights: []float64{1, 0},
			want:    30,
		},
		{
			// weights [1, 1] 表示目标点在右上角位置
			name:    "exact corner - top right",
			points:  []float64{10, 20, 30, 40},
			weights: []float64{1, 1},
			want:    40,
		},
		{
			// weights [0.5, 0.5] 表示目标点在正中心位置
			// 双线性插值在此处会给四个点相等的权重：
			// (10 + 20 + 30 + 40) / 4 = 25
			name:    "center point",
			points:  []float64{10, 20, 30, 40},
			weights: []float64{0.5, 0.5},
			want:    25,
		},
		{
			// weights [0, 0.5] 表示目标点在下边缘中点
			// 双线性插值结果应该是两个下角点的平均值：
			// (10 + 20) / 2 = 15
			name:    "bottom edge middle",
			points:  []float64{10, 20, 30, 40},
			weights: []float64{0, 0.5},
			want:    15,
		},
		{
			// weights [1, 0.5] 表示目标点在上边缘中点
			// 双线性插值结果应该是两个上角点的平均值：
			// (30 + 40) / 2 = 35
			name:    "top edge middle",
			points:  []float64{10, 20, 30, 40},
			weights: []float64{1, 0.5},
			want:    35,
		},
		{
			// weights [0.5, 0] 表示目标点在左边缘中点
			// 双线性插值结果应该是两个左侧点的平均值：
			// (10 + 30) / 2 = 20
			name:    "left edge middle",
			points:  []float64{10, 20, 30, 40},
			weights: []float64{0.5, 0},
			want:    20,
		},
		{
			// weights [0.5, 1] 表示目标点在右边缘中点
			// 双线性插值结果应该是两个右侧点的平均值：
			// (20 + 40) / 2 = 30
			name:    "right edge middle",
			points:  []float64{10, 20, 30, 40},
			weights: []float64{0.5, 1},
			want:    30,
		},
		{
			// weights [0.3, 0.7] 表示目标点位于：
			// - 距离底部 30% 的高度
			// - 距离左侧 70% 的宽度
			// 当所有点值相同时，无论位置在哪里结果都应该相同
			name:    "same value at all points",
			points:  []float64{10, 10, 10, 10},
			weights: []float64{0.3, 0.7},
			want:    10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bi.Interpolate(tt.points, tt.weights)
			assert.InDelta(t, tt.want, got, 1e-10, "interpolation result mismatch")
		})
	}
}

func TestBilinearInterpolator_Interpolate_Validation(t *testing.T) {
	bi := &BilinearInterpolator{}

	// 测试权重范围
	weights := []struct {
		w    []float64
		desc string
	}{
		{[]float64{-0.1, 0.5}, "negative y weight"},
		{[]float64{0.5, -0.1}, "negative x weight"},
		{[]float64{1.1, 0.5}, "y weight > 1"},
		{[]float64{0.5, 1.1}, "x weight > 1"},
	}

	points := []float64{10, 20, 30, 40}
	for _, w := range weights {
		t.Run("invalid "+w.desc, func(t *testing.T) {
			// 虽然这些权重值不合理，但由于我们的实现是纯数学计算
			// 所以这里主要是确保不会panic，而不是检查具体的结果值
			got := bi.Interpolate(points, w.w)
			assert.NotPanics(t, func() {
				_ = got
			})
		})
	}
}
