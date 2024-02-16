package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	baseURL   = "https://online-prospekt.com/pdf/"
	outputDir = "results/"
)

var companies = []string{"lidl", "aldi-nord", "aldi-sued",
	"mueller", "real", "netto", "netto-marken-discount", "edeka"}

func fetchPDF(baseURL string, companyName string, wg *sync.WaitGroup) {
	defer wg.Done()

	start := time.Now()

	url := baseURL + companyName + "-prospekt.pdf"

	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching %s: %v\n", url, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Error: %s returned status %d\n", url, response.StatusCode)
		return
	}

	filePath := filepath.Join(outputDir, companyName+".pdf")
	savePDF(response.Body, filePath)

	fmt.Printf("PDF fetched successfully: %s in %s\n", url, time.Since(start))

}

func savePDF(body io.ReadCloser, filePath string) {

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	if _, err := io.Copy(file, body); err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}
}

func main() {
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(len(companies))

	for _, company := range companies {
		go fetchPDF(baseURL, company, &wg)
	}

	wg.Wait()

	fmt.Printf("Total time taken: %s\n", time.Since(start))
}
