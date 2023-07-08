package utils

import (
	"net/url"
	"strings"
)

func IsSameDomain(u *url.URL, domain string) bool {
	return strings.HasSuffix(u.Host, domain)
}

func IsSectionLink(href string) bool {
	return strings.HasPrefix(href, "#")
}

func ResolveURL(baseURL *url.URL, href string) *url.URL {
	relURL, err := url.Parse(href)
	if err != nil {
		return nil
	}
	return baseURL.ResolveReference(relURL)
}
