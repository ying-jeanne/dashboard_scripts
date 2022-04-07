package search

import (
	"bytes"
	"dashboard_scripts/folder"
	"dashboard_scripts/request"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func GetDashboardsInFolder(c *http.Client, fileName string, folder folder.GetFolder, endPoint string, dashboardDataFolder string) error {
	parameters := map[string]string{
		"query":      "",
		"starred":    "false",
		"skipRecent": "true",
		"folderId":   strconv.Itoa(folder.ID),
		"layout":     "folders",
		"prevSort":   "null",
	}
	response := request.SendRequest(c, http.MethodGet, endPoint+"search", parameters)

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Can't get current repository. %+v", err)
		return err
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, response, "", "\t")
	if err != nil {
		log.Fatalf("Bad search request. %+v", err)
		return err
	}
	err = os.WriteFile(filepath.Join(pwd, dashboardDataFolder, folder.UID, fileName), prettyJSON.Bytes(), 0644)
	return err
}
