package gaussian

import (
	"math"
)

// legendrePolynomial calculates the value of nth order Legendre polynomial at point x
func legendrePolynomial(n int, x float64) float64 {
	if n == 0 {
		return 1.0
	}
	if n == 1 {
		return x
	}

	// Using recurrence relation: (n+1)Pn+1(x) = (2n+1)xPn(x) - nPn-1(x)
	p0, p1 := 1.0, x
	var pn float64

	for i := 2; i <= n; i++ {
		pn = ((2*float64(i)-1)*x*p1 - (float64(i)-1)*p0) / float64(i)
		p0, p1 = p1, pn
	}
	return p1
}

// legendrePolynomialDerivative calculates the derivative value of nth order Legendre polynomial at point x
func legendrePolynomialDerivative(n int, x float64) float64 {
	if n == 0 {
		return 0.0
	}
	if n == 1 {
		return 1.0
	}
	return (float64(n) * (x*legendrePolynomial(n, x) - legendrePolynomial(n-1, x))) / (x*x - 1.0)
}

// gaussLegendreZeros calculates zeros of nth order Legendre polynomial
// Returns latitudes (in degrees) sorted in ascending order (from South Pole to North Pole)
func gaussLegendreZeros(n int) []float64 {
	zeros := make([]float64, n)

	// For each zero, use Newton's iteration method to solve
	for i := 0; i < n; i++ {
		// Initial estimate
		x := math.Cos(math.Pi * (float64(4*i+3) / float64(4*n+2)))

		// Newton iteration
		for iter := 0; iter < 10; iter++ { // Usually 10 iterations are sufficient
			dx := -legendrePolynomial(n, x) / legendrePolynomialDerivative(n, x)
			x += dx
			if math.Abs(dx) < 1e-15 {
				break
			}
		}

		// Convert zeros in [-1,1] interval to latitude values [-90°,90°]
		zeros[i] = math.Asin(x) * 180.0 / math.Pi
	}

	return zeros
}
