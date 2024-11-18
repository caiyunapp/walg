package interpolation

import (
	"fmt"
	"testing"

	"github.com/scorix/walg/pkg/geo/grids"
	"github.com/scorix/walg/pkg/interpolation/interpolators"
	"github.com/stretchr/testify/assert"
)

// mockGrid 是一个用于测试的简单网格实现
// 实现了 2x2 的网格结构：
//
//	(31.0,120.0) --- (31.0,121.0)
//	      |               |
//	      |               |
//	(30.0,120.0) --- (30.0,121.0)
type mockGrid struct {
	lats []float64
	lons []float64
}

// 简化的网格接口实现
func (g *mockGrid) Size() int             { return len(g.lats) * len(g.lons) }
func (g *mockGrid) Latitudes() []float64  { return g.lats }
func (g *mockGrid) Longitudes() []float64 { return g.lons }
func (g *mockGrid) GetNearestIndex(lat, lon float64) (latIdx, lonIdx int) {
	return 0, 0
}
func (g *mockGrid) ScanMode() grids.ScanMode { return 0 }

// mockReader 模拟数据读取器
// 存储了2x2网格的四个角点值：
// - 左下角(30.0,120.0): 10.0
// - 右下角(30.0,121.0): 20.0
// - 左上角(31.0,120.0): 15.0
// - 右上角(31.0,121.0): 25.0
type mockReader struct {
	values map[int]float64
}

func (r *mockReader) ReadValueAt(timeStep, gridIndex int) (float64, error) {
	if timeStep < 0 {
		return 0, fmt.Errorf("negative time step: %d", timeStep)
	}
	if val, ok := r.values[gridIndex]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("invalid grid index: %d", gridIndex)
}

func TestGridInterpolator_InterpolateAt(t *testing.T) {
	// 创建测试网格和数据
	grid := &mockGrid{
		lats: []float64{30.0, 31.0},   // 1度间隔的纬度
		lons: []float64{120.0, 121.0}, // 1度间隔的经度
	}

	reader := &mockReader{
		values: map[int]float64{
			0: 10.0, // 左下角 (30.0, 120.0)
			1: 20.0, // 右下角 (30.0, 121.0)
			2: 15.0, // 左上角 (31.0, 120.0)
			3: 25.0, // 右上角 (31.0, 121.0)
		},
	}

	// 使用默认的双线性插值器
	interpolator := NewGridInterpolator(reader, grid, nil)

	tests := []struct {
		name     string
		timeStep int
		lat      float64
		lon      float64
		want     float64
	}{
		{
			// 网格中心点的插值
			// 位置：(30.5, 120.5)
			// 计算：四个角点的平均值 (10 + 20 + 15 + 25) / 4 = 17.5
			name:     "grid center point",
			timeStep: 0,
			lat:      30.5,
			lon:      120.5,
			want:     17.5,
		},
		{
			// 下边缘中点的插值
			// 位置：(30.0, 120.5)
			// 计算：下边两个点的平均值 (10 + 20) / 2 = 15.0
			name:     "bottom edge middle",
			timeStep: 0,
			lat:      30.0,
			lon:      120.5,
			want:     15.0,
		},
		{
			// 网格点的精确值
			// 位置：左下角 (30.0, 120.0)
			// 应直接返回该点的值 10.0
			name:     "exact grid point",
			timeStep: 0,
			lat:      30.0,
			lon:      120.0,
			want:     10.0,
		},
		{
			// 上边缘中点的插值
			// 位置：(31.0, 120.5)
			// 计算：上边两个点的平均值 (15 + 25) / 2 = 20.0
			name:     "top edge middle",
			timeStep: 0,
			lat:      31.0,
			lon:      120.5,
			want:     20.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := interpolator.InterpolateAt(tt.timeStep, tt.lat, tt.lon)
			assert.NoError(t, err)
			assert.InDelta(t, tt.want, got, 0.0001)
		})
	}
}

func TestGridInterpolator_InterpolateAt_Error(t *testing.T) {
	grid := &mockGrid{
		lats: []float64{30.0, 31.0},
		lons: []float64{120.0, 121.0},
	}

	reader := &mockReader{
		values: map[int]float64{},
	}

	// 使用默认的双线性插值器
	interpolator := NewGridInterpolator(reader, grid, nil)

	tests := []struct {
		name     string
		timeStep int
		lat      float64
		lon      float64
	}{
		{
			name:     "latitude out of bounds",
			timeStep: 0,
			lat:      32.0, // 超出网格范围的纬度
			lon:      120.5,
		},
		{
			name:     "longitude out of bounds",
			timeStep: 0,
			lat:      30.5,
			lon:      122.0, // 超出网格范围的经度
		},
		{
			name:     "negative time step",
			timeStep: -1,
			lat:      30.5,
			lon:      120.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := interpolator.InterpolateAt(tt.timeStep, tt.lat, tt.lon)
			assert.Error(t, err)
		})
	}
}

func TestGridInterpolator_WithDifferentInterpolators(t *testing.T) {
	grid := &mockGrid{
		lats: []float64{30.0, 31.0},
		lons: []float64{120.0, 121.0},
	}

	reader := &mockReader{
		values: map[int]float64{
			0: 10.0, // 左下角
			1: 20.0, // 右下角
			2: 15.0, // 左上角
			3: 25.0, // 右上角
		},
	}

	tests := []struct {
		name         string
		interpolator interpolators.Interpolator
		lat          float64
		lon          float64
		want         float64
	}{
		{
			name:         "bilinear interpolation",
			interpolator: &interpolators.BilinearInterpolator{},
			lat:          30.5,
			lon:          120.5,
			want:         17.5,
		},
		{
			name:         "idw interpolation",
			interpolator: interpolators.NewIDWInterpolator(2.0),
			lat:          30.5,
			lon:          120.5,
			want:         17.5,
		},
		{
			name:         "nearest interpolation",
			interpolator: &interpolators.NearestInterpolator{},
			lat:          30.2,
			lon:          120.2,
			want:         10.0, // 最近的是左下角点
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interpolator := NewGridInterpolator(reader, grid, tt.interpolator)
			got, err := interpolator.InterpolateAt(0, tt.lat, tt.lon)
			assert.NoError(t, err)
			assert.InDelta(t, tt.want, got, 0.0001)
		})
	}
}
