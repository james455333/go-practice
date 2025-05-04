package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"net/http"
	"os"
)

func main() {
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
		proxy  = flag.String("proxy", "", "Optional comma-separated list of URLs to proxy uppercase requests")
	)
	flag.Parse()

	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "listen", *listen, "caller", log.DefaultCaller)

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
	svc = proxyingMiddleware(context.Background(), *proxy, logger)(svc)
	svc = loggingMiddleware{logger, svc}
	svc = instrumentingMiddleware{requestCount, requestLatency, counterResult, svc}

	uppercaseHandler := httptransport.NewServer(
		makeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
	)

	countHandler := httptransport.NewServer(
		makeCountEndpoint(svc),
		decodeCountRequest,
		encodeResponse,
	)

	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP", "addr", *listen)
	logger.Log("msg", http.ListenAndServe(*listen, nil))
}

func encodeRequest(_ context.Context, httpRequest *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	httpRequest.Body = io.NopCloser(&buf)
	return nil
}

func decodeUppercaseResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var response uppercaseResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
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
