package main

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type myTransport struct {
}

func (t *myTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	start := time.Now()
	response, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		return nil, err
	}
	elapsed := time.Since(start)
	logrus.Infof("Proxy: %s %s %s %d %d", elapsed, request.Method, request.RequestURI, response.StatusCode, response.ContentLength)
	return response, err
}

type proxyUrl struct {
	application string
	internalUrl string
	externalUrl string
	token       string
	enabled     bool
}
