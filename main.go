package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	testEnvEndPoint = "http://localhost:3000/api/"
)

func httpClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	return client
}

func sendRequest(client *http.Client, method string, endpoint string, parameters map[string]string) []byte {
	jsonData, err := json.Marshal(parameters)
	if err != nil {
		log.Fatalf("Error Occurred. %+v", err)
	}

	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error Occurred. %+v", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("apikey", os.Getenv("GRAFANA_API_KEY"))

	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request to API endpoint. %+v", err)
	}

	// Close the connection to reuse it
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Couldn't parse response body. %+v", err)
	}

	return body
}

func searchDashboards(c *http.Client) {
	parameters := map[string]string{
		"query":      "",
		"starred":    "false",
		"skipRecent": "true",
		"folderId":   "0",
		"layout":     "folders",
		"prevSort":   "null",
	}
	response := sendRequest(c, http.MethodGet, testEnvEndPoint+"search", parameters)
	log.Println("Response Body:", string(response))
}

func main() {
	c := httpClient()
	searchDashboards(c)
}
