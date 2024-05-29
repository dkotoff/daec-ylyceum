package config

import (
	"os"
	"strconv"

	"github.com/dkotoff/daec-ylyceum/server/logger"
)

type Config struct {
	TimeAddition       int
	TimeSubtraction    int
	TimeMultiplication int
	TimeDivision       int
}

const (
	defaultTimeAddition       = "1000"
	defaultTimeSubtraction    = "1000"
	defaultTimeMultiplication = "1000"
	defaultTimeDivision       = "1000"
)

func LoadFromEnv() (*Config, error) {
	conf := &Config{}

	var err error

	timeAddition := os.Getenv("TIME_ADDITION_MS")
	if timeAddition == "" {
		timeAddition = defaultTimeAddition
	}
	timeSubtraction := os.Getenv("TIME_SUBTRACTION_MS")
	if timeSubtraction == "" {
		timeSubtraction = defaultTimeSubtraction
	}
	timeMultiplication := os.Getenv("TIME_MULTIPLICATION_MS")
	if timeMultiplication == "" {
		timeMultiplication = defaultTimeMultiplication
	}
	timeDivision := os.Getenv("TIME_DIVISION_MS")
	if timeDivision == "" {
		timeDivision = defaultTimeDivision
	}

	conf.TimeAddition, err = strconv.Atoi(timeAddition)
	if err != nil {
		logger.Fatal("Failed to parse %s as int: %v", timeAddition, err)
		return nil, err

	}

	conf.TimeSubtraction, err = strconv.Atoi(timeSubtraction)
	if err != nil {
		logger.Fatal("Failed to parse %s as int: %v", timeSubtraction, err)
		return nil, err
	}

	conf.TimeMultiplication, err = strconv.Atoi(timeMultiplication)
	if err != nil {
		logger.Fatal("Failed to parse %s as int: %v", timeMultiplication, err)
		return nil, err
	}

	conf.TimeDivision, err = strconv.Atoi(timeDivision)
	if err != nil {
		logger.Fatal("Failed to parse %s as int: %v", timeDivision, err)
		return nil, err
	}

	return conf, nil

}
