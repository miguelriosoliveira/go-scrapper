package main

import (
	"net/url"

	"scrapper/src/crawler"
)

func main() {
	startURL := "https://parserdigital.com/"
	c := crawler.NewCrawler()
	startUrlParsed, _ := url.Parse(startURL)
	c.Crawl(startURL, startUrlParsed.Host)
}
