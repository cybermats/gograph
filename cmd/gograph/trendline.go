package main

type trendline struct {
	X []int
	Y []float32
}

type trendlineMaker struct {
	sumX  float32
	sumY  float32
	sumXX float32
	sumXY float32
	n     int
}

func (tm *trendlineMaker) addSample(x, y float32) {
	tm.sumX += x
	tm.sumY += y
	tm.sumXX += x * x
	tm.sumXY += x * y
	tm.n++
}

func (tm *trendlineMaker) trendline(offset int) trendline {
	if tm.n == 0 || tm.sumX == 0 {
		return trendline{[]int{0, 0}, []float32{0, 0}}
	}
	alpha := (float32(tm.n)*tm.sumXY - tm.sumX*tm.sumY) /
		(float32(tm.n)*tm.sumXX - tm.sumX*tm.sumX)
	beta := (tm.sumY - alpha*tm.sumX) / float32(tm.n)
	xaxis := []int{offset, offset + tm.n - 1}
	yaxis := []float32{alpha*float32(offset) + beta, alpha*float32(offset+tm.n-1) + beta}
	return trendline{xaxis, yaxis}
}
