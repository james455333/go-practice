package main

import (
	"context"
	"encoding/json"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)

	fieldKey := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_count",
		Help:      "Number of request received.",
	}, fieldKey)

	requestLatency := kitprometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKey)

	counterResult := kitprometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "counter_result",
		Help:      "the result of each count result",
	}, fieldKey)

	var svc StringService
	svc = stringService{}
	svc = loggingMiddleware{logger, svc}
	svc = instrumentingMiddleware{requestCount, requestLatency, counterResult, svc}

	uppercaseHandler := httptransport.NewServer(
		makeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
	)

	var count StringService
	count = stringService{}
	count = loggingMiddleware{logger, count}
	count = instrumentingMiddleware{requestCount, requestLatency, counterResult, count}
	countHandler := httptransport.NewServer(
		makeCountEndpoint(count),
		decodeCountRequest,
		encodeResponse,
	)

	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "listening on :8080")
	logger.Log("msg", http.ListenAndServe(":8080", nil))
}

func decodeCountRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	var countRequest countRequest
	if err := json.NewDecoder(request.Body).Decode(&countRequest); err != nil {
		return nil, err
	}
	return countRequest, nil
}

func decodeUppercaseRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	var uppercaseRequest uppercaseRequest
	if err := json.NewDecoder(request.Body).Decode(&uppercaseRequest); err != nil {
		return nil, err
	}
	return uppercaseRequest, nil
}

func encodeResponse(ctx context.Context, writer http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(writer).Encode(response)
}
