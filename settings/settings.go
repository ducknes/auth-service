package settings

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	_appName    = "storage-service"
	_configPath = ".config"

	_envKey = "ENV"

	_prod  ENV = "prod"
	_local ENV = "local"
	_dev   ENV = "dev"
)

type ENV string

func AppName() string {
	return _appName
}

func ReadConfig() (Config, error) {
	path := fmt.Sprintf("%s/%s.json", _configPath, GetEnv())
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err = json.Unmarshal(configBytes, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func GetEnv() ENV {
	env := os.Getenv(_envKey)
	if env == "" {
		return _local
	}

	return ENV(env)
}

func LocalEnv() ENV {
	return _local
}
