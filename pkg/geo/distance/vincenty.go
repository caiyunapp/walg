package distance

import (
	"math"
)

// WGS-84 ellipsoid parameters
const (
	a = 6378137.0         // semi-major axis in meters
	f = 1 / 298.257223563 // flattening
	b = a * (1 - f)       // semi-minor axis
)

// Vincenty calculates the distance between two points on the Earth's surface using Vincenty's formula
// Returns distance in kilometers. Returns -1 if the algorithm fails to converge.
func Vincenty(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert degrees to radians
	φ1 := lat1 * math.Pi / 180.0
	λ1 := lon1 * math.Pi / 180.0
	φ2 := lat2 * math.Pi / 180.0
	λ2 := lon2 * math.Pi / 180.0

	// Calculate reduced latitudes
	U1 := math.Atan((1 - f) * math.Tan(φ1))
	U2 := math.Atan((1 - f) * math.Tan(φ2))

	L := λ2 - λ1
	λ := L

	sinU1 := math.Sin(U1)
	cosU1 := math.Cos(U1)
	sinU2 := math.Sin(U2)
	cosU2 := math.Cos(U2)

	// Iterate until convergence
	const maxIterations = 100
	const epsilon = 1e-12

	var (
		sinλ, cosλ, sinσ, cosσ, σ, sinα, cos2α, cos2σm float64
		C, λʹ                                          float64
	)

	for i := 0; i < maxIterations; i++ {
		sinλ = math.Sin(λ)
		cosλ = math.Cos(λ)

		// Calculate sigma
		term1 := cosU2 * sinλ
		term2 := cosU1*sinU2 - sinU1*cosU2*cosλ
		sinσ = math.Sqrt(term1*term1 + term2*term2)
		cosσ = sinU1*sinU2 + cosU1*cosU2*cosλ
		σ = math.Atan2(sinσ, cosσ)

		// Calculate alpha
		sinα = cosU1 * cosU2 * sinλ / sinσ
		cos2α = 1 - sinα*sinα
		cos2σm = cosσ - 2*sinU1*sinU2/cos2α

		C = f / 16 * cos2α * (4 + f*(4-3*cos2α))
		λʹ = λ
		λ = L + (1-C)*f*sinα*(σ+C*sinσ*(cos2σm+C*cosσ*(-1+2*cos2σm*cos2σm)))

		if math.Abs(λ-λʹ) < epsilon {
			// Calculate distance
			u2 := cos2α * (a*a - b*b) / (b * b)
			A := 1 + u2/16384*(4096+u2*(-768+u2*(320-175*u2)))
			B := u2 / 1024 * (256 + u2*(-128+u2*(74-47*u2)))
			Δσ := B * sinσ * (cos2σm + B/4*(cosσ*(-1+2*cos2σm*cos2σm)-
				B/6*cos2σm*(-3+4*sinσ*sinσ)*(-3+4*cos2σm*cos2σm)))

			// Return distance in kilometers
			return b * A * (σ - Δσ) / 1000
		}
	}

	// Algorithm didn't converge
	return -1
}
