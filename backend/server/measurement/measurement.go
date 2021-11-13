package measurement

import (
	"time"

	"github.com/rs/zerolog/log"
)

type Measurement struct {
	Position  Position  `json:"position"`
	Timestamp time.Time `json:"timestamp"`
	O2        int       `json:"o2"`
}

type Position struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

var ExampleMeasurements = []Measurement{
	{
		Position: Position{
			Lat:  21.37,
			Long: 42.69,
		},
		Timestamp: mustTimeParse(time.RFC3339, "2021-10-10T15:00:00Z"),
		O2:        80,
	},
	{
		Position: Position{
			Lat:  21.38,
			Long: 42.67,
		},
		Timestamp: mustTimeParse(time.RFC3339, "2021-10-10T15:00:02Z"),
		O2:        70,
	},
	{
		Position: Position{
			Lat:  21.37,
			Long: 42.70,
		},
		Timestamp: mustTimeParse(time.RFC3339, "2021-10-10T15:00:01Z"),
		O2:        82,
	},
	{
		Position: Position{
			Lat:  21.39,
			Long: 42.69,
		},
		Timestamp: mustTimeParse(time.RFC3339, "2021-10-10T14:00:00Z"),
		O2:        81,
	},
}

func mustTimeParse(format, timeStr string) time.Time {
	t, err := time.Parse(format, timeStr)
	if err != nil {
		log.Fatal().Err(err).Msg("parsing time")
	}
	return t
}
