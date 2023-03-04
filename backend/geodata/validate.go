package geodata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fafeitsch/city-knowledge-contest/backend/types"
)

var client = http.Client{Timeout: 60 * time.Second}

type nominatimReverseResponse struct {
	Address struct {
		Road string `json:"road"`
	} `json:"address"`
}

func sendNominatimRequest(c types.Coordinate) (nominatimReverseResponse, error) {
	url := fmt.Sprintf("%s/reverse?format=json&lat=%f&lon=%f&zoom=17&addressdetails=1", NominatimServer, c.Lat, c.Lng)
	resp, err := client.Get(url)

	if err != nil {
		return nominatimReverseResponse{}, fmt.Errorf("could not query Nominatim address: %v", err)
	}

	var nomiReponse nominatimReverseResponse
	err = json.NewDecoder(resp.Body).Decode(&nomiReponse)

	if err != nil {
		return nominatimReverseResponse{}, fmt.Errorf("could not parse response from Nominatim: %v", err)
	}

	return nomiReponse, nil
}

func VerifyAnswer(guess types.Coordinate, answer string) (bool, error) {
	result, err := sendNominatimRequest(guess)
	if err != nil {
		return false, err
	}

	if result.Address.Road == answer {
		return true, nil
	}

	return false, nil
}
