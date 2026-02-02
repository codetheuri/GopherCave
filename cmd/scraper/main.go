package main

import (
	"fmt"
	"os"
	"scraper/internal/crawler"
	"scraper/internal/server"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: scraper <url>")
		return
	}

	url := os.Args[1]
	
	fmt.Println("Starting crawl...")
	crawler.StartCrawl(url)
	
	fmt.Println("Crawl complete.")
	server.Serve(url, "8080")
}