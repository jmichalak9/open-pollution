package measurement_test

import (
	"github.com/jmichalak9/open-pollution/server/measurement"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InMemoryCache", func() {
	var (
		measurement1    measurement.Measurement
		measurement2    measurement.Measurement
		allMeasurements []measurement.Measurement
	)
	BeforeEach(func() {
		measurement1 = measurement.Measurement{
			Position: measurement.Position{
				Lat:  10,
				Long: 10,
			},
		}
		measurement1 = measurement.Measurement{
			Position: measurement.Position{
				Lat:  20,
				Long: 20,
			},
		}

		allMeasurements = []measurement.Measurement{measurement1, measurement2}
	})
	It("returns all cached measurements", func() {
		cache := measurement.NewInMemoryCache(allMeasurements)

		Expect(cache.GetMeasurements()).To(Equal(allMeasurements))
	})
	It("appends cached measurements", func() {
		cache := measurement.NewInMemoryCache([]measurement.Measurement{measurement1})
		cache.AppendMeasurements([]measurement.Measurement{measurement2})

		Expect(cache.GetMeasurements()).To(Equal(allMeasurements))
	})

})
