package parser

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/url"
)

type PageData struct {
	Links []string
}

func ParseHTML(body io.Reader, baseURL *url.URL) (PageData, error) {
    doc, err := goquery.NewDocumentFromReader(body)
    if err != nil {
        return PageData{}, err
    }

    var links []string
    doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
        href, _ := s.Attr("href")
        
        // Parse the link
        u, err := url.Parse(href)
        if err != nil {
            return
        }
        
        // Resolve relative to absolute
        links = append(links, baseURL.ResolveReference(u).String())
    })

    return PageData{Links: links}, nil
}