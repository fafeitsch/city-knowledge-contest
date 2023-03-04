package geodata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	url2 "net/url"
	"path/filepath"
	"sync"

	"github.com/fafeitsch/city-knowledge-contest/backend/types"
)

var StreetListDirectory = "./streetlists"
var NominatimServer = "https://nominatim.openstreetmap.org"

type StreetListHeader struct {
	FileName string
	Name     string           `json:"name"`
	Center   types.Coordinate `json:"center"`
}

func ReadStreetLists() ([]StreetListHeader, error) {
	files, err := ioutil.ReadDir(StreetListDirectory)
	if err != nil {
		log.Printf("could not list files: %v", err)
		return nil, fmt.Errorf("could not list files")
	}
	result := make([]StreetListHeader, 0, 0)
	for _, file := range files {
		var streetListFile StreetListHeader
		fileContent, err := ioutil.ReadFile(filepath.Join(StreetListDirectory, file.Name()))
		if err != nil {
			log.Printf("could not read file \"%s\": %v", file.Name(), err)
			continue
		}
		streetListFile.FileName = file.Name()
		err = json.Unmarshal(fileContent, &streetListFile)
		if err != nil {
			log.Printf("could not parse file \"%s\" as streetListFile: %v", file.Name(), err)
			continue
		}
		result = append(result, streetListFile)
	}
	return result, nil
}

var streetLists = make(map[string]*StreetList)
var mapMutex = &sync.Mutex{}

func ReadStreetList(fileName string) (*StreetList, error) {
	mapMutex.Lock()
	defer mapMutex.Unlock()
	if streetList, ok := streetLists[fileName]; ok {
		return streetList, nil
	}
	var streetList *StreetList
	fileContent, err := ioutil.ReadFile(filepath.Join(StreetListDirectory, fileName))
	if err != nil {
		return streetList, fmt.Errorf("could not read file \"%s\"", fileName)
	}
	err = json.Unmarshal(fileContent, &streetList)
	if err != nil {
		return streetList, fmt.Errorf("could not parse file \"%s\" as streetListFile", fileName)
	}
	streetList.fileName = fileName
	if len(streetList.Streets) == 0 {
		return streetList, fmt.Errorf("file \"%s\" does not contain any streets", fileName)
	}
	streetLists[fileName] = streetList
	return streetList, err
}

type StreetList struct {
	mutex    sync.Mutex
	fileName string
	Country  string           `json:"country"`
	City     string           `json:"city"`
	Name     string           `json:"name"`
	Center   types.Coordinate `json:"center"`
	Streets  []string         `json:"streets"`
}

type Street struct {
	Name  string
	Lines []Line
}

type Line struct {
	Coordinates []types.Coordinate
}

func (s *StreetList) GetRandomStreet(random *rand.Rand) (Street, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	index := random.Intn(len(s.Streets))
	street := s.Streets[index]

	template := NominatimServer + "/search?street=%s&format=json&city=%s&country=%s&dedupe=0&polygon_geojson=1"
	url := fmt.Sprintf(template, url2.QueryEscape(street), url2.QueryEscape(s.City), url2.QueryEscape(s.Country))
	response, err := client.Get(url)
	if err != nil {
		log.Printf("could not query nominatim using url \"%s\": %v", url, err)
		return Street{}, err
	}

	var nominatimResponses []struct {
		Lat     string `json:"lat"`
		Lon     string `json:"lon"`
		GeoJson struct {
			Type        string      `json:"type"`
			Coordinates [][]float64 `json:"coordinates"`
		} `json:"geojson"`
	}

	err = json.NewDecoder(response.Body).Decode(&nominatimResponses)
	if err != nil {
		log.Printf("could not parse response from nominatim, using url \"%s\": %v", url, err)
		return Street{}, err
	}
	if len(nominatimResponses) == 0 {
		log.Printf("could not find street \"%s\" in nominatim using url \"%s\"", street, url)
		return Street{}, err
	}

	log.Printf("found street \"%s\" in nominatim using url \"%s\", updating file", street, url)

	var lines []Line

	for _, response := range nominatimResponses {
		var line Line
		for _, cor := range response.GeoJson.Coordinates {
			line.Coordinates = append(line.Coordinates,
				types.Coordinate{Lat: cor[1], Lng: cor[0]})
		}
		lines = append(lines, line)
	}

	return Street{Name: street, Lines: lines}, nil
}
