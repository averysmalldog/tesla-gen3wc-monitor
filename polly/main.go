package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	// Grab HPWC_IP from envvars:
	hpwcIP := os.Getenv("HPWC_IP")
    // Create a new client using an InfluxDB server base URL and an authentication token
    // and set batch size to 20 
    client := influxdb2.NewClientWithOptions("http://localhost:8086", "my-token",
        influxdb2.DefaultOptions().SetBatchSize(20))
    // Get non-blocking write client
    writeAPI := client.WriteAPI("admin","tesla")
    // write some points
    for {
		// query the HPWC for info
		var data map[string]interface{}
		resp, _ := http.Get(fmt.Sprintf("http://%s/api/1/vitals",hpwcIP))
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &data)
		fmt.Printf("%+v\n", data)
        // create point
        p := influxdb2.NewPoint(
            "hpwc",
            map[string]string{
                "product":  "Gen3 HPWC",
                "vendor":   "Tesla",
                "location": "Garage",
            },
            data,
            time.Now())
        // write asynchronously
        writeAPI.WritePoint(p)
		time.Sleep(time.Millisecond*1000)
    }
    // Force all unwritten data to be sent
    writeAPI.Flush()
    // Ensures background processes finishes
    client.Close()
}
