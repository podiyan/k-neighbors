package persistence

import "github.com/ericlagergren/go-kml"

// PersistenceConfig Allows configuration of service backends
type PersistenceConfig struct {
	Type string  `json:"type"`
	Size int	 `json:"size"`
}


// Repository interface abstracts the backend for proxi service
type Repository interface {
	AddPlaceMark(kml.Placemark) error
	FindKNearestPlaceMarks(float64, float64, int) []kml.Placemark
	GetIndexSize() int
}