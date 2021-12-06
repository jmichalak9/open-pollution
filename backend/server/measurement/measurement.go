package measurement

import (
	"time"

	"github.com/rs/zerolog/log"
)

type Measurement struct {
	Position    Position  `json:"position"`
	Timestamp   time.Time `json:"timestamp"`
	O3          int       `json:"levelO3,omitempty"`
	PM10        int       `json:"levelPM10,omitempty"`
	PM25        int       `json:"levelPM25,omitempty"`
	Temperature int       `json:"temperature,omitempty"`
	SO2         int       `json:"levelSO2,omitempty"`
}

type Position struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

var ExampleMeasurements = []Measurement{
	{
		Position: Position{
			Lat:  52.2,
			Long: 21.0,
		},
		Timestamp: mustTimeParse(time.RFC3339, "2021-10-10T15:00:00Z"),
		O3:        80,
	},
	{
		Position: Position{
			Lat:  52.3,
			Long: 21.0,
		},
		Timestamp: mustTimeParse(time.RFC3339, "2021-10-10T15:00:02Z"),
		O3:        70,
	},
	{
		Position: Position{
			Lat:  52.25,
			Long: 21.0,
		},
		Timestamp: mustTimeParse(time.RFC3339, "2021-10-10T15:00:01Z"),
		O3:        82,
	},
	{
		Position: Position{
			Lat:  52.2,
			Long: 21.2,
		},
		Timestamp: mustTimeParse(time.RFC3339, "2021-10-10T14:00:00Z"),
		O3:        81,
	},
}

func mustTimeParse(format, timeStr string) time.Time {
	t, err := time.Parse(format, timeStr)
	if err != nil {
		log.Fatal().Err(err).Msg("parsing time")
	}
	return t
}
