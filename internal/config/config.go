package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}
type State struct {
	Config *Config
}

func Read() (Config, error) {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	finalDir := homeDirectory + "/" + configFileName

	data, err := os.ReadFile(finalDir)
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func (cfg *Config) SetUser(username string) error {
	instance := Config{Db_url: cfg.Db_url, Current_user_name: username}
	data, err := json.Marshal(instance)
	if err != nil {
		return err
	}
	permissions := os.FileMode(0644)
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error getting home directory!")
		return err
	}
	finalDir := homeDirectory + "/" + configFileName

	err = os.WriteFile(finalDir, data, permissions)
	if err != nil {
		fmt.Println("error writing to file!")
		return err
	}
	return nil

}
