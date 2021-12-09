package server_test

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/jmichalak9/open-pollution/server/measurement"
)

var (
	measurements []measurement.Measurement
)

var _ = BeforeEach(func() {
	measurements = []measurement.Measurement{
		{
			Position: measurement.Position{
				Lat:  21.0,
				Long: 37.0,
			},
			Timestamp: time.Date(2005, 04, 02, 21, 37, 0, 0, time.UTC),
			O3:        42,
		},
	}
})

var _ = DescribeTable("OpenPollution API", func(t opAPITest) {
	if t.recordMocks != nil {
		t.recordMocks()
	}

	resp, err := http.Get(fmt.Sprintf("http://localhost:%v%s", port, t.endpoint))
	Expect(err).NotTo(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(t.expectedCode))
	Expect(resp.Header.Get("Content-Type")).To(Equal("application/json"))
	body, err := io.ReadAll(resp.Body)
	Expect(err).NotTo(HaveOccurred())
	Expect(body).To(MatchJSON(t.expectedBody))

},
	Entry("measurements endpoint", opAPITest{
		recordMocks: func() {
			expectCacheQuery()
		},
		endpoint:     "/measurements",
		expectedCode: http.StatusOK,
		expectedBody: `[
			{
			  "position": {
				"lat": 21,
				"long": 37
			  },
			  "timestamp": "2005-04-02T21:37:00Z",
			  "levelO3": 42
			}
		  ]`,
	}),
)

type opAPITest struct {
	recordMocks  func()
	endpoint     string
	expectedCode int
	expectedBody string
}

func expectCacheQuery() *gomock.Call {
	return cache.EXPECT().GetMeasurements().
		DoAndReturn(func() []measurement.Measurement {
			return measurements
		})
}
