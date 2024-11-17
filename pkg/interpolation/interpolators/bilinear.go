package interpolators

// BilinearInterpolator 双线性插值实现
// 双线性插值是一种在二维空间上进行插值的方法，通过对 x 和 y 两个方向进行线性插值得到结果
// 适用场景：
// 1. 规则网格数据的插值（如气象数据、地理信息数据）
// 2. 图像处理中的图像缩放
// 3. 需要平滑过渡的数据插值
//
// 算法过程：
// 1. 在 x 方向进行两次线性插值
// 2. 在 y 方向对上述结果再进行一次线性插值
// 3. 最终得到插值点的值
//
// points 数组中四个点的位置：
// points[0]: 左下角 (x0,y0)
// points[1]: 右下角 (x1,y0)
// points[2]: 左上角 (x0,y1)
// points[3]: 右上角 (x1,y1)
//
// weights 数组含义：
// weights[0]: y方向的权重 ((y-y0)/(y1-y0))
// weights[1]: x方向的权重 ((x-x0)/(x1-x0))
type BilinearInterpolator struct{}

func (bi *BilinearInterpolator) Interpolate(points []float64, weights []float64) float64 {
	return (1-weights[0])*(1-weights[1])*points[0] +
		(1-weights[0])*weights[1]*points[1] +
		weights[0]*(1-weights[1])*points[2] +
		weights[0]*weights[1]*points[3]
}
