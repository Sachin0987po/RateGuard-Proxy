package config

import (
	"encoding/json"
	"log"
	"io/ioutil"
)

type Endpoint struct {
	Id			   int 	  `json:"id"`
	Path           string `json:"path"`
	RequestsPerSec int    `json:"RequestsPerSec"`
}

type endPointDetialsConfig struct {
	Endpoints []Endpoint `json:"endpoints"`
}

var endptsDetailsConf endPointDetialsConfig


func readEndpointFromJson() bool {
	filePath := "./config/endpoints.json"
	
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
		return false
	}

	err = json.Unmarshal(data, &endptsDetailsConf)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
		return false
	}

	return true
}


func GetEndpointDetail(path string, epDetail *Endpoint) bool {
	if len(endptsDetailsConf.Endpoints) == 0 && !readEndpointFromJson(){
		return false
	}
	for _, endpoint := range endptsDetailsConf.Endpoints {
		if endpoint.Path == path {
			epDetail.Id = endpoint.Id
			epDetail.Path = endpoint.Path
			epDetail.RequestsPerSec = endpoint.RequestsPerSec
			return true
		}
	}
	return false
}