package crawler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/schollz/progressbar/v3"
	"scraper/internal/extractor"
	"scraper/internal/fetcher"
	"scraper/internal/parser"
	"scraper/internal/saver"
)

var (
	visited    = make(map[string]bool)
	allResults []string
	visitedMux sync.Mutex
	wg         sync.WaitGroup
	maxWorkers = make(chan struct{}, 5)
	bar        = progressbar.Default(-1, "Scraping")
)

func Crawl(targetURL string, depth int) {
	if depth > 2 { return }

	visitedMux.Lock()
	if visited[targetURL] {
		visitedMux.Unlock()
		return
	}
	visited[targetURL] = true
	allResults = append(allResults, targetURL)
	visitedMux.Unlock()

	wg.Add(1)
	go func(tURL string, d int) {
		defer wg.Done()
		maxWorkers <- struct{}{}
		defer func() { <-maxWorkers }()

		time.Sleep(100 * time.Millisecond)
		bar.Add(1)

		resp, err := fetcher.Fetch(tURL)
		if err != nil { return }
		defer resp.Body.Close()

		bodyBytes, _ := io.ReadAll(resp.Body)
		u, _ := url.Parse(tURL)
		doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(bodyBytes))

		// Process Assets
		doc.Find("img, link, script").Each(func(i int, s *goquery.Selection) {
			attr := "src"
			if s.Is("link") { attr = "href" }
			val := s.AttrOr(attr, "")
			if val == "" { return }

			assetURL, _ := u.Parse(val)
			fileName := filepath.Base(assetURL.Path)
			
			wg.Add(1)
			go func(aURL string) {
				defer wg.Done()
				saver.DownloadAsset(filepath.Join("output", u.Host), aURL)
			}(assetURL.String())
			
			s.SetAttr(attr, "/assets/"+fileName)
		})

		// Fix Links for Sidebar
		doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			if href != "" && !strings.HasPrefix(href, "http") && !strings.HasPrefix(href, "#") {
				lURL, _ := u.Parse(href)
				s.SetAttr("href", filepath.Join("/", lURL.Path, "index.html"))
			}
		})

		htmlStr, _ := doc.Html()
		
		// CRITICAL FIX: Clean the path and remove leading slashes so it stays in ./output/
		relativeDir := strings.TrimPrefix(filepath.Clean(u.Path), "/")
		if relativeDir == "" || relativeDir == "." {
			relativeDir = "." 
		}
		
		targetDir := filepath.Join("output", u.Host, relativeDir)
		targetFile := filepath.Join(targetDir, "index.html")
		
		// Save HTML
		saver.SavePage(targetFile, strings.NewReader(htmlStr))

		// Save Metadata
		comps := extractor.ExtractComponents(doc)
		jsonData, _ := json.MarshalIndent(comps, "", "  ")
		os.MkdirAll(targetDir, 0755)
		os.WriteFile(filepath.Join(targetDir, "components.json"), jsonData, 0644)

		// Crawl deeper
		pageData, _ := parser.ParseHTML(bytes.NewReader(bodyBytes), u)
		for _, link := range pageData.Links {
			lURL, err := url.Parse(link)
			if err == nil && lURL.Host == u.Host {
				Crawl(link, d+1)
			}
		}
	}(targetURL, depth)
}

func StartCrawl(targetURL string) {
	Crawl(targetURL, 0)
	wg.Wait()
	bar.Finish()
	createMasterIndex(targetURL)
}

func createMasterIndex(startURL string) {
	u, _ := url.Parse(startURL)
	baseDir := filepath.Join("output", u.Host)
	os.MkdirAll(baseDir, 0755)
	
	var html bytes.Buffer
	html.WriteString("<html><body style='background:#111; color:#eee; font-family:sans-serif; padding:50px;'><h1>Sitemap</h1><ul>")
	for _, link := range allResults {
		p, _ := url.Parse(link)
		linkPath := filepath.Join("/", p.Path, "index.html")
		html.WriteString(fmt.Sprintf(`<li><a style='color:#4fb' href="%s">%s</a></li>`, linkPath, link))
	}
	html.WriteString("</ul></body></html>")
	os.WriteFile(filepath.Join(baseDir, "index.html"), html.Bytes(), 0644)
}