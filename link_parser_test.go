package main

import (
	"net/url"
	"testing"
)

func TestBasicLinkParsing(t *testing.T) {
	html := `
        <html>... <a href="http://example.com/page"> ...</html>
    `
	expected, _ := url.Parse("http://example.com/page")
	currentURL, _ := url.Parse("http://example.com")

	links := ExtractLinks(html, *currentURL)

	if len(links) != 1 {
		t.Error("Expected 1 link, got", len(links))
	}

	if links[0].String() != expected.String() {
		t.Error("Unexpected link", links[0].String())
	}
}

func TestSrcLinkParsing(t *testing.T) {
	html := `
        <html>... <img src="http://example.com/image.jpg"> ...</html>
    `
	expected, _ := url.Parse("http://example.com/image.jpg")
	currentURL, _ := url.Parse("http://example.com")

	links := ExtractLinks(html, *currentURL)

	if links[0].String() != expected.String() {
		t.Error("Unexpected result", links[0].String())
	}
}

func TestIgnoresMailtoLink(t *testing.T) {
	html := `
        <html>... <a mailto="name@example.com"> <a href="mailto:name@example.com"> ...</html>
    `
	currentURL, _ := url.Parse("http://example.com")

	links := ExtractLinks(html, *currentURL)

	if len(links) > 0 {
		t.Error("Expected 0 links, got", len(links))
	}
}

func TestIgnoresSomeRandomProtocol(t *testing.T) {
	html := `
        <html>... <a href="somethingelse:things"> ...</html>
    `
	currentURL, _ := url.Parse("http://example.com")

	links := ExtractLinks(html, *currentURL)

	if len(links) > 0 {
		t.Error("Expected 0 links, got", len(links))
	}
}

func TestIgnoresFragment(t *testing.T) {
	html := `
        <html>... <a mailto="#anchor"> ...</html>
    `
	currentURL, _ := url.Parse("http://example.com")

	links := ExtractLinks(html, *currentURL)

	if len(links) > 0 {
		t.Error("Expected 0 links, got", len(links))
	}
}

func TestExpandsRelativeLinks(t *testing.T) {
	html := `
        <html>... <a href="page"> ...</html>
    `
	currentURL, _ := url.Parse("http://example.com")
	expected, _ := url.Parse("http://example.com/page")

	links := ExtractLinks(html, *currentURL)

	if links[0].String() != expected.String() {
		t.Error("Unexpected result", links[0].String())
	}
}

func TestExpandsAbsoluteLinks(t *testing.T) {
	html := `
        <html>... <a href="/page"> ...</html>
    `
	currentURL, _ := url.Parse("http://example.com/some/nested/page")
	expected, _ := url.Parse("http://example.com/page")

	links := ExtractLinks(html, *currentURL)

	if links[0].String() != expected.String() {
		t.Error("Unexpected result", links[0].String())
	}
}
