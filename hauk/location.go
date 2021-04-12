package hauk

// Location contains GPS data for hauk
type Location struct {
	Latitude     float64
	Longitude    float64
	Time         float64
	AccuracyMode int
	SID          string
}
