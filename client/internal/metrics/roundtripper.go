package metrics

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	lblCode     = "code"
	lblEndpoint = "endpoint"
)

// TODO TASK-5.2: Implement RoundTripper using promhttp.RoundTripperFunc and record request duration
func InstrumentRoundTripperDuration(obs prometheus.ObserverVec, next http.RoundTripper) promhttp.RoundTripperFunc {
	_ = obs.MustCurryWith(prometheus.Labels{
		lblCode:     "",
		lblEndpoint: "",
	})

	return func(r *http.Request) (*http.Response, error) {
		start := time.Now()
		resp, err := next.RoundTrip(r)
		if err == nil {
			labels := prometheus.Labels{
				lblCode:     strconv.Itoa(resp.StatusCode),
				lblEndpoint: fmt.Sprintf("%s %s", r.Method, r.URL.Path),
			}

			obs.With(labels).Observe(time.Since(start).Seconds())
		}
		return resp, err
	}
}
