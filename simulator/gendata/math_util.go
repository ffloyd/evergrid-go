package gendata

import "math/rand"

func uniformDistr(minValue int, maxValue int) int {
	if minValue > maxValue {
		panic("Incorrect arguments")
	}

	delta := maxValue - minValue + 1

	return minValue + rand.Intn(delta)
}

func coin(probability float64) bool {
	return rand.Float64() <= probability
}

func uniformDistrF64(minValue float64, maxValue float64) float64 {
	if minValue > maxValue {
		panic("Incorrect arguments")
	}

	delta := maxValue - minValue + 1

	return minValue + rand.Float64()*delta
}
