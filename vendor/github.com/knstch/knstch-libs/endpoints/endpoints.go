package endpoints

import (
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/knstch/knstch-libs/middleware"
	metrics "github.com/knstch/knstch-libs/prometeus"
	"github.com/knstch/knstch-libs/tracing"
	"github.com/knstch/knstch-libs/transport"

	"github.com/go-chi/chi/v5"
)

type Endpoint struct {
	Method  string
	Path    string
	Handler endpoint.Endpoint
	Decoder httptransport.DecodeRequestFunc
	Encoder httptransport.EncodeResponseFunc
	Mdw     []middleware.Middleware
	Opts    []httptransport.ServerOption
}

func InitHttpEndpoints(serviceName string, endpoints []Endpoint) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.WithTrackingRequests)

	metrics.InitBasicMetrics()

	for _, ep := range endpoints {
		handler := ep.Handler

		handler = tracing.WithTracing()(handler)

		for _, mw := range ep.Mdw {
			handler = mw(handler)
		}

		opts := append(ep.Opts,
			httptransport.ServerErrorEncoder(transport.EncodeError),
			httptransport.ServerBefore(httptransport.PopulateRequestContext),
		)

		r.Method(ep.Method, ep.Path, httptransport.NewServer(
			handler,
			ep.Decoder,
			ep.Encoder,
			opts...,
		))
	}

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	r.Method(http.MethodGet, "/metrics", promhttp.Handler())

	return r
}
