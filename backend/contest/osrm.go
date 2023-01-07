package contest

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var OsrmServer = "http://127.0.0.1:5000"
var client = http.Client{Timeout: 60 * time.Second}

type osrmAddressResponse struct {
	Waypoints []struct {
		Name string `json:"name"`
	} `json:"waypoints"`
}

func queryStreet(c Coordinate, random *rand.Rand) (string, error) {
	url := fmt.Sprintf("%s/nearest/v1/driving/%f,%f.json?number=10", OsrmServer, c.Lng, c.Lat)
	log.Printf(url)
	log.Printf("https://www.openstreetmap.org/?mlat=%f&mlon=%f#map=18/%f/%f", c.Lat, c.Lng, c.Lat, c.Lng)
	osrmResp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("could not query OSRM address: %v", err)
	}
	var osrmWaypoints osrmAddressResponse
	err = json.NewDecoder(osrmResp.Body).Decode(&osrmWaypoints)
	log.Printf("%v", osrmWaypoints)
	if err != nil {
		return "", fmt.Errorf("could not parse response from osrm: %v", err)
	}
	name := ""
	random.Shuffle(len(osrmWaypoints.Waypoints), func(i, j int) {
		osrmWaypoints.Waypoints[i], osrmWaypoints.Waypoints[j] = osrmWaypoints.Waypoints[j], osrmWaypoints.Waypoints[i]
	})
	for _, wp := range osrmWaypoints.Waypoints {
		if wp.Name != "" {
			name = wp.Name
		}
	}
	return name, nil
}
