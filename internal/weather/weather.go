package weather

func FToC(f float64) float64 {
	return f - 32.0*(5.0/9.0)
}

func CToF(c float64) float64 {
	return (9.0/5.0)*c + 32.0
}
