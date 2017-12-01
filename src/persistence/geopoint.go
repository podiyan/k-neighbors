package persistence

import (
	"github.com/hongshibao/go-kdtree"
	"github.com/ericlagergren/go-kml"
	"strconv"
)

// GeoPoint
// A simple Geo coordinate to be used for the kdTree library
// This operates in Cartesian mode, will not be accurate for Polar coordinates
type GeoPoint struct {
	kdtree.Point
	Vec []float64
	Data kml.Placemark
}

func (p GeoPoint) Dim() int {
	return len(p.Vec)
}

func (p GeoPoint) GetValue(dim int) float64 {
	return p.Vec[dim]
}

func (p GeoPoint) Distance(other kdtree.Point) float64 {
	var ret float64
	for i := 0; i < p.Dim(); i++ {
		tmp := p.GetValue(i) - other.GetValue(i)
		ret += tmp * tmp
	}
	return ret
}

func (p GeoPoint) PlaneDistance(val float64, dim int) float64 {
	tmp := p.GetValue(dim) - val
	return tmp * tmp
}

func NewGeoPoint(placemark kml.Placemark) *GeoPoint {
	var lat, lon float64

	latFound := false
	lonFound := false

	for _,data := range placemark.ExtendedData.SchemaData.SimpleData {
		if latFound && lonFound {
			break
		}

		if !latFound && data.Name == "LAT" {
			lat, _ = strconv.ParseFloat(data.Value, 64)
			latFound = true
			continue
		}

		if !lonFound && data.Name == "LON" {
			lon, _ = strconv.ParseFloat(data.Value, 64)
			lonFound = true
			continue
		}
	}

	if latFound && lonFound {
		g :=  &GeoPoint{Data:placemark}
		g.Vec = append(g.Vec, lat)
		g.Vec = append(g.Vec, lon)

		return g
	}

	return nil
}

func GeoPointFromCoordinates(lat, lon float64) *GeoPoint {
	g := &GeoPoint{}

	g.Vec = append(g.Vec, lat)
	g.Vec = append(g.Vec, lon)

	return g
}