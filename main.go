package main

import (
	"dashboard_scripts/folder"
	"dashboard_scripts/search"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	endPoint               = "http://localhost:3000/api/"
	getFolderResultFile    = "get_folder_result.json"
	getDashboardResultFile = "get_dashboard_result.json"
	dashboardDataFolder    = "data"
)

func httpClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	return client
}

func getDashboard(c *http.Client, UIDs []string, folderName string) error {
	return nil
}

func main() {
	// clean up the dashboard folder
	err := recreateDataFolder()
	if err != nil {
		log.Fatalf("Error recreating data folder for dashboards %+v", err)
		os.Exit(1)
	}

	// create http client for sending requests
	c := httpClient()

	// search folders and store the map of folder UID <-> folder name
	var folders *folder.GetResponse
	folders, err = folder.GetFolders(c, getFolderResultFile, endPoint)
	if err != nil {
		log.Fatalf("Error occurred during get folders. %+v", err)
		os.Exit(1)
	}

	// create folders in data according to API response
	err = createFoldersInDataRepository(folders, dashboardDataFolder)
	if err != nil {
		log.Fatalf("Error occurred during create folders in data repository. %+v", err)
		os.Exit(1)
	}

	// get dashboards in folders
	for _, folder := range *folders {
		var dbsInFolder *search.GetDBsResponse
		dbsInFolder, err := search.GetDashboardsInFolder(c, getDashboardResultFile, folder, endPoint, dashboardDataFolder)
		if err != nil {
			log.Fatalf("Error occurred during search dashboards. %+v", err)
		}

		for _, dashboard := range *dbsInFolder {
			err := getDashboard(c, dashboardDataFolder, folder, dashboard)
			if err != nil {
				log.Fatalf("Error occurred during get dashboard. %+v", err)
			}
		}
	}
}

func createFoldersInDataRepository(folders *folder.GetResponse, dataFolder string) error {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Can't get current repository. %+v", err)
		return err
	}

	for _, folder := range *folders {
		err = os.Mkdir(filepath.Join(pwd, dataFolder, folder.UID), 0755)
		if err != nil {
			log.Fatalf("Can't create folder. %+v", err)
			return err
		}
	}
	return nil
}

func recreateDataFolder() error {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Can't get current repository. %+v", err)
		return err
	}

	err = os.RemoveAll(filepath.Join(pwd, dashboardDataFolder))
	if err != nil {
		log.Fatalf("Can't remove folder. %+v", err)
		return err
	}

	err = os.Mkdir(filepath.Join(pwd, dashboardDataFolder), 0755)
	if err != nil {
		log.Fatalf("Can't create folder. %+v", err)
		return err
	}
	return nil
}
