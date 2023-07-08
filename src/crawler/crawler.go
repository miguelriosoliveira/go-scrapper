package crawler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"scrapper/src/utils"
)

type Crawler struct {
	visited map[string]bool
	count   int
}

func NewCrawler() *Crawler {
	return &Crawler{
		visited: make(map[string]bool),
		count:   0,
	}
}

func (c *Crawler) Crawl(urlStr, baseDomain string) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		log.Printf("Error parsing URL: %v", err)
		return
	}

	if c.visited[parsedURL.String()] || !utils.IsSameDomain(parsedURL, baseDomain) {
		return
	}

	c.visited[parsedURL.String()] = true

	resp, err := http.Get(urlStr)
	if err != nil {
		log.Printf("Error fetching URL: %v", err)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Error loading HTML document: %v", err)
		return
	}

	fmt.Println("Visited:", parsedURL)

	links := getLinks(doc, parsedURL)
	if len(links) == 0 {
		fmt.Println("No links found!")
		return
	}
	fmt.Printf("Links found:\n- %s\n", strings.Join(links, "\n- "))
	for _, link := range links {
		c.Crawl(link, baseDomain)
	}
}

func getLinks(doc *goquery.Document, baseURL *url.URL) []string {
	control := make(map[string]bool)
	links := make([]string, 0)
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")

		if !exists {
			return;
		}
		absURL := utils.ResolveURL(baseURL, href)
		if absURL == nil || !utils.IsSameDomain(absURL, baseURL.Host) || utils.IsSectionLink(href) || control[absURL.String()] {
			return;
		}
		control[absURL.String()] = true
		links = append(links, absURL.String())
	})
	return links
}
