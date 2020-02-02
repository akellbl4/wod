package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Domain struct is data about domain what we neet to evaluate
type Domain struct {
	Name        string `json:"name"`
	PricePerDay int    `json:"price_per_day,omitempty"`
	Rate        int    `json:"rate,omitempty"`
}

// Config struct pis data about domais
// and common values for domains
// without overrided PricePerDay and Rate
type Config struct {
	Domains     []Domain `json:"domains"`
	PricePerDay int      `json:"price_per_day"`
	Rate        int      `json:"rate"`
}

// Parse returns content of file
// which was readed by provided filepath
func Parse(filepath string) (config Config, err error) {
	log.Printf("Read config file: %q\n", filepath)
	configBytes, err := ioutil.ReadFile(filepath)

	if err != nil {
		return Config{}, err
	}

	err = json.Unmarshal(configBytes, &config)

	if err != nil {
		return Config{}, err
	}

	return config, nil
}
