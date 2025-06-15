package client

import (
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"sync"
	"time"
)

type CloudsResponses []struct {
	URL          string `json:"url"`
	Ping         int    `json:"ping"`
	ResponseCode string `json:"responsecode"`
}

type CloudResponse struct {
	URL          string `json:"url"`
	Ping         int    `json:"ping"`
	ResponseCode string `json:"responsecode"`
}

type Clouds []struct {
	URL string `yaml:"url"`
}

type Cloud struct {
	URL string `yaml:"url"`
}

func GetCloudsInfo() (clouds Clouds) {
	yamlFile, error := os.ReadFile("clouds.yaml")
	if error != nil {
		panic(error)
	}
	error = yaml.Unmarshal(yamlFile, &clouds)
	if error != nil {
		panic(error)
	}

	return
}

func Get(url string) (cloudResponse CloudResponse) {
	start := time.Now()
	client := &http.Client{
		Timeout: 5 * time.Second, // добавляем таймаут для запроса
	}

	response, err := client.Get(url)
	if err != nil {
		end := time.Now()
		duration := end.Sub(start).Milliseconds()
		
		cloudResponse.URL = url
		cloudResponse.ResponseCode = "ERROR: " + err.Error()
		cloudResponse.Ping = int(duration)
		return
	}

	end := time.Now()
	duration := end.Sub(start).Milliseconds()

	cloudResponse.URL = url
	cloudResponse.ResponseCode = response.Status
	cloudResponse.Ping = int(duration)

	defer response.Body.Close()
	return
}

func PingClouds() CloudsResponses {
	cloudsInfo := GetCloudsInfo()
	responses := make(CloudsResponses, len(cloudsInfo))

	var wg sync.WaitGroup
	wg.Add(len(cloudsInfo))

	for i, cloud := range cloudsInfo {
		go func(index int, url string) {
			defer wg.Done()
			responses[index] = Get(url)
		}(i, cloud.URL)
	}

	wg.Wait()
	return responses
}
