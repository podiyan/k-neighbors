package persistence

import "github.com/ericlagergren/go-kml"
import (
	"github.com/hongshibao/go-kdtree"
	"log"
	"github.com/dropbox/godropbox/errors"
)

// MemoryStore is a simple in-memory repository for Proxi service
// It implements Repository interface
type MemoryStore struct {
	store []kdtree.Point
	tree *kdtree.KDTree
	capacity int
	reindexPoints bool
}

func (ms *MemoryStore) AddPlaceMark(p kml.Placemark) error {

	log.Printf("Index size: %v, capacity: %v", len(ms.store), ms.capacity)

	if ms.capacity > 0 && len(ms.store) >= ms.capacity  {
		return errors.New("Repo full")
	}

	gp := NewGeoPoint(p)
	if gp != nil {
		log.Printf("Geopoint %#v", gp)

		ms.store = append(ms.store, *NewGeoPoint(p))
		ms.reindexPoints = true
	}

	return nil
}

func (ms *MemoryStore) FindKNearestPlaceMarks(lat float64, long float64, k int) []kml.Placemark {

	if len(ms.store) <= 0 {
		return []kml.Placemark{}
	}

	if ms.reindexPoints {
		// sic - need to do it due to limitations of the kd tree library used
		ms.tree = kdtree.NewKDTree(ms.store)
		ms.reindexPoints = false
	}

	nearestK := ms.tree.KNN(*GeoPointFromCoordinates(lat, long), k)
	ret := make([]kml.Placemark, len(nearestK))

	for i,point := range nearestK {
		g,_ := point.(GeoPoint)
		ret[i] = g.Data
	}

	return ret
}

func (ms *MemoryStore) GetIndexSize() int {
	return len(ms.store)
}

func NewMemoryStore(pc PersistenceConfig) *MemoryStore{
	ms := &MemoryStore{reindexPoints:true, capacity:pc.Size}
	ms.store = make([]kdtree.Point,0)
	return ms
}