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
		Name     string     `json:"name"`
		Location [2]float64 `json:"location"`
	} `json:"waypoints"`
}

func queryStreet(c Coordinate, random *rand.Rand) (string, Coordinate, error) {
	url, osrmWaypoints, err := sendOsrmRequest(c)
	if err != nil {
		return "", Coordinate{}, err
	}
	random.Shuffle(len(osrmWaypoints.Waypoints), func(i, j int) {
		osrmWaypoints.Waypoints[i], osrmWaypoints.Waypoints[j] = osrmWaypoints.Waypoints[j], osrmWaypoints.Waypoints[i]
	})
	name := ""
	coordinate := Coordinate{}
	for _, wp := range osrmWaypoints.Waypoints {
		if wp.Name != "" {
			name, coordinate = wp.Name, Coordinate{Lng: wp.Location[0], Lat: wp.Location[1]}
		}
	}
	log.Printf("Selected random point: %s, https://www.openstreetmap.org/?mlat=%f&mlon=%f#map=18/%f/%f → %s", url,
		c.Lat, c.Lng, c.Lat,
		c.Lng, name)
	return name, coordinate, nil
}

func sendOsrmRequest(c Coordinate) (string, osrmAddressResponse, error) {
	url := fmt.Sprintf("%s/nearest/v1/driving/%f,%f.json?number=10", OsrmServer, c.Lng, c.Lat)
	osrmResp, err := client.Get(url)
	if err != nil {
		return "", osrmAddressResponse{}, fmt.Errorf("could not query OSRM address: %v", err)
	}
	var osrmWaypoints osrmAddressResponse
	err = json.NewDecoder(osrmResp.Body).Decode(&osrmWaypoints)
	if err != nil {
		return "", osrmAddressResponse{}, fmt.Errorf("could not parse response from osrm: %v", err)
	}
	return url, osrmWaypoints, nil
}
func verifyAnswer(guess Coordinate, answer string) (bool, error) {
	_, waypoints, err := sendOsrmRequest(guess)
	if err != nil {
		return false, err
	}
	for _, wp := range waypoints.Waypoints {
		if wp.Name == answer {
			return true, nil
		}
	}
	return false, nil
}
