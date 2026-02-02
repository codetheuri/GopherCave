package server

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
)

func Serve(targetURL, port string) {
	u, _ := url.Parse(targetURL)
	rootPath, _ := filepath.Abs(filepath.Join("output", u.Host))
	
	// Create a handler that serves the rootPath
	fs := http.FileServer(http.Dir(rootPath))
	
	fmt.Printf("\n--- Preview Server ---\n")
	fmt.Printf("Viewing results at: http://localhost:%s/\n", port)
	
	// Start server
	err := http.ListenAndServe(":"+port, fs)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}