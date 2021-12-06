package measurement_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"

	"github.com/jmichalak9/open-pollution/server/measurement"
)

var _ = Describe("InMemoryCache", func() {
	var (
		measurement1    measurement.Measurement
		measurement1a   measurement.Measurement
		measurement1b   measurement.Measurement
		measurement2    measurement.Measurement
		allMeasurements []measurement.Measurement
	)
	BeforeEach(func() {
		measurement1 = measurement.Measurement{
			Position: measurement.Position{
				Lat:  10,
				Long: 10,
			},
			Timestamp: mustTimeParse(time.RFC3339, "2021-10-10T15:00:00Z"),
		}
		measurement1a = measurement.Measurement{
			Position: measurement.Position{
				Lat:  10,
				Long: 10,
			},
			Timestamp: mustTimeParse(time.RFC3339, "2021-10-10T16:00:00Z"),
		}
		measurement1b = measurement.Measurement{
			Position: measurement.Position{
				Lat:  10,
				Long: 10,
			},
			Timestamp: mustTimeParse(time.RFC3339, "2021-10-10T10:00:00Z"),
		}
		measurement2 = measurement.Measurement{
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
		cache.UpdateMeasurements([]measurement.Measurement{measurement2})

		Expect(cache.GetMeasurements()).To(Equal(allMeasurements))
	})
	It("overrides old measurements", func() {
		cache := measurement.NewInMemoryCache([]measurement.Measurement{measurement1})
		cache.UpdateMeasurements([]measurement.Measurement{measurement1a})

		Expect(cache.GetMeasurements()).To(Equal([]measurement.Measurement{measurement1a}))
	})
	It("does not override newer measurement", func() {
		cache := measurement.NewInMemoryCache([]measurement.Measurement{measurement1})
		cache.UpdateMeasurements([]measurement.Measurement{measurement1b})

		Expect(cache.GetMeasurements()).To(Equal([]measurement.Measurement{measurement1}))
	})
})

func mustTimeParse(format, timeStr string) time.Time {
	t, err := time.Parse(format, timeStr)
	if err != nil {
		log.Fatal().Err(err).Msg("parsing time")
	}
	return t
}
