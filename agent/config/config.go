package config

import (
	"os"
	"strconv"

	"github.com/dkotoff/daec-ylyceum/agent/logger"
)

const (
	defaultServerPort     = "8000"
	defaultComputingPower = "4"
)

type Config struct {
	ComputingPower int
	ServerPort     string
}

func LoadFromEnv() (*Config, error) {
	conf := &Config{}

	var err error

	computingPower := os.Getenv("COMPUTING_POWER")
	if computingPower == "" {
		computingPower = defaultComputingPower
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = defaultServerPort
	}

	conf.ComputingPower, err = strconv.Atoi(computingPower)
	if err != nil {
		logger.Error("Failed to parse %s as int: %v", computingPower, err)
		return nil, err

	}

	if _, err := strconv.Atoi(serverPort); err != nil {
		logger.Error("Failed to parse %s as int: %v", serverPort, err)
		return nil, err
	}
	conf.ServerPort = serverPort

	return conf, nil

}
