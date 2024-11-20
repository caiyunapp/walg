package interpolation

import (
	"github.com/scorix/walg/pkg/geo/grids"
	"github.com/scorix/walg/pkg/interpolation/interpolators"
)

// ValueReader 定义了获取网格数据的接口
type ValueReader interface {
	// ReadValueAt 读取指定时间步和网格索引的值
	ReadValueAt(timeStep, gridIndex int) (float64, error)
}

// GridInterpolator 网格插值器
type GridInterpolator struct {
	reader       ValueReader
	grid         grids.Grid
	scanningMode grids.ScanMode
	interpolator interpolators.Interpolator
}

// NewGridInterpolator 创建新的网格插值器
func NewGridInterpolator(reader ValueReader, grid grids.Grid, scanningMode grids.ScanMode, interpolator interpolators.Interpolator) *GridInterpolator {
	if interpolator == nil {
		interpolator = &interpolators.BilinearInterpolator{} // 默认使用双线性插值
	}
	return &GridInterpolator{
		reader:       reader,
		grid:         grid,
		scanningMode: scanningMode,
		interpolator: interpolator,
	}
}

// InterpolateAt 在指定时间步和位置进行插值
func (g *GridInterpolator) InterpolateAt(timeStep int, lat, lon float64) (float64, error) {
	// 获取网格索引
	latIdx, lonIdx := g.grid.GetNearestIndex(lat, lon)

	// 获取四个相邻点的网格索引
	indices := []int{
		grids.GridIndexFromIndices(g.grid, latIdx, lonIdx, g.scanningMode),
		grids.GridIndexFromIndices(g.grid, latIdx, lonIdx+1, g.scanningMode),
		grids.GridIndexFromIndices(g.grid, latIdx+1, lonIdx, g.scanningMode),
		grids.GridIndexFromIndices(g.grid, latIdx+1, lonIdx+1, g.scanningMode),
	}

	// 获取相邻点的值
	points := make([]float64, 4)
	for i, idx := range indices {
		value, err := g.reader.ReadValueAt(timeStep, idx)
		if err != nil {
			return 0, err
		}
		points[i] = value
	}

	// 计算权重
	lats := g.grid.Latitudes()
	lons := g.grid.Longitudes()
	weights := []float64{
		(lat - lats[latIdx]) / (lats[latIdx+1] - lats[latIdx]),
		(lon - lons[lonIdx]) / (lons[lonIdx+1] - lons[lonIdx]),
	}

	// 使用选定的插值算法进行计算
	return g.interpolator.Interpolate(points, weights), nil
}
