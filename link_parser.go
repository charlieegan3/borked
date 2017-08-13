package main

import (
	"fmt"
	"net/url"
	"regexp"
)

// ExtractLinks returns a slice of parsed URLs found in a given HTML document
// internal links are expanded to full URLs
func ExtractLinks(html string, currentURL url.URL) []url.URL {
	linkPattern := regexp.MustCompile("(href|src)=\"(\\S+)\"")
	matches := linkPattern.FindAllStringSubmatch(html, -1)

	var links []string

	for _, group := range matches {
		links = append(links, group[2])
	}

	var upgradedLinks []url.URL

	for _, link := range links {
		url, err := url.Parse(link)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if len(link) == 0 {
			continue
		}

		re := regexp.MustCompile(`^\w+:`)
		if url.Scheme != "http" && url.Scheme != "https" && re.MatchString(url.String()) {
			continue
		}

		if len(url.Fragment) > 0 && url.String()[0] == '#' {
			continue
		}

		url.Fragment = ""

		if url.Host == "" && len(url.Path) > 0 && url.Path[0] != '/' {
			if currentURL.Path == "" || currentURL.Path[len(currentURL.Path)-1] != '/' {
				url.Path = "/" + url.Path
			}
			url.Path = currentURL.Path + url.Path
		}

		if url.Scheme == "" {
			url.Scheme = currentURL.Scheme
		}

		if url.Host == "" {
			url.Host = currentURL.Host
		}

		upgradedLinks = append(upgradedLinks, *url)
	}

	return upgradedLinks
}
