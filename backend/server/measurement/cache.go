package measurement

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
)

// if distance between positions is less than epsilon, it's considered as equal
const epsilon = 100.0

// Cache is an interface for caching measurements.
type Cache interface {
	GetMeasurements() []Measurement
	UpdateMeasurements([]Measurement)
}

// InMemoryCache implements Cache interface by storing objects in memory.
type InMemoryCache struct {
	measurements []Measurement
}

// NewInMemoryCache creates a new instance of InMemoryCache.
func NewInMemoryCache(measurements []Measurement) *InMemoryCache {
	return &InMemoryCache{
		measurements: measurements,
	}
}

func (c *InMemoryCache) GetMeasurements() []Measurement {
	return c.measurements
}

func (c *InMemoryCache) UpdateMeasurements(newMeasurements []Measurement) {
	// TODO: make this thread-safe
	for _, m := range newMeasurements {
		c.measurements = reconcile(c.measurements, m)
	}
}

func reconcile(old []Measurement, candidate Measurement) []Measurement {
	p1 := orb.Point{
		candidate.Position.Long,
		candidate.Position.Lat,
	}
	for i, o := range old {
		p2 := orb.Point{
			o.Position.Long,
			o.Position.Lat,
		}
		if geo.DistanceHaversine(p1, p2) < epsilon {
			if candidate.Timestamp.After(o.Timestamp) {
				old[i] = candidate
			}
			return old
		}
	}
	// there's no measurement near old ones, so we just add it
	return append(old, candidate)
}
