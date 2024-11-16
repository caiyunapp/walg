package latlon

import (
	"math"

	"github.com/scorix/walg/pkg/geo/grids"
)

// latLon 表示一个基于经纬度的等间距网格系统
type latLon struct {
	minLat   float64        // 最小纬度
	maxLat   float64        // 最大纬度
	minLon   float64        // 最小经度
	maxLon   float64        // 最大经度
	latStep  float64        // 纬度步长
	lonStep  float64        // 经度步长
	latCount int            // 纬度方向的网格数量
	lonCount int            // 经度方向的网格数量
	lats     []float64      // 缓存的纬度值
	lons     []float64      // 缓存的经度值
	scanMode grids.ScanMode // 新增: 扫描模式
}

type latlonOption func(*latLon)

func WithScanMode(scanMode grids.ScanMode) latlonOption {
	return func(ll *latLon) {
		ll.scanMode = scanMode
	}
}

// NewLatLonGrid 创建一个新的经纬度网格
func NewLatLonGrid(minLat, maxLat, minLon, maxLon, latStep, lonStep float64, opts ...latlonOption) *latLon {
	latCount := int((maxLat-minLat)/latStep) + 1
	lonCount := int((maxLon-minLon)/lonStep) + 1

	// 预计算所有的纬度和经度值
	lats := make([]float64, latCount)
	for i := 0; i < latCount; i++ {
		lats[i] = minLat + float64(i)*latStep
	}

	lons := make([]float64, lonCount)
	for i := 0; i < lonCount; i++ {
		lons[i] = minLon + float64(i)*lonStep
	}

	ll := &latLon{
		minLat:   minLat,
		maxLat:   maxLat,
		minLon:   minLon,
		maxLon:   maxLon,
		latStep:  latStep,
		lonStep:  lonStep,
		latCount: latCount,
		lonCount: lonCount,
		lats:     lats,
		lons:     lons,
	}

	for _, o := range opts {
		o(ll)
	}

	return ll
}

// Size 返回网格总数
func (g *latLon) Size() int {
	return g.latCount * g.lonCount
}

// Latitudes 返回所有纬度值
func (g *latLon) Latitudes() []float64 {
	return g.lats
}

// Longitudes 返回所有经度值
func (g *latLon) Longitudes() []float64 {
	return g.lons
}

// GridIndex 根据经纬度返回网格索引
func (g *latLon) GridIndex(lat, lon float64) int {
	latIdx := g.GetLatitudeIndex(lat)
	lonIdx := g.GetLongitudeIndex(lat, lon)

	return g.GridIndexFromIndices(latIdx, lonIdx)
}

// GridIndexFromIndices 根据纬度和经度的索引计算网格索引
func (g *latLon) GridIndexFromIndices(latIdx, lonIdx int) int {
	if latIdx < 0 || latIdx >= g.latCount || lonIdx < 0 || lonIdx >= g.lonCount {
		return -1
	}

	// 处理负方向扫描
	if g.scanMode.IsNegativeI() {
		lonIdx = g.lonCount - 1 - lonIdx
	}

	// 处理正向 J 扫描（默认是负向）
	if g.scanMode.IsPositiveJ() {
		latIdx = g.latCount - 1 - latIdx
	}

	// 处理连续 J 方向
	if g.scanMode.IsConsecutiveJ() {
		return lonIdx*g.latCount + latIdx
	}

	// 处理交替行
	if g.scanMode.HasOppositeRows() && latIdx%2 == 1 {
		lonIdx = g.lonCount - 1 - lonIdx
	}

	// 默认是连续 I 方向
	return latIdx*g.lonCount + lonIdx
}

// GridPoint 根据网格索引返回对应的经纬度
func (g *latLon) GridPoint(index int) (lat, lon float64) {
	if index < 0 || index >= g.Size() {
		return math.NaN(), math.NaN()
	}

	var latIdx, lonIdx int
	if g.scanMode.IsConsecutiveJ() {
		lonIdx = index / g.latCount
		latIdx = index % g.latCount
	} else {
		latIdx = index / g.lonCount
		lonIdx = index % g.lonCount
	}

	// 处理负方向扫描
	if g.scanMode.IsNegativeI() {
		lonIdx = g.lonCount - 1 - lonIdx
	}

	// 处理正向 J 扫描
	if g.scanMode.IsPositiveJ() {
		latIdx = g.latCount - 1 - latIdx
	}

	// 处理交替行
	if g.scanMode.HasOppositeRows() && latIdx%2 == 1 {
		lonIdx = g.lonCount - 1 - lonIdx
	}

	return g.lats[latIdx], g.lons[lonIdx]
}

// normalize 将一个值映射到指定区间 [minVal, maxVal]
func (g *latLon) normalize(minVal, maxVal, val float64) float64 {
	// Handle values outside the range
	interval := maxVal - minVal
	normalized := val

	// Wrap around until we're in range
	for normalized < minVal {
		normalized += interval
	}
	for normalized > maxVal {
		normalized -= interval
	}

	return normalized
}

// GetLatitudeIndex returns the index of the latitude band containing the given latitude
func (g *latLon) GetLatitudeIndex(lat float64) int {
	normalized := g.normalize(g.minLat, g.maxLat, lat)
	relativePos := normalized - g.minLat
	idx := int(math.Floor(relativePos / g.latStep))
	if idx == g.latCount { // Handle edge case for maximum value
		idx = idx - 1
	}

	return idx
}

// GetLongitudeIndex returns the index of the longitude band containing the given longitude
func (g *latLon) GetLongitudeIndex(lat, lon float64) int {
	normalized := g.normalize(g.minLon, g.maxLon, lon)
	relativePos := normalized - g.minLon
	idx := int(math.Floor(relativePos / g.lonStep))
	if idx == g.lonCount { // Handle edge case for maximum value
		idx = idx - 1
	}

	return idx
}
