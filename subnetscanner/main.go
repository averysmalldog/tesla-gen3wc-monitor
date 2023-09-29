package main

import (
    "encoding/json"
    "fmt"
    "net"
    "net/http"
)

type Vital struct {
    Host string `json:"host"`
    Status int `json:"status"`
}

func main() {
    // Get the subnet information
    subnet := getSubnet()

    // Scan all the routable IPs within the subnet to identify active IPs
    activeIPs := scanSubnet(subnet)

    // Test each active IP with a simple HTTP GET at the following endpoint: /api/1/vitals
    var respondingIPs []string
    for _, ip := range activeIPs {
        // Create an HTTP client
        client := &http.Client{}

        // Create an HTTP GET request
        req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/api/1/vitals", ip), nil)
        if err != nil {
            fmt.Println(err)
            continue
        }

        // Set the request headers
        req.Header.Set("Accept", "application/json")

        // Send the request and get the response
        resp, err := client.Do(req)
        if err != nil {
            fmt.Println(err)
            continue
        }

        // Check the response status code
        if resp.StatusCode != 200 {
            fmt.Println(fmt.Sprintf("GET %s/api/1/vitals failed with status code %d", ip, resp.StatusCode))
            continue
        }

        // Decode the JSON response
        var vital Vital
        err = json.NewDecoder(resp.Body).Decode(&vital)
        if err != nil {
            fmt.Println(err)
            continue
        }

        // Store the IP if the endpoint responded with JSON
        respondingIPs = append(respondingIPs, ip)
    }

    // Return a slice of IPs that respond with JSON
    fmt.Println(respondingIPs)
}

func getSubnet() *net.IPNet {
    // Get the default gateway
    interfaces, err := net.Interfaces()
    if err != nil {
        fmt.Println(err)
        return nil
    }

    for _, i := range interfaces {
        addrs, err := i.Addrs()
        if err != nil {
            fmt.Println(err)
            continue
        }

        for _, addr := range addrs {
            // Check if the address is a gateway
            if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() && ipnet.IP.To4() != nil {
                // Return the subnet information
                return ipnet
            }
        }
    }

    return nil
}

func scanSubnet(subnet *net.IPNet) []string {
    // Create a slice to store the active IPs
    activeIPs := []string{}

    // Iterate over all the routable IPs in the subnet
    for i := subnet.IP.Mask(subnet.Mask).To4(); i.Cmp(subnet.IP.Mask(subnet.Mask).To4().Add(subnet.Mask.Size())); i.Inc(1) {
        // Try to connect to the IP on port 80
        conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: i})
        if err != nil {
            // The IP is not active
            continue
        }

        // Close the connection
        conn.Close()

        // The IP is active
        activeIPs = append(activeIPs, i.String())
    }

    return activeIPs
}