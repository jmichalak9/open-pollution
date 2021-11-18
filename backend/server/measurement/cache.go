package measurement

// Cache is an interface for caching measurements.
type Cache interface {
	GetMeasurements() []Measurement
	AppendMeasurements([]Measurement)
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

func (c *InMemoryCache) AppendMeasurements(newMeasurements []Measurement) {
	// TODO: make this thread-safe
	c.measurements = append(c.measurements, newMeasurements...)
}
