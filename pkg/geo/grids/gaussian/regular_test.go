package gaussian_test

import (
	"testing"

	"github.com/scorix/walg/pkg/geo/grids"
	"github.com/scorix/walg/pkg/geo/grids/gaussian"
	"github.com/stretchr/testify/assert"
)

func TestRegular_F768(t *testing.T) {
	g := gaussian.NewRegular(768)

	assert.Equal(t, 1536, len(g.Latitudes()))
	assert.Equal(t, 3072, len(g.Longitudes()))
}

// https://confluence.ecmwf.int/display/UDOC/N48
func TestRegular_F48(t *testing.T) {
	g := gaussian.NewRegular(48)

	lats := g.Latitudes()
	lons := g.Longitudes()
	t.Logf("lats: %v", lats)
	t.Logf("lons: %v", lons)

	assert.Equal(t, 96, len(lats))
	assert.Equal(t, 192, len(lons))

	assert.InDelta(t, 88.572169, lats[0], 1e-6)
	assert.InDelta(t, 86.722531, lats[1], 1e-6)
	assert.InDelta(t, 84.861970, lats[2], 1e-6)
	assert.InDelta(t, 82.998942, lats[3], 1e-6)
	assert.InDelta(t, 81.134977, lats[4], 1e-6)
	assert.InDelta(t, 79.270559, lats[5], 1e-6)
	assert.InDelta(t, 77.405888, lats[6], 1e-6)
	assert.InDelta(t, 75.541061, lats[7], 1e-6)
	assert.InDelta(t, 73.676132, lats[8], 1e-6)
	assert.InDelta(t, 71.811132, lats[9], 1e-6)
	assert.InDelta(t, 69.946081, lats[10], 1e-6)
	assert.InDelta(t, 68.080991, lats[11], 1e-6)
	assert.InDelta(t, 66.215872, lats[12], 1e-6)
	assert.InDelta(t, 64.350730, lats[13], 1e-6)
	assert.InDelta(t, 62.485571, lats[14], 1e-6)
	assert.InDelta(t, 60.620396, lats[15], 1e-6)
	assert.InDelta(t, 58.755209, lats[16], 1e-6)
	assert.InDelta(t, 56.890013, lats[17], 1e-6)
	assert.InDelta(t, 55.024808, lats[18], 1e-6)
	assert.InDelta(t, 53.159595, lats[19], 1e-6)
	assert.InDelta(t, 51.294377, lats[20], 1e-6)
	assert.InDelta(t, 49.429154, lats[21], 1e-6)
	assert.InDelta(t, 47.563926, lats[22], 1e-6)
	assert.InDelta(t, 45.698694, lats[23], 1e-6)
	assert.InDelta(t, 43.833459, lats[24], 1e-6)
	assert.InDelta(t, 41.968220, lats[25], 1e-6)
	assert.InDelta(t, 40.102979, lats[26], 1e-6)
	assert.InDelta(t, 38.237736, lats[27], 1e-6)
	assert.InDelta(t, 36.372491, lats[28], 1e-6)
	assert.InDelta(t, 34.507243, lats[29], 1e-6)
	assert.InDelta(t, 32.641994, lats[30], 1e-6)
	assert.InDelta(t, 30.776744, lats[31], 1e-6)
	assert.InDelta(t, 28.911492, lats[32], 1e-6)
	assert.InDelta(t, 27.046239, lats[33], 1e-6)
	assert.InDelta(t, 25.180986, lats[34], 1e-6)
	assert.InDelta(t, 23.315731, lats[35], 1e-6)
	assert.InDelta(t, 21.450475, lats[36], 1e-6)
	assert.InDelta(t, 19.585219, lats[37], 1e-6)
	assert.InDelta(t, 17.719962, lats[38], 1e-6)
	assert.InDelta(t, 15.854704, lats[39], 1e-6)
	assert.InDelta(t, 13.989446, lats[40], 1e-6)
	assert.InDelta(t, 12.124187, lats[41], 1e-6)
	assert.InDelta(t, 10.258928, lats[42], 1e-6)
	assert.InDelta(t, 8.393669, lats[43], 1e-6)
	assert.InDelta(t, 6.528409, lats[44], 1e-6)
	assert.InDelta(t, 4.663150, lats[45], 1e-6)
	assert.InDelta(t, 2.797890, lats[46], 1e-6)
	assert.InDelta(t, 0.932630, lats[47], 1e-6)
	assert.InDelta(t, -0.932630, lats[48], 1e-6)
	assert.InDelta(t, -2.797890, lats[49], 1e-6)
	assert.InDelta(t, -4.663150, lats[50], 1e-6)
	assert.InDelta(t, -6.528409, lats[51], 1e-6)
	assert.InDelta(t, -8.393669, lats[52], 1e-6)
	assert.InDelta(t, -10.258928, lats[53], 1e-6)
	assert.InDelta(t, -12.124187, lats[54], 1e-6)
	assert.InDelta(t, -13.989446, lats[55], 1e-6)
	assert.InDelta(t, -15.854704, lats[56], 1e-6)
	assert.InDelta(t, -17.719962, lats[57], 1e-6)
	assert.InDelta(t, -19.585219, lats[58], 1e-6)
	assert.InDelta(t, -21.450475, lats[59], 1e-6)
	assert.InDelta(t, -23.315731, lats[60], 1e-6)
	assert.InDelta(t, -25.180986, lats[61], 1e-6)
	assert.InDelta(t, -27.046239, lats[62], 1e-6)
	assert.InDelta(t, -28.911492, lats[63], 1e-6)
	assert.InDelta(t, -30.776744, lats[64], 1e-6)
	assert.InDelta(t, -32.641994, lats[65], 1e-6)
	assert.InDelta(t, -34.507243, lats[66], 1e-6)
	assert.InDelta(t, -36.372491, lats[67], 1e-6)
	assert.InDelta(t, -38.237736, lats[68], 1e-6)
	assert.InDelta(t, -40.102979, lats[69], 1e-6)
	assert.InDelta(t, -41.968220, lats[70], 1e-6)
	assert.InDelta(t, -43.833459, lats[71], 1e-6)
	assert.InDelta(t, -45.698694, lats[72], 1e-6)
	assert.InDelta(t, -47.563926, lats[73], 1e-6)
	assert.InDelta(t, -49.429154, lats[74], 1e-6)
	assert.InDelta(t, -51.294377, lats[75], 1e-6)
	assert.InDelta(t, -53.159595, lats[76], 1e-6)
	assert.InDelta(t, -55.024808, lats[77], 1e-6)
	assert.InDelta(t, -56.890013, lats[78], 1e-6)
	assert.InDelta(t, -58.755209, lats[79], 1e-6)
	assert.InDelta(t, -60.620396, lats[80], 1e-6)
	assert.InDelta(t, -62.485571, lats[81], 1e-6)
	assert.InDelta(t, -64.350730, lats[82], 1e-6)
	assert.InDelta(t, -66.215872, lats[83], 1e-6)
	assert.InDelta(t, -68.080991, lats[84], 1e-6)
	assert.InDelta(t, -69.946081, lats[85], 1e-6)
	assert.InDelta(t, -71.811132, lats[86], 1e-6)
	assert.InDelta(t, -73.676132, lats[87], 1e-6)
	assert.InDelta(t, -75.541061, lats[88], 1e-6)
	assert.InDelta(t, -77.405888, lats[89], 1e-6)
	assert.InDelta(t, -79.270559, lats[90], 1e-6)
	assert.InDelta(t, -81.134977, lats[91], 1e-6)
	assert.InDelta(t, -82.998942, lats[92], 1e-6)
	assert.InDelta(t, -84.861970, lats[93], 1e-6)
	assert.InDelta(t, -86.722531, lats[94], 1e-6)
	assert.InDelta(t, -88.572169, lats[95], 1e-6)

	assert.InDelta(t, 0.0, lons[0], 1e-6)
	assert.InDelta(t, 1.875, lons[1], 1e-6)
	assert.InDelta(t, 358.125, lons[191], 1e-6)
}

func TestRegular_GridIndexFromIndices(t *testing.T) {
	tests := []struct {
		name   string
		n      int
		mode   grids.ScanMode
		latIdx int
		lonIdx int
		want   int
	}{
		{
			name:   "consecutive i (default)",
			n:      48,
			mode:   grids.ScanModePositiveI,
			latIdx: 1,
			lonIdx: 2,
			want:   1*192 + 2, // latIdx * longitudesSize + lonIdx
		},
		{
			name:   "consecutive j",
			n:      48,
			mode:   grids.ScanModeConsecutiveJ,
			latIdx: 1,
			lonIdx: 2,
			want:   2*96 + 1, // lonIdx * latitudesSize + latIdx
		},
		{
			name:   "consecutive i with opposite rows - even row",
			n:      48,
			mode:   grids.ScanModePositiveI | grids.ScanModeOppositeRows,
			latIdx: 0,
			lonIdx: 2,
			want:   0*192 + 2, // even row - normal order
		},
		{
			name:   "consecutive i with opposite rows - odd row",
			n:      48,
			mode:   grids.ScanModePositiveI | grids.ScanModeOppositeRows,
			latIdx: 1,
			lonIdx: 2,
			want:   1*192 + (191 - 2), // odd row - reverse longitude order
		},
		{
			name:   "consecutive j with opposite rows - even column",
			n:      48,
			mode:   grids.ScanModeConsecutiveJ | grids.ScanModeOppositeRows,
			latIdx: 1,
			lonIdx: 0,
			want:   0*96 + 1, // even column - normal order
		},
		{
			name:   "consecutive j with opposite rows - odd column",
			n:      48,
			mode:   grids.ScanModeConsecutiveJ | grids.ScanModeOppositeRows,
			latIdx: 1,
			lonIdx: 1,
			want:   1*96 + (95 - 1), // odd column - reverse latitude order
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gaussian.NewRegular(tt.n, gaussian.WithScanMode(tt.mode))
			got := grids.GridIndexFromIndices(g, tt.latIdx, tt.lonIdx)
			assert.Equal(t, tt.want, got)
		})
	}
}
