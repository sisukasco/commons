package utils

import (
	"strings"
)

func MatchDomain(spec string, domain string) bool {
	if spec == "" {
		return true
	}
	specDomains := strings.Split(spec, ",")
	for _, sd := range specDomains {
		subDomainsSpec := strings.Split(sd, ".")
		subDomains := strings.Split(domain, ".")
		if len(subDomainsSpec) != len(subDomains) {
			continue
		}
		matches := true

		for p := 0; p < len(subDomainsSpec); p++ {
			if subDomainsSpec[p] == "*" {
				continue
			}
			if subDomainsSpec[p] != subDomains[p] {
				matches = false
				break
			}
		}
		if matches {
			return true
		}
	}
	return false
}
