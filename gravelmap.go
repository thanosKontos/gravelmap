package gravelmap

type Elevation interface {
	Find(lat, lng float64) (int64, error)
}