package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"
)

type Endpoint struct {
	Id             int    `json:"id"`
	Path           string `json:"path"`
	RequestsPerMin int    `json:"RequestsPerMin"`
	pattern        *regexp.Regexp
}

type endPointDetialsConfig struct {
	Endpoints []Endpoint `json:"endpoints"`
}

var endptsDetailsConf endPointDetialsConfig

func compileRegexp(endptsDetailsConf *endPointDetialsConfig) {
	for i, endpoint := range endptsDetailsConf.Endpoints {
		regex, _ := regexp.Compile(endpoint.Path)
		endptsDetailsConf.Endpoints[i].pattern = regex
	}
}

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

	compileRegexp(&endptsDetailsConf)

	return true
}

func GetEndpointDetail(path string, epDetail *Endpoint) bool {
	if len(endptsDetailsConf.Endpoints) == 0 && !readEndpointFromJson() {
		return false
	}
	for _, endpoint := range endptsDetailsConf.Endpoints {
		if endpoint.Path == path || endpoint.pattern.MatchString(path) {
			epDetail.Id = endpoint.Id
			epDetail.Path = endpoint.Path
			epDetail.RequestsPerMin = endpoint.RequestsPerMin
			return true
		}
	}
	return false
}
