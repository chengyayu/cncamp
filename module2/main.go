package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"m2/lg"
)

const VERSION = "VERSION"

func init() {
	os.Setenv(VERSION, "2.0")
	lg.Init()
}

func main() {
	http.HandleFunc("/health", mid(healthHandler))
	http.HandleFunc("/", mid(headerHandler))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logrus.Fatal(err)
	}
}

func headerHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "===================Details of the http request header:============\n")
	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	content, _ := json.Marshal(Resp{Code: 200, Data: "ok"})
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

func mid(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Add(VERSION, os.Getenv(VERSION))

		clientIP, err := getIP(r)
		if err != nil {
			logrus.Fatalln(err)
		}

		f(w, r)

		logrus.WithField("clientIP", clientIP).WithField("respCode", w.Header().Get("Status-Code")).Infoln()
	}
}

func getIP(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i, nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	return "", errors.New("no valid ip found")
}

type Resp struct {
	Code int
	Data string
}
