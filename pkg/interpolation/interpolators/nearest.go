package interpolators

// NearestInterpolator 最近邻插值实现
// 最近邻插值是最简单的插值方法，直接使用距离目标点最近的已知点的值
// 适用场景：
// 1. 分类数据的插值（如土地利用类型）
// 2. 需要保持原始数据值的情况
// 3. 对计算速度要求高的场景
// 4. 离散数据的插值
//
// 算法过程：
// 1. 计算目标点到四个相邻点的距离
// 2. 选择距离最近的点的值作为插值结果
//
// weights 数组表示目标点在两个方向上的相对位置：
// weights[0] < 0.5 表示更靠近下方的点
// weights[1] < 0.5 表示更靠近左边的点
type NearestInterpolator struct{}

func (n *NearestInterpolator) Interpolate(points []float64, weights []float64) float64 {
	if weights[0] < 0.5 {
		if weights[1] < 0.5 {
			return points[0] // 左下角点最近
		}
		return points[1] // 右下角点最近
	} else {
		if weights[1] < 0.5 {
			return points[2] // 左上角点最近
		}
		return points[3] // 右上角点最近
	}
}
