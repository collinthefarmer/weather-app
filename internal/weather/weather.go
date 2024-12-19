package weather

import "math"

const openMeteoPrecision = 2

// https://stackoverflow.com/a/29786394
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func FToC(f float64) float64 {
	return toFixed(f-32.0*(5.0/9.0), openMeteoPrecision)
}

func CToF(c float64) float64 {
	return toFixed((9.0/5.0)*c+32.0, openMeteoPrecision)
}
