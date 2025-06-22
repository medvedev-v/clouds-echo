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

type job struct {
	index int
	url   string
}

type result struct {
	index int
	resp  CloudResponse
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

	jobs := make(chan job, len(cloudsInfo))
	results := make(chan result, len(cloudsInfo))

	var wg sync.WaitGroup
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				results <- result{job.index, Get(job.url)}
			}
		}()
	}

	for i, cloud := range cloudsInfo {
		jobs <- job{i, cloud.URL}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		responses[res.index] = res.resp
	}

	return responses
}
