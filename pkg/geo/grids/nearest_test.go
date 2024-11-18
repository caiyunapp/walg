package grids_test

import (
	"testing"

	"github.com/scorix/walg/pkg/geo/grids"
	"github.com/stretchr/testify/assert"
)

func TestNearest_FindNearestIndices(t *testing.T) {
	t.Run("ascending order", func(t *testing.T) {
		latitudes := []float64{-90, -80, -70, -60, -50, -40, -30, -20, -10, 0, 10, 20, 30, 40, 50, 60, 70, 80, 90}

		tests := []struct {
			name  string
			value float64
			want  [2]int
		}{
			{
				name:  "outside left",
				value: -91,
				want:  [2]int{0, 0},
			},
			{
				name:  "outside right",
				value: 91,
				want:  [2]int{18, 18},
			},
			{
				name:  "value in the middle",
				value: 45,
				want:  [2]int{13, 14},
			},
			{
				name:  "value at the left edge",
				value: -90,
				want:  [2]int{0, 0},
			},
			{
				name:  "value at the right edge",
				value: 90,
				want:  [2]int{18, 18},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.want, grids.FindNearestIndices(tt.value, latitudes))
			})
		}
	})

	t.Run("descending order", func(t *testing.T) {
		latitudes := []float64{90, 80, 70, 60, 50, 40, 30, 20, 10, 0, -10, -20, -30, -40, -50, -60, -70, -80, -90}

		tests := []struct {
			name  string
			value float64
			want  [2]int
		}{
			{
				name:  "outside left",
				value: 91,
				want:  [2]int{0, 0},
			},
			{
				name:  "outside right",
				value: -91,
				want:  [2]int{18, 18},
			},
			{
				name:  "value in the middle",
				value: 45,
				want:  [2]int{4, 5},
			},
			{
				name:  "value at the left edge",
				value: 90,
				want:  [2]int{0, 0},
			},
			{
				name:  "value at the right edge",
				value: -90,
				want:  [2]int{18, 18},
			},
			{
				name:  "value between [0, 1]",
				value: 85,
				want:  [2]int{0, 1},
			},
			{
				name:  "value between [1, 2]",
				value: 75,
				want:  [2]int{1, 2},
			},
			{
				name:  "value between [17, 18]",
				value: -85,
				want:  [2]int{17, 18},
			},
			{
				name:  "positive value between negative and positive",
				value: 5,
				want:  [2]int{8, 9},
			},
			{
				name:  "negative value between negative and positive",
				value: -5,
				want:  [2]int{9, 10},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.want, grids.FindNearestIndices(tt.value, latitudes))
			})
		}
	})
}
