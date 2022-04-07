package folder

import (
	"bytes"
	"dashboard_scripts/request"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type GetFolder struct {
	ID            int    `json:"id"`
	UID           string `json:"uid"`
	TITLE         string `json:"title"`
	DashboardUIDs []string
}

type GetResponse []GetFolder

// Get folder ids/names from Grafana API
func GetFolders(c *http.Client, fileName string, endpoint string) (*GetResponse, error) {
	response := request.SendRequest(c, http.MethodGet, endpoint+"folders", nil)
	var result GetResponse
	if err := json.Unmarshal(response, &result); err != nil { // Parse []byte to go struct pointer
		return nil, err
	}

	result = append(result, GetFolder{ID: 0, UID: "0", TITLE: "General"})
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Can't get current repository. %+v", err)
		return &result, err
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, response, "", "\t")
	if err != nil {
		log.Fatalf("Bad search request. %+v", err)
		return &result, err
	}
	err = os.WriteFile(filepath.Join(pwd, fileName), prettyJSON.Bytes(), 0644)
	return &result, err
}
