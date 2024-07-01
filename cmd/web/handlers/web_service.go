package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"news-aggregator/constant"
	"news-aggregator/entity/source"
	"os"
	"regexp"
	"strings"
)

func WriteSourcesToFile(sources []source.Source) error {
	file, err := os.Create(constant.PathToStorage)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Print("Error closing file: ", err)
		}
	}(file)

	encoder := json.NewEncoder(file)
	err = encoder.Encode(&sources)
	if err != nil {
		return err
	}

	return nil
}

func ReadSourcesFromFile() []source.Source {
	file, err := os.Open(constant.PathToStorage)
	if err != nil {
		if os.IsNotExist(err) {
			return []source.Source{}
		}
		fmt.Println("Error opening sources file:", err)
		return nil
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Print("Error closing file: ", err)
		}
	}(file)

	var sources []source.Source
	if err := json.NewDecoder(file).Decode(&sources); err != nil {
		fmt.Println("Error decoding sources file:", err)
		return nil
	}
	return sources
}

func GetRssFeedLink(w http.ResponseWriter, url string) (error, string) {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to download page", http.StatusInternalServerError)
		return err, ""
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Print("Error closing file: ", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read page content", http.StatusInternalServerError)
		return err, ""
	}

	re := regexp.MustCompile(`(?i)<link[^>]+type="application/rss\+xml"[^>]+href="([^"]+)"`)
	matches := re.FindStringSubmatch(string(body))

	if len(matches) < 2 {
		http.Error(w, "RSS link not found", http.StatusBadRequest)
		return nil, ""
	}

	rssURL := matches[1]
	return nil, rssURL
}

func ExtractDomainName(url string) string {
	re := regexp.MustCompile(`https?://(www\.)?([^/]+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 3 {
		return "unknown"
	}
	domain := matches[2]
	domain = strings.Split(domain, ".")[0]
	return domain
}

func IsSourceExists(name source.Name) bool {
	sources := ReadSourcesFromFile()
	for _, s := range sources {
		if s.Name == name {
			return true
		}
	}
	return false
}

func AddSourceToStorage(newSource source.Source) {
	sources := append(ReadSourcesFromFile(), newSource)

	file, err := os.Create(constant.PathToStorage)
	if err != nil {
		fmt.Println("Error creating sources file:", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Print("Error closing file: ", err)
		}
	}(file)

	if err := json.NewEncoder(file).Encode(sources); err != nil {
		fmt.Println("Error encoding sources to file:", err)
	}
}
