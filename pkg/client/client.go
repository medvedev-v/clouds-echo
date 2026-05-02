package client

import (
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"sync"
	"time"
)

type RequestResults []struct {
	URL  string `json:"url"`
	Ping int    `json:"ping"`
	Code string `json:"code"`
}

type RequestResult struct {
	URL  string `json:"url"`
	Ping int    `json:"ping"`
	Code string `json:"code"`
}

type Clouds []struct {
	URL string `yaml:"url"`
}

type Cloud struct {
	URL string `yaml:"url"`
}

type ForbiddenResources []struct {
	URL string `yaml:"url"`
}

type ForbiddenResource struct {
	URL string `yaml:"url"`
}

type AllowedResources []struct {
	URL string `yaml:"url"`
}

type AllowedResource struct {
	URL string `yaml:"url"`
}

type job struct {
	index int
	url   string
}

type result struct {
	index int
	resp  RequestResult
}

func Get(url string) (requestResult RequestResult) {
	start := time.Now()
	client := &http.Client{
		Timeout: 5 * time.Second, // добавляем таймаут для запроса
	}

	response, err := client.Get(url)
	if err != nil {
		end := time.Now()
		duration := end.Sub(start).Milliseconds()

		requestResult.URL = url
		requestResult.Code = "ERROR: " + err.Error()
		requestResult.Ping = int(duration)
		return
	}

	end := time.Now()
	duration := end.Sub(start).Milliseconds()

	requestResult.URL = url
	requestResult.Code = response.Status
	requestResult.Ping = int(duration)

	defer response.Body.Close()
	return
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

func PingClouds() RequestResults {
	cloudsInfo := GetCloudsInfo()
	responses := make(RequestResults, len(cloudsInfo))

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

func GetForbiddenInfo() (forbiddenResources ForbiddenResources) {
	yamlFile, error := os.ReadFile("forbidden.yaml")
	if error != nil {
		panic(error)
	}
	error = yaml.Unmarshal(yamlFile, &forbiddenResources)
	if error != nil {
		panic(error)
	}

	return
}

func PingForbidden() RequestResults {
	forbiddenInfo := GetForbiddenInfo()
	forbiddenLength := len(forbiddenInfo)

	responses := make(RequestResults, forbiddenLength)

	jobs := make(chan job, forbiddenLength)
	results := make(chan result, forbiddenLength)

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

	for i, cloud := range forbiddenInfo {
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

func GetAllowedInfo() (allowedResources AllowedResources) {
	yamlFile, error := os.ReadFile("allowed.yaml")
	if error != nil {
		panic(error)
	}
	error = yaml.Unmarshal(yamlFile, &allowedResources)
	if error != nil {
		panic(error)
	}

	return
}

func PingAllowed() RequestResults {
	allowedInfo := GetAllowedInfo()
	allowedLength := len(allowedInfo)

	responses := make(RequestResults, allowedLength)

	jobs := make(chan job, allowedLength)
	results := make(chan result, allowedLength)

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

	for i, cloud := range allowedInfo {
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
