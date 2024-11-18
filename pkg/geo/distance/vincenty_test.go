package distance_test

import (
	"testing"

	"github.com/scorix/walg/pkg/geo/distance"
	"github.com/stretchr/testify/assert"
)

func TestVincenty(t *testing.T) {
	tests := []struct {
		name string
		lat1 float64
		lon1 float64
		lat2 float64
		lon2 float64
	}{
		{name: "Identical Points", lat1: 0, lon1: 0, lat2: 0, lon2: 0},
		{name: "Poles", lat1: 90, lon1: 0, lat2: -90, lon2: 0},
		{name: "Equator", lat1: 0, lon1: 0, lat2: 0, lon2: 180},
		{name: "Prime Meridian", lat1: 0, lon1: 0, lat2: 90, lon2: 0},
		{name: "Date Line", lat1: 0, lon1: 0, lat2: 0, lon2: -180},
		{name: "Antipodes", lat1: 0, lon1: 0, lat2: -90, lon2: 180},
		{name: "Equator and Prime Meridian", lat1: 0, lon1: 0, lat2: 90, lon2: 180},
		{name: "Equator and Date Line", lat1: 0, lon1: 0, lat2: 0, lon2: -180},
		{name: "Prime Meridian and Equator", lat1: 90, lon1: 0, lat2: 0, lon2: 180},
		{name: "Prime Meridian and Date Line", lat1: 90, lon1: 0, lat2: 0, lon2: -180},
		{name: "Shanghai and Beijing", lat1: 31.2304, lon1: 121.4737, lat2: 39.9042, lon2: 116.4074},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vDist := distance.Vincenty(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			hDist := distance.Haversine(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			assert.InDelta(t, vDist, hDist, 100, "Vincenty distance is %f, haversine distance is %f", vDist, hDist)
		})
	}
}
