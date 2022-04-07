package request

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func SendRequest(client *http.Client, method string, endpoint string, parameters map[string]string) []byte {
	jsonData, err := json.Marshal(parameters)
	if err != nil {
		log.Fatalf("Error Occurred. %+v", err)
	}

	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error Occurred. %+v", err)
	}

	req.Header.Add("Accept", "application/json")
	// req.Header.Add("apikey", os.Getenv("GRAFANA_API_KEY"))
	req.Header.Add("Authorization", "Bearer eyJrIjoieXZ2SUs0clVFeWZXSzZHNkZsM0U5NWJ6SjlYYjRidGsiLCJuIjoic2NyaXB0IiwiaWQiOjF9")

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
