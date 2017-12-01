package main

import (
	//"encoding/json"
	"github.com/ericlagergren/go-kml"
	"github.com/paulmach/go.geojson"
	"net/http"
	"io"
)

func PlacemarkToGeoJSON(placemark kml.Placemark) *geojson.FeatureCollection {
	var fc *geojson.FeatureCollection = geojson.NewFeatureCollection()

	return fc
}


func reportBadRequest(w http.ResponseWriter, s string) {
	w.WriteHeader(http.StatusBadRequest)
	io.WriteString(w, s)
}

func resultsToKML(p []kml.Placemark) *kml.Document {
	d := &kml.Document{}
	d.Folder = &kml.Folder{}
	d.Folder.Name = "Results"
	d.Folder.Placemark = make([]kml.Placemark,len(p))

	for i, placemark := range p {
		d.Folder.Placemark[i] = placemark
	}

	return d
}