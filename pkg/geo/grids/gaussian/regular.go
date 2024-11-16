package gaussian

type regular struct {
	n int
}

func NewRegular(n int) *regular {
	return &regular{n: n}
}

func (g *regular) Size() int {
	return g.longitudesSize() * g.latitudesSize()
}

func (g *regular) latitudesSize() int {
	return g.n * 2
}

func (g *regular) longitudesSize() int {
	return g.n * 4
}

func (g *regular) GridIndex(lat, lon float64) int {
	return 0
}

func (g *regular) GridPoint(index int) (lat, lon float64) {
	return 0, 0
}

func (g *regular) Latitudes() []float64 {
	latitudes := gaussLegendreZeros(g.latitudesSize())

	return latitudes
}

func (g *regular) Longitudes() []float64 {
	length := g.longitudesSize()
	longitudes := make([]float64, length)

	for i := 0; i < length; i++ {
		longitudes[i] = 360.0 * float64(i) / float64(length)
	}

	return longitudes
}
