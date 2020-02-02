package domain

import (
	"math"
	"time"

	"github.com/domainr/whois"
	whoisparser "github.com/likexian/whois-parser-go"
)

// GetDomainInfo makes request for domain info
func GetDomainInfo(domain string) (whoisparser.WhoisInfo, error) {
	request, err := whois.NewRequest(domain)

	if err != nil {
		return whoisparser.WhoisInfo{}, err
	}

	response, err := whois.DefaultClient.Fetch(request)

	if err != nil {
		return whoisparser.WhoisInfo{}, err
	}

	domainInfo, err := whoisparser.Parse(string(response.Body))

	if err != nil {
		return whoisparser.WhoisInfo{}, err
	}

	return domainInfo, nil
}

// GetDaysFromCreation get days from creation date
func GetDaysFromCreation(creationDate time.Time) int {
	now := time.Now()
	days := int(math.Ceil(now.Sub(creationDate).Hours() / 24))

	return days
}
