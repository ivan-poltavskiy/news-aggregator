package main

import (
	"NewsAggregator/entity/source"
	. "NewsAggregator/initialization-data"
	. "NewsAggregator/service"
	"flag"
	"fmt"
	"strings"
	"time"
)

func main() {
	var sources string
	var keywords string
	var startDateStr string
	var endDateStr string

	flag.StringVar(&sources, "sources", "", "Specify news sources separated by comma")
	flag.StringVar(&keywords, "keywords", "", "Specify keywords to filter news articles")
	flag.StringVar(&startDateStr, "startDate", "", "Specify start date (YYYY-MM-DD)")
	flag.StringVar(&endDateStr, "endDate", "", "Specify end date (YYYY-MM-DD)")

	flag.Parse()

	if sources == "" {
		fmt.Println("Please specify at least one news source using --sources flag")
		return
	}

	InitializeSource()

	sourceNames := strings.Split(sources, ",")
	uniqueNames := make(map[string]struct{})
	var sourceNameObjects []source.Name

	for _, name := range sourceNames {
		if _, ok := uniqueNames[name]; !ok {
			uniqueNames[name] = struct{}{}
			sourceNameObjects = append(sourceNameObjects, source.Name(name))
		}
	}

	articles, errorMessage := FindNews(sourceNameObjects)

	if errorMessage != "" {
		fmt.Println(errorMessage)
	}

	if keywords != "" {
		articles = FilterNewsByKeyword(keywords, articles)
	}

	var startDate, endDate time.Time
	var err error
	if startDateStr != "" || endDateStr != "" {
		if startDateStr == "" || endDateStr == "" {
			fmt.Println("Please specify both start date and end date or omit them")
			return
		}
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			fmt.Println("Error parsing start date:", err)
			return
		}
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			fmt.Println("Error parsing end date:", err)
			return
		}

		articles = FilterByDate(startDate, endDate, articles)
	}

	for _, article := range articles {
		fmt.Println("\n")
		fmt.Println("Title of article:", article.Title)
		fmt.Println("Description:", article.Description)
		fmt.Println("Link:", article.Link)
		fmt.Println("Date:", article.Date)

	}
}
