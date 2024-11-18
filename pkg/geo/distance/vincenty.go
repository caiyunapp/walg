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
	return VincentyIterations(lat1, lon1, lat2, lon2, 100)
}

func VincentyIterations(lat1, lon1, lat2, lon2 float64, iterations int) float64 {
	// Convert degrees to radians
	φ1 := lat1 * math.Pi / 180.0
	λ1 := lon1 * math.Pi / 180.0
	φ2 := lat2 * math.Pi / 180.0
	λ2 := lon2 * math.Pi / 180.0

	// Calculate L (difference in longitude)
	L := λ2 - λ1

	// Calculate U1 and U2 (reduced latitude)
	tanU1 := (1 - f) * math.Tan(φ1)
	tanU2 := (1 - f) * math.Tan(φ2)
	U1 := math.Atan(tanU1)
	U2 := math.Atan(tanU2)

	// Initial value for lambda
	λ := L
	sinλ, cosλ := math.Sincos(λ)

	// Initialize variables for iteration
	var sinσ, cosσ, σ, sinα, cos2α, cos2σm float64
	var λʹ float64 // lambda prime
	const maxIterations = 100
	const epsilon = 1e-12

	// Iterate until convergence
	for i := 0; i < iterations; i++ {
		sinU1, cosU1 := math.Sincos(U1)
		sinU2, cosU2 := math.Sincos(U2)

		// Calculate sigma
		sinσ = math.Sqrt(math.Pow(cosU2*sinλ, 2) +
			math.Pow(cosU1*sinU2-sinU1*cosU2*cosλ, 2))
		if sinσ == 0 {
			return 0 // Points are coincident
		}

		cosσ = sinU1*sinU2 + cosU1*cosU2*cosλ
		σ = math.Atan2(sinσ, cosσ)

		// Calculate alpha (azimuth)
		sinα = cosU1 * cosU2 * sinλ / sinσ
		cos2α = 1 - sinα*sinα

		// Calculate cos(2σm)
		cos2σm = cosσ - 2*sinU1*sinU2/cos2α
		if math.IsNaN(cos2σm) {
			cos2σm = 0 // Equatorial line
		}

		C := f / 16 * cos2α * (4 + f*(4-3*cos2α))
		λʹ = λ
		λ = L + (1-C)*f*sinα*(σ+C*sinσ*(cos2σm+C*cosσ*(-1+2*math.Pow(cos2σm, 2))))

		// Check for convergence
		if math.Abs(λ-λʹ) < epsilon {
			// Calculate final distance
			u2 := cos2α * (a*a - b*b) / (b * b)
			A := 1 + u2/16384*(4096+u2*(-768+u2*(320-175*u2)))
			B := u2 / 1024 * (256 + u2*(-128+u2*(74-47*u2)))
			Δσ := B * sinσ * (cos2σm + B/4*(cosσ*(-1+2*math.Pow(cos2σm, 2))-
				B/6*cos2σm*(-3+4*math.Pow(sinσ, 2))*(-3+4*math.Pow(cos2σm, 2))))

			// Return distance in kilometers
			return b * A * (σ - Δσ) / 1000
		}
	}

	return -1 // Failed to converge
}
