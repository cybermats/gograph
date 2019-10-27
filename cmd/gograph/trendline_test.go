package main

import "testing"

func TestTrendlineEmpty(t *testing.T) {
	maker := trendlineMaker{}
	expected := trendline{[2]int{0, 0}, [2]float32{0, 0}}
	actual := maker.trendline(0)

	if expected.X != actual.X ||
		expected.Y != actual.Y {
		t.Fail()
	}
}

func TestTrendline(t *testing.T) {
	maker := trendlineMaker{}
	expected := trendline{[2]int{0, 9}, [2]float32{8.738181818, 9.441818182}}

	ys := [10]float32{9, 8.8, 8.7, 8.8, 9.1, 9.2, 9.2, 9, 9.6, 9.5}

	for x, y := range ys {
		maker.addSample(float32(x), y)
	}

	actual := maker.trendline(0)

	if expected.X != actual.X ||
		expected.Y != actual.Y {
		t.Fail()
	}
}
