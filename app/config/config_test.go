package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	var testSet = []struct{
		name   string
		path   string
		errMsg string
		result Config
	} {
		{
			name: "failed to read config.json",
			path: "config.json",
			errMsg: "open config.json: no such file or directory",
		},
		{
			name: "failed to parse config.json",
			path: "testdata/config_wrong.json",
			errMsg: "invalid character '\n' in string literal",
		},
		{
			name: "no error expected to valid config.json",
			path: "testdata/config.json",
			errMsg: "",
			result: Config{
				Domains: []Domain{
					{Name: "google.com", Rate: 10},
				},
				PricePerDay: 1,
				Rate:        1,
			},
		},
	}

	for _, tt := range testSet {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.path)

			if tt.errMsg != "" {
				assert.Error(t, err, tt.errMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.result, result)
			}
		})
	}
}
