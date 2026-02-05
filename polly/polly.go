package polly

import (
	"encoding/json"
	"fmt"
	"io/ioutil" // todo(averysmalldog): refactor deprecated package
	"net/http"
	"os"
	"time"
	"strings"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)


// Host structure to hold TWC name/IP pair
type host struct {
    name string
    ip string
}


// Function to create a new host structure
func newHost(name string, ip string) *host {
    h := host{name: name}
    h.ip  = ip
    return &h
}


// InfluxAsyncGet polls the Tesla Gen3 Wall Connector and writes results to
// InfluxDB with an aync, non-blocking client you supply. You must also
// supply the IP of the wall conenctor.
func InfluxAsyncGet(writeAPI *api.WriteAPI, measurementName string, wcIP string) {
	client := *writeAPI
	var data map[string]interface{}

	url := fmt.Sprintf("http://%s/api/1/vitals", wcIP)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error - during GET of '", url, "'. Do you have the right IP?")
		return
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &data)

	// Output a dot (.) for every successful GET against the Wall Connector
	// This helps people like me who need to see something to know it works
	fmt.Printf(".")

	p := influxdb2.NewPoint(
		measurementName,
		map[string]string{
			"product":  "Gen3 HPWC",
			"vendor":   "Tesla",
			"location": wcIP,
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
	hpwcIPs := strings.Split(os.Getenv("HPWC_IP"), ",")
	influxIP := os.Getenv("INFLUX_IP")
	InfluxToken := os.Getenv("INFLUX_TOKEN")
	InfluxOrg := os.Getenv("INFLUX_ORG")
	InfluxBucket := os.Getenv("INFLUX_BUCKET")


	// Parse wall connector host:IPs
	var hosts = make([]host,0)

	if (len(hpwcIPs) == 1 && len(strings.Split(hpwcIPs[0], ":")) == 1) {
		// Presume legacy format of just an IP address
		hosts = append(hosts, *newHost("hpwc", hpwcIPs[0]))
	} else {
		// Parse new name:IP Address format
		for _,hpwcIP := range hpwcIPs {
			var pair = strings.Split(hpwcIP, ":")
			if (len(pair) != 2) {
				fmt.Print("Error in HPWC_IP declaration: '", hpwcIP, "'\n")
				os.Exit(1)
			}
			hosts = append(hosts, *newHost(pair[0], pair[1]))
		}
	}

	client := influxdb2.NewClientWithOptions(
		fmt.Sprintf("http://%s:8086", influxIP),
		InfluxToken,
		influxdb2.DefaultOptions().SetBatchSize(20),
	)
	writeAPI := client.WriteAPI(InfluxOrg, InfluxBucket)

	// The way this is set up, these likely don't get executed on ^C.
	defer client.Close()
	defer writeAPI.Flush()

	// Simple, isn't it?
	for {
		for _,host := range hosts {
			go InfluxAsyncGet(&writeAPI, host.name, host.ip)
		}
		time.Sleep(time.Millisecond * 1000)
	}
}
