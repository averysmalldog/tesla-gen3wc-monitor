package polly

import (
	"encoding/json"
	"fmt"
	"io/ioutil" // todo(averysmalldog): refactor deprecated package
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

// InfluxAsyncGet polls the Tesla Gen3 Wall Connector and writes results to
// InfluxDB with an aync, non-blocking client you supply. You must also
// supply the IP of the wall conenctor.
func InfluxAsyncGet(writeAPI *api.WriteAPI, wcIP string) {
	client := *writeAPI
	var data map[string]interface{}

	resp, err := http.Get(fmt.Sprintf("http://%s/api/1/vitals", wcIP))
	if err != nil {
		fmt.Println("error - during GET of hpwc. Do you have the right IP?")
		return
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &data)

	// Output a dot (.) for every successful GET against the Wall Connector
	// This helps people like me who need to see something to know it works
	fmt.Printf(".")

	p := influxdb2.NewPoint(
		"hpwc",
		map[string]string{
			"product":  "Gen3 HPWC",
			"vendor":   "Tesla",
			"location": "Garage",
		},
		data,
		time.Now(),
	)
	client.WritePoint(p)
}

// Execute simply runs the totality of polly in your program. It is
// recommended you run this as a goroutine so your program can do
// other things.
func Execute() {
	hpwcIP, exists := os.LookupEnv("HPWC_IP")
	if !exists {
		fmt.Println("env var HPWC_IP not set")
		os.Exit(1)
	}
	influxIP, exists := os.LookupEnv("INFLUX_IP")
	if !exists {
		fmt.Println("env var INFLUX_IP not set")
		os.Exit(1)
	}

	// Check to see if INFLUX_IP has a port specified. If not, supply the default
	if !strings.Contains(influxIP, ":") {
		influxIP += ":8086"
	}

	client := influxdb2.NewClientWithOptions(
		fmt.Sprintf("http://%s", influxIP),
		"my-token",                                  // todo(averysmalldog) un-hardcode this
		influxdb2.DefaultOptions().SetBatchSize(20),
	)
	writeAPI := client.WriteAPI("admin", "admin123")

	// The way this is set up, these likely don't get executed on ^C.
	defer client.Close()
	defer writeAPI.Flush()

	// Simple, isn't it?
	for {
		go InfluxAsyncGet(&writeAPI, hpwcIP)
		time.Sleep(time.Millisecond * 1000)
	}
}
