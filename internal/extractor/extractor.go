package extractor

import (
	"github.com/PuerkitoBio/goquery"
)

type Components struct {
	Header string
	Footer string
	Cards  []string
}

func ExtractComponents(doc *goquery.Document) Components {
	var comps Components

	// Header
	if html, _ := doc.Find("header").Html(); html != "" {
		comps.Header = html
	}

	// Footer
	if html, _ := doc.Find("footer").Html(); html != "" {
		comps.Footer = html
	}

	// Cards
	doc.Find(".card").Each(func(i int, s *goquery.Selection) {
		h, _ := s.Html()
		comps.Cards = append(comps.Cards, h)
	})

	return comps
}
