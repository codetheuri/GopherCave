# GopherCave

**High-performance, concurrent website archiver and offline mirror creator written in Go.**

GopherCave crawls a target website, downloads HTML pages together with all referenced assets (CSS, JavaScript, images, favicons, etc.), rewrites internal links to work offline, and saves everything in a structured local directory — giving you a fully browsable static copy of the site.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## ✨ Features

- **Recursive crawling** — follows internal links up to a configurable depth
- **Concurrent fetching** — uses goroutines + semaphore to download multiple resources safely and efficiently
- **Asset mirroring** — saves stylesheets, scripts, images and other files into `/assets/`
- **Link rewriting** — converts absolute, root-relative and protocol-relative URLs so the offline copy navigates correctly
- **Built-in preview server** — instantly view your archived site in the browser
- **Metadata extraction** — saves structural elements (headers, footers, cards, etc.) as clean JSON
- **Polite crawling** — configurable delay between requests to avoid overwhelming servers

## 📂 Project Layout

```text
GopherCave/
├── cmd/
│   └── scraper/           # CLI entry point
├── internal/
│   ├── crawler/           # Recursion logic & concurrency control
│   ├── fetcher/           # HTTP client with UA, timeouts, redirects
│   ├── parser/            # goquery-based HTML parsing & link rewriting
│   ├── saver/             # Filesystem layout, directory creation, asset storage
│   └── server/            # Simple static HTTP preview server
├── go.mod
└── README.md
```


## 📥 Installation
# 1. Clone the repository
git clone https://github.com/codetheuri/GopherCave.git
cd GopherCave

## 2. Download dependencies

```text
go mod tidy
```

## 🚀 Quick Start
# Basic usage — crawl and preview
```
go run ./cmd/scraper https://example.com/blog/
```

# After crawling finishes, the preview server starts automatically.
# Open in browser: http://localhost:8080
## ⚙️ Configuration (for now)
All important limits are currently defined as constants in internal/crawler/crawler.go. Edit and recompile to change them:
```
const (
    MaxDepth         = 2
    MaxConcurrency   = 5
    PolitenessDelay  = 100 * time.Millisecond
    // ...
)
```
## 🛠️ Example Output Structure
After running on https://example.com/:
 ```
 output/
├── index.html
├── blog/
│   ├── post-1/
│   │   └── index.html
│   └── post-2/
│       └── index.html
└── assets/
    ├── css/
    │   └── style.min.css
    ├── js/
    │   └── main.bundle.js
    ├── images/
    │   ├── logo.svg
    │   └── hero-01.jpg
    └── favicon.ico
    ```

