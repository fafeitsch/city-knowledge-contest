package geodata

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/city-knowledge-contest/backend/types"
	"io/ioutil"
	"log"
	"math/rand"
	url2 "net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"
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
	Streets  []Street         `json:"streets"`
}

type Street struct {
	Name       string
	Coordinate *types.Coordinate `json:"coord"`
}

func (s *StreetList) GetRandomStreet(random *rand.Rand) Street {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	index := random.Intn(len(s.Streets))
	street := s.Streets[index]
	if street.Coordinate != nil {
		return street
	}
	template := NominatimServer + "/search?street=%s&format=json&city=%s&country=%s"
	url := fmt.Sprintf(template, url2.QueryEscape(street.Name), url2.QueryEscape(s.City), url2.QueryEscape(s.Country))
	response, err := client.Get(url)
	if err != nil {
		log.Printf("could not query nominatim using url \"%s\": %v", url, err)
		return street
	}
	var nominatimResponse []struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	}
	err = json.NewDecoder(response.Body).Decode(&nominatimResponse)
	if err != nil {
		log.Printf("could not parse response from nominatim, using url \"%s\": %v", url, err)
		return street
	}
	if len(nominatimResponse) == 0 {
		log.Printf("could not find street \"%s\" in nominatim using url \"%s\"", street.Name, url)
		return street
	}
	log.Printf("found street \"%s\" in nominatim, updating file", street.Name)
	lat, _ := strconv.ParseFloat(nominatimResponse[0].Lat, 64)
	lon, _ := strconv.ParseFloat(nominatimResponse[0].Lon, 64)
	street.Coordinate = &types.Coordinate{Lat: lat, Lng: lon}
	s.Streets[index] = street
	fileContent, _ := json.MarshalIndent(s, "", " ")
	file, err := os.OpenFile(filepath.Join(StreetListDirectory, s.fileName), os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("could not open file \"%s\": %v", s.fileName, err)
		return street
	}
	_, err = file.Write(fileContent)
	if err != nil {
		log.Printf("could not write file \"%s\": %v", s.fileName, err)
		return street
	}
	return street
}
