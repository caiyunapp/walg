package interpolators

// Interpolator 定义插值算法接口
type Interpolator interface {
	// Interpolate 执行插值计算
	// points: 相邻点的值
	// weights: 插值权重
	Interpolate(points []float64, weights []float64) float64
}
