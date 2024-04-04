package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

// TestMain mocks a Tesla Gen3 Wall Connector and an InfluxDB instance. Since
// both of these dependencies are "up", we should expect to see a "." in
// standard out. Since I don't have a good way to terminate main, these aren't
// real tests and the output keeps stacking over the 20s runtime. However,
// we should expect to see only dots (".") for the first 5 seconds, and then
// that plus various error messages stacking for the next 15 seconds.
func TestMain(t *testing.T) {
	testTable := []struct {
		name         string
		hpwcServer   *httptest.Server
		influxServer *httptest.Server
	}{
		{
			name: "happy-wall-connector-and-happy-influx",
			hpwcServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"voltage": 240, "contactor_status": "1", "amperage": 1}`))
			})),
			influxServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status": "OK"}`))
			})),
		},
		{
			name: "happy-wall-connector-and-sad-influx",
			hpwcServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"voltage": 240, "contactor_status": "1", "amperage": 1}`))
			})),
			influxServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			})),
		},
		{
			name: "sad-wall-connector-and-happy-influx",
			hpwcServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			})),
			influxServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status": "OK"}`))
			})),
		},
		{
			name: "sad-wall-connector-and-sad-influx",
			hpwcServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			})),
			influxServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			})),
		},
	}
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			defer tc.hpwcServer.Close()
			defer tc.influxServer.Close()

			// Find the IP:port of the two testing servers and set
			// (and schedule unset of) env vars
			fmt.Println(tc.hpwcServer.URL)
			fmt.Println(tc.influxServer.URL)
			myTeslaUrl, err := url.Parse(tc.hpwcServer.URL)
			if err != nil {
				t.Fail()
			}
			myInfluxUrl, err := url.Parse(tc.hpwcServer.URL)
			if err != nil {
				t.Fail()
			}
			teslaIP := myTeslaUrl.Host
			influxIP := myInfluxUrl.Host
			os.Setenv("HPWC_IP", fmt.Sprint(teslaIP))
			os.Setenv("INFLUX_IP", fmt.Sprint(influxIP))
			defer os.Unsetenv("HPWC_IP")
			defer os.Unsetenv("INFLUX_IP")

			// main will just live on forever and keep outputting
			// throughout the test
			go main()
			time.Sleep(time.Millisecond * 5000)
		})
	}
}
