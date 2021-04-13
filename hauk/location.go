package hauk

// Location contains GPS data for hauk
type Location struct {
	Latitude     *float64
	Longitude    *float64
	Altitude     *int
	Velocity     *int
	Time         *float64
	Accuracy     *int
	AccuracyMode *int
	SID          *string
}
