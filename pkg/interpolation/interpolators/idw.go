package interpolators

import "math"

// IDWInterpolator 反距离加权插值实现
// IDW(Inverse Distance Weighting)是一种基于距离权重的插值方法
// 距离越近的点对插值结果的影响越大
// 适用场景：
// 1. 不规则分布的数据点插值
// 2. 地理数据分析（如降水、温度等）
// 3. 空间数据的局部特征保持
//
// 算法过程：
// 1. 计算目标点到各已知点的距离
// 2. 计算每个点的权重（距离的倒数的 power 次方）
// 3. 对所有点进行加权平均
//
// Power参数的影响：
// - Power值越大，距离较近的点的影响越大
// - Power=2 是最常用的值
// - Power=1 时为线性权重
// - Power>2 时会加强局部特征
type IDWInterpolator struct {
	Power float64
}

func NewIDWInterpolator(power float64) *IDWInterpolator {
	return &IDWInterpolator{Power: power}
}

func (i *IDWInterpolator) Interpolate(points []float64, weights []float64) float64 {
	var weightSum, valueSum float64

	for idx, point := range points {
		// 计算到四个角点的距离
		var dist float64
		switch idx {
		case 0: // 左下角点距离
			dist = math.Sqrt(weights[0]*weights[0] + weights[1]*weights[1])
		case 1: // 右下角点距离
			dist = math.Sqrt(weights[0]*weights[0] + (1-weights[1])*(1-weights[1]))
		case 2: // 左上角点距离
			dist = math.Sqrt((1-weights[0])*(1-weights[0]) + weights[1]*weights[1])
		case 3: // 右上角点距离
			dist = math.Sqrt((1-weights[0])*(1-weights[0]) + (1-weights[1])*(1-weights[1]))
		}

		if dist < 1e-10 {
			return point // 如果距离极小，直接返回该点的值
		}

		// 计算权重并累加
		weight := 1.0 / math.Pow(dist, i.Power)
		weightSum += weight
		valueSum += weight * point
	}

	return valueSum / weightSum
}
