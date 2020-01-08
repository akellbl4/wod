package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/domainr/whois"
	"github.com/jessevdk/go-flags"
)

// Options with all cli commands and flags
type Options struct {
	Path     string `long:"config" short:"c" env:"WOD_CONFIG" description:"Path to config.json" default:"./config.json"`
	Dest     string `long:"destination" short:"d" env:"WOD_DESTINATION" description:"Path to destination directory" default:"./web"`
	Template string `long:"template" short:"t" env:"WOD_TEMPLATE_PATH" description:"Path to html template" default:"./template.html"`
}

// Domain struct
type Domain struct {
	Name        string `json:"name"`
	PricePerDay int    `json:"price_per_day"`
	Rate        int    `json:"rate"`
}

// Config struct
type Config struct {
	Domains     []Domain `json:"domains"`
	PricePerDay int      `json:"price_per_day"`
	Rate        int      `json:"rate"`
}

// TemplateData struct
type TemplateData struct {
	Domain          string
	Price           int
	RegistratedDate string
}

func main() {
	var opts Options

	p := flags.NewParser(&opts, flags.Default)

	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	fmt.Printf("Read config file: '%s'\n", opts.Path)
	configBytes, err := ioutil.ReadFile(opts.Path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tmpl := template.Must(template.ParseFiles(opts.Template))

	var config Config
	json.Unmarshal(configBytes, &config)

	if _, err = os.Stat(opts.Dest); os.IsNotExist(err) {
		os.Mkdir(opts.Dest, os.ModePerm)
	}

	for i := 0; i < len(config.Domains); i++ {
		name := config.Domains[i].Name
		if name == "" {
			fmt.Println("Domain name not defined")
			os.Exit(1)
		}
		pricePerDay := config.Domains[i].PricePerDay
		if pricePerDay == 0 {
			pricePerDay = config.PricePerDay
		}
		if pricePerDay == 0 {
			fmt.Printf("Price per day for %s domain not defined, global price also, not defined\n", name)
			os.Exit(1)
		}
		rate := config.Domains[i].Rate
		if rate == 0 {
			rate = 1
		}
		domainInfo := getDomainInfo(name)
		creationDate, _ := getDomainCreationDate(domainInfo)
		daysFromCreation := getDaysFromCreation(creationDate)
		price := daysFromCreation * 3 * 10
		data := TemplateData{name, price, creationDate.Format("02-Jan-2006")}
		file, err := os.Create(filepath.Join(opts.Dest, "/", name+".html"))
		tmpl.Execute(file, data)
		if err != nil {
			fmt.Print("execute: ", err)
			return
		}
		file.Close()
		fmt.Printf("Domain: %s, price: %v\n", name, price)
	}

}

func getDomainInfo(domain string) string {
	request, _ := whois.NewRequest(domain)
	response, _ := whois.DefaultClient.Fetch(request)
	responsebody := string(response.Body)

	return responsebody
}

func getDomainCreationDate(domainInfo string) (creationDate time.Time, err error) {
	reMatch := regexp.MustCompile(`(?:Creation Date|created):\s+([\w\d:-]*)`)
	match := reMatch.FindStringSubmatch(domainInfo)

	if len(match) == 0 {
		return time.Time{}, errors.New("can't find creation date")
	}

	creationDate, err = time.Parse(time.RFC3339, match[1])

	if err != nil {
		return time.Time{}, errors.New("can't parse creation date")
	}

	return creationDate, nil
}

func getDaysFromCreation(creationDate time.Time) int {
	now := time.Now()
	days := int(math.Round(now.Sub(creationDate).Hours() / 24))

	return days
}
