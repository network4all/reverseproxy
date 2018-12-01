package main

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

func handler404(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
	fmt.Fprint(w, "ReverseProxy: 404 Page not found!")
	logrus.Infof("404 error on URL %s", r.URL)
}
