package geodata

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/city-knowledge-contest/backend/types"
	"net/http"
	"time"
)

var OsrmServer = "http://127.0.0.1:5000"
var client = http.Client{Timeout: 60 * time.Second}

type osrmAddressResponse struct {
	Waypoints []struct {
		Name     string     `json:"name"`
		Location [2]float64 `json:"location"`
		Distance float64    `json:"distance"`
	} `json:"waypoints"`
}

func sendOsrmRequest(c types.Coordinate) (string, osrmAddressResponse, error) {
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

func VerifyAnswer(guess types.Coordinate, answer string) (bool, error) {
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
