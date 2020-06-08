package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TODO TASK-5.2: Implement RoundTripper using promhttp.RoundTripperFunc and record request duration
func InstrumentRoundTripperDuration(obs prometheus.ObserverVec, next http.RoundTripper) promhttp.RoundTripperFunc {
	return func(r *http.Request) (*http.Response, error) {
		return next.RoundTrip(r)
	}
}
