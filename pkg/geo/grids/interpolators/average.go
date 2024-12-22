package interpolators

// AverageInterpolator 简单平均插值实现
// 简单平均插值是最基础的插值方法，直接计算所有相邻点的算术平均值
// 适用场景：
// 1. 需要平滑噪声数据
// 2. 相邻点的值差异不大的情况
// 3. 不考虑空间位置关系的快速估计
// 4. 数据预处理或初步分析
//
// 算法过程：
// 1. 对所有相邻点的值进行求和
// 2. 除以点的数量得到平均值
//
// 注意：这种方法不考虑点的空间位置关系，
// 仅适用于相邻点具有相似重要性的情况
type AverageInterpolator struct{}

func (a *AverageInterpolator) Interpolate(points []float64, weights []float64) float64 {
	sum := 0.0
	for _, point := range points {
		sum += point
	}
	return sum / float64(len(points))
}
