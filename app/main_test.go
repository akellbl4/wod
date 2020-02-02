package main

import (
	"html/template"
	"io/ioutil"
	"os"
	"testing"

	"github.com/akellbl4/wod/app/config"
	"github.com/stretchr/testify/assert"
)

func Test_createPages(t *testing.T) {
	var testSet = []struct {
		name   string
		config config.Config
		tmpl   *template.Template
		errMsg string
	}{
		{
			name: "with empty domain name error",
			config: config.Config{
				Domains: []config.Domain{
					{
						Name: "",
					},
				},
				PricePerDay: 1,
			},
			tmpl: &template.Template{},
			errMsg: "domain name not defined for record #0 in config",
		},
		{
			name: "with empty price per day error",
			config: config.Config{
				Domains: []config.Domain{
					{
						Name: "example.com",
					},
				},
			},
			tmpl: &template.Template{},
			errMsg: "price per day for example.com domain not defined for record #0 in config, global price also, not defined",
		},
		{
			name: "with bad domain name",
			config: config.Config{
				Domains: []config.Domain{
					{
						Name: "asdasd",
					},
				},
				PricePerDay: 1,
			},
			tmpl: &template.Template{},
			errMsg: "Domain whois data invalid.",
		},
		{
			name: "successful flow",
			config: config.Config{
				Domains: []config.Domain{
					{
						Name: "google.com",
					},
				},
				PricePerDay: 1,
			},
			tmpl: template.Must(template.New("template").Parse("{{.Domain}},{{.Price}},{{.RegistrationDate}}")),
		},
	}

	dir, err := ioutil.TempDir(os.TempDir(), "wod")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	for _, item := range testSet {
		t.Run(item.name, func(t *testing.T) {
			err := createPages(item.config, item.tmpl, dir)
			if item.errMsg != "" {
				assert.Error(t, err, item.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

