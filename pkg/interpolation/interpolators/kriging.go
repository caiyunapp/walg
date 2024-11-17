package interpolators

import "math"

// KrigingInterpolator 克里金插值实现
// Kriging是一种基于空间统计学的最优插值方法
// 适用场景：
// 1. 地理统计数据分析
// 2. 矿产资源评估
// 3. 气象数据空间分布
//
// 参数说明：
// - Sill(基台值): 表示空间自相关的总体变异程度
// - Range(变程): 表示空间自相关的影响范围
// - Nugget(块金值): 表示测量误差和微观尺度变异
//
// 算法过程：
// 1. 计算目标点到四个已知点的距离
// 2. 根据距离计算变异函数值：
//   - 当距离 <= Range时：γ(h) = Nugget + Sill * (1.5h/Range - 0.5(h/Range)³)
//   - 当距离 > Range时：γ(h) = Nugget + Sill
//
// 3. 使用变异函数值的倒数作为权重
// 4. 进行加权平均得到最终结果
//
// 计算示例：
// 假设目标点位于 (0.2, 0.2)，四个已知点值为 [10, 20, 30, 40]
// 1. 计算距离：
//   - 到左下角(0,0)距离: sqrt(0.2² + 0.2²) ≈ 0.283
//   - 到右下角(1,0)距离: sqrt(0.8² + 0.2²) ≈ 0.825
//   - 到左上角(0,1)距离: sqrt(0.2² + 0.8²) ≈ 0.825
//   - 到右上角(1,1)距离: sqrt(0.8² + 0.8²) ≈ 1.131
//
// 2. 计算变异函数值（假设 Range=1.0, Sill=1.0, Nugget=0.1）
// 3. 计算权重并加权平均
type KrigingInterpolator struct {
	// Sill (基台值) 表示空间自相关的总体变异程度
	// - 物理意义：表示当两点距离足够远时的最大变异值
	// - 取值范围：通常大于0，常用1.0作为标准化值
	// - 参数影响：
	//   * 值越大，表示空间变异性越强
	//   * 值越小，表示空间分布越均匀
	// - 估算方法：
	//   1. 可以使用样本数据的方差作为参考
	//   2. 或通过变异函数拟合确定
	Sill float64

	// Range (变程) 表示空间自相关的影响范围
	// - 物理意义：超过这个距离的两点基本不相关
	// - 取值范围：在标准化空间(0-1)中通常取0.3-1.0
	// - 参数影响：
	//   * 值越大，远处点的影响越大
	//   * 值越小，插值结果更依赖近邻点
	// - 估算方法：
	//   1. 通过实验变异函数确定
	//   2. 或根据实际物理现象确定影响范围
	Range float64

	// Nugget (块金值) 表示测量误差和微观尺度变异
	// - 物理意义：表示极小尺度上的随机变异
	// - 取值范围：通常是Sill的0-30%
	// - 参数影响：
	//   * 值越大，空间相关性越弱
	//   * 值越小，插值结果更光滑
	// - 估算方法：
	//   1. 通过重复测量评估测量误差
	//   2. 或通过变异函数在原点处的截距确定
	Nugget float64
}

func NewKrigingInterpolator(sill, range_, nugget float64) *KrigingInterpolator {
	return &KrigingInterpolator{
		Sill:   sill,
		Range:  range_,
		Nugget: nugget,
	}
}

func (k *KrigingInterpolator) Interpolate(points []float64, weights []float64) float64 {
	// weights[1]是x坐标，weights[0]是y坐标
	// 这里交换是为了匹配常见的(x,y)坐标系表示方式
	x, y := weights[1], weights[0]

	// 当目标点正好在四个角点位置时，直接返回对应角点的值
	// 这避免了距离为0时的除零错误，同时保证了插值的连续性
	if x == 0 && y == 0 {
		return points[0] // 左下角
	}
	if x == 1 && y == 0 {
		return points[1] // 右下角
	}
	if x == 0 && y == 1 {
		return points[2] // 左上角
	}
	if x == 1 && y == 1 {
		return points[3] // 右上角
	}

	// 使用变异函数计算每个已知点的权重并进行加权平均
	var totalWeight, weightedSum float64
	for i, point := range points {
		// 计算目标点到当前已知点的相对位置
		var dx, dy float64
		switch i {
		case 0: // 左下角点(0,0)
			dx, dy = x, y
		case 1: // 右下角点(1,0)
			dx, dy = 1-x, y
		case 2: // 左上角点(0,1)
			dx, dy = x, 1-y
		case 3: // 右上角点(1,1)
			dx, dy = 1-x, 1-y
		}

		// 计算欧几里得距离
		dist := math.Sqrt(dx*dx + dy*dy)
		// 如果距离极小，认为在已知点上，直接返回该点的值
		if dist < 1e-10 {
			return point
		}

		// 计算变异函数值
		// 变异函数描述了空间自相关性，距离越远，变异程度越大
		var gamma float64
		if dist <= k.Range {
			// 在变程范围内使用球状模型
			// h 是标准化距离
			h := dist / k.Range
			// 球状模型：γ(h) = Nugget + Sill * (1.5h - 0.5h³)
			gamma = k.Nugget + k.Sill*(1.5*h-0.5*h*h*h)
		} else {
			// 超出变程范围时，变异函数值为基台值
			gamma = k.Nugget + k.Sill
		}

		// 使用变异函数值的倒数作为权重
		// 加上 Nugget 是为了避免变异函数值为0时的除零错误
		weight := 1.0 / (gamma + k.Nugget)
		totalWeight += weight
		weightedSum += weight * point
	}

	// 返回加权平均结果
	return weightedSum / totalWeight
}
