package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	TimeAddition       int
	TimeSubtraction    int
	TimeMultiplication int
	TimeDivision       int
}

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
		log.Fatalf("Failed to parse %s as int: %v", timeAddition, err)
		return nil, err

	}

	conf.TimeSubtraction, err = strconv.Atoi(timeSubtraction)
	if err != nil {
		log.Fatalf("Failed to parse %s as int: %v", timeSubtraction, err)
		return nil, err
	}

	conf.TimeMultiplication, err = strconv.Atoi(timeMultiplication)
	if err != nil {
		log.Fatalf("Failed to parse %s as int: %v", timeMultiplication, err)
		return nil, err
	}

	conf.TimeDivision, err = strconv.Atoi(timeDivision)
	if err != nil {
		log.Fatalf("Failed to parse %s as int: %v", timeDivision, err)
		return nil, err
	}

	return conf, nil

}