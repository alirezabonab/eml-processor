package project

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Projects []Project `json:"projects"`
}

type Project struct {
	Keywords  []string `json:"keywords"`
	DestDir   string   `json:"destDir"`
	SourceDir string   `json:"SourceDir"`
	ClearDest bool     `json:"clearDest"`
	Name      string   `json:"name"`
}

func GetAllProjects(configFilePath string) ([]Project, error) {
	var config Config
	// read config file
	r, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// unmarshal config file
	err = json.Unmarshal(r, &config)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return config.Projects, nil
}
