package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jessevdk/go-flags"

	"github.com/akellbl4/wod/app/config"
	"github.com/akellbl4/wod/app/domain"
	"github.com/akellbl4/wod/app/helpers"
)

// Opts with all cli commands and flags
type Opts struct {
	ConfigPath      string `long:"config" short:"c" env:"WOD_CONFIG" description:"Path to config.json" default:"./config.json"`
	TemplatePath    string `long:"template" short:"t" env:"WOD_TEMPLATE_PATH" description:"Path to html template" default:"./template.html.tmpl"`
	DestinationPath string `long:"destination" short:"d" env:"WOD_DESTINATION" description:"Path to destination directory" default:"./web"`
}

// TemplateData with all template variables
type TemplateData struct {
	Domain           string
	Price            int
	RegistrationDate time.Time
}

func main() {
	var opts Opts

	if _, err := flags.Parse(&opts); err != nil {
		log.Printf("error on reading flags: %v", err)
		os.Exit(1)
	}

	cfg, err := config.Parse(opts.ConfigPath)

	if err != nil {
		log.Printf("error on parse config: %v", err)
		os.Exit(1)
	}

	tmpl := template.Must(template.ParseFiles(opts.TemplatePath))
	err = helpers.CreateFolder(opts.DestinationPath)

	if err != nil {
		log.Printf("error on creation destination folder: %v", err)
		os.Exit(1)
	}

	err = createPages(cfg, tmpl, opts.DestinationPath)

	if err != nil {
		log.Printf("error on creation pages: %v", err)
		os.Exit(1)
	}
}

func createPages(cfg config.Config, tmpl *template.Template, destPath string) error {
	for i, item := range cfg.Domains {
		if item.Name == "" {
			return fmt.Errorf("domain name not defined for record #%d in config", i)
		}

		if item.PricePerDay == 0 {
			item.PricePerDay = cfg.PricePerDay
		}

		if item.PricePerDay == 0 {
			return fmt.Errorf("price per day for %s domain not defined for record #%d in config, global price also, not defined", item.Name, i)
		}

		if item.Rate == 0 {
			item.Rate = 1
		}

		domainInfo, err := domain.GetDomainInfo(item.Name)

		if err != nil {
			return err
		}

		creationDate, err := time.Parse(time.RFC3339, domainInfo.Domain.CreatedDate)

		if err != nil {
			return err
		}

		days := helpers.GetDaysFromDate(creationDate)
		price := days * item.PricePerDay * item.Rate
		file, err := os.Create(filepath.Join(destPath, "/", item.Name+".html"))

		if err != nil {
			return err
		}

		data := TemplateData{
			Domain:           item.Name,
			Price:            price,
			RegistrationDate: creationDate,
		}
		err = tmpl.Execute(file, data)

		if err != nil {
			return err
		}

		err = file.Close()

		if err != nil {
			return err
		}

		log.Printf("domain: %s, price: %v\n", item.Name, price)
	}

	return nil
}
