package main

import (
	"context"
	"encoding/json"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"net/http"
	"os"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	var svc StringService
	svc = stringService{}
	svc = loggingMiddleware{logger, svc}

	uppercaseHandler := httptransport.NewServer(
		makeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
	)

	var count StringService
	count = stringService{}
	count = loggingMiddleware{logger, count}
	countHandler := httptransport.NewServer(
		makeCountEndpoint(count),
		decodeCountRequest,
		encodeResponse,
	)

	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
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
