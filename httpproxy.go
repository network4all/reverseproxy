package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/koding/websocketproxy"
	"github.com/sirupsen/logrus"
)

func addURL(r proxyUrl) error {
	// decode url
	url, err := url.Parse(r.internalUrl)
	if err != nil {
		return err
	}

	// ws or wss
	if strings.HasPrefix(r.internalUrl, "ws") {
		addproxy := websocketproxy.NewProxy(url)
		addproxy.Director = func(req *http.Request, out http.Header) {
			out.Set("devicename", req.Header.Get("devicename"))
		}
		http.HandleFunc(r.externalUrl, handlerWebsocketProxy(addproxy, r.token))

		return nil
	}

	// http or https
	addproxy := httputil.NewSingleHostReverseProxy(url)
	addproxy.Transport = &myTransport{}
	http.HandleFunc(r.externalUrl, handlerSecureProxy(addproxy, r.token))
	return nil
}

func handlerSecureProxy(p *httputil.ReverseProxy, requiredtoken string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authenticate")
		if token != requiredtoken {
			logrus.Warningf("Proxy: Unauthenticated token found for this service: '%s' for %s\n", token, r.URL)
			fmt.Fprint(w, "ReverseProxy: Authentication failed!")
		}
		p.ServeHTTP(w, r)
	}
}

func handlerWebsocketProxy(p *websocketproxy.WebsocketProxy, requiredtoken string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if requiredtoken != "" {
			token := r.Header.Get("Authenticate")
			devicename := r.Header.Get("devicename")
			if token != requiredtoken {

				logrus.Warningf("Proxy: Unauthenticated token found for this service: '%s' for %s (%s)\n", token, devicename, r.URL)
				fmt.Fprint(w, "ReverseProxy: Authentication failed!")
				return
			}
			logrus.Infof("Proxy: Authenticated with token: '%s' for %s\n", devicename, r.URL)
		}
		p.ServeHTTP(w, r)
	}
}
