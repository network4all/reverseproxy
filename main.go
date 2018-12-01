package main

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
)

var (
	whiteList        = os.Getenv("WHITELIST")
	certificatesPath = os.Getenv("CERT_PATH")
	logfile			 = os.Getenv("NETWORK4ALL_LOGPATH")
)

func main() {
	// disable cert checking for internal (selfsigned) webservers
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// split up domains
	domains := strings.Split(whiteList, ";")
	if len(domains) == 0 || len(whiteList) == 0 {
		logrus.Warningf("could not get the domain(s) from whitelist %s", domains)
	}

	// user certmanager to autogenerate certificates
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domains...),
		Cache:      autocert.DirCache(certificatesPath + "certs"),
	}

	// create a website with certificates from letsencrypt
	server := &http.Server{
		Addr: ":https",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	// autocert & redirect port 80
	go func() {
		logrus.Fatal(http.ListenAndServe(":http", certManager.HTTPHandler(nil)))
	}()

	// add sites
	configuration := []byte(`[{"application":"website1", "internalurl":"http://127.0.0.1:6000/", "externalurl":"www.network4all.nl/aa/", "token":"", "enabled":true}, {"application":"websocket2", "internalurl":"ws://127.0.0.1:6001/", "externalurl":"www.network4all.nl/bb/", "token":"aaa", "enabled":true}, {"application":"websocket3","internalurl":"ws://127.0.0.1:6002/", "externalurl":"www.network4all.nl/cc/", "token":"bbb", "enabled":false}]`)
	var websites []proxyUrl
	err := json.Unmarshal(configuration, &websites)
	if err != nil {
		logrus.Fatalf("could not decode configuration file %v", err)
	}
	for _, website := range websites {
		if website.enabled {
			err := addURL(website)
			if err != nil {
				logrus.Warningf("could not add proxy for %v", err)
			}
			logrus.Info("redirecting the application %s with url %s to %s", website.application, website.externalUrl, website.internalUrl)
		}
	}

	// catch rest
	http.HandleFunc("/", handler404)
	logrus.Infoln("Starting the network4all proxy server")
	logrus.Fatal(server.ListenAndServeTLS("", ""))
}
