package client

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"news-aggregator/entity/news"
	"news-aggregator/filter"
	"news-aggregator/sorter"
	"news-aggregator/storage/source"
	"strings"
)

type WebClient struct {
	aggregator       Aggregator
	Sources          []string
	sortBy           string
	sortingBySources bool
	help             bool
	DateSorter       sorter.DateSorter
	filters          []filter.NewsFilter
	output           http.ResponseWriter
	sourceStorage    source.Storage
}

// NewWebClient creates and initializes a new web client with the provided aggregator.
func NewWebClient(r http.Request, w http.ResponseWriter, aggregator Aggregator, sourceStorage source.Storage) Client {
	queryParams := r.URL.Query()
	webClient := &WebClient{aggregator: aggregator}
	webClient.Sources = checkUnique(strings.Split(queryParams.Get("sources"), ","))
	webClient.sortBy = queryParams.Get("sortBy")
	webClient.sortingBySources = queryParams.Get("sortingBySources") == "true"
	webClient.help = queryParams.Get("help") == "true"
	webClient.DateSorter = sorter.DateSorter{}
	webClient.filters = buildKeywordFilter(queryParams.Get("keywords"), webClient.filters)
	filters, err := buildDateFilters(queryParams.Get("startDate"), queryParams.Get("endDate"), webClient.filters)
	if err != nil {
		logrus.Error("New web client initialization error: ", err)
	}
	webClient.filters = filters
	webClient.output = w
	webClient.sourceStorage = sourceStorage
	logrus.Info("New web client initialized")
	return webClient
}

// FetchNews retrieves articles based on arguments provided as params.
func (webClient *WebClient) FetchNews() ([]news.News, error) {
	if webClient.help {
		webClient.printUsage()
		return nil, nil
	}

	articles, err := webClient.aggregator.Aggregate(webClient.Sources, webClient.filters...)
	if err != nil {
		return nil, err
	}
	logrus.Info("Web client: articles aggregate successfully. Length: ", len(articles))

	articles, fetchParametersError := webClient.DateSorter.SortNews(articles, webClient.sortBy)
	if fetchParametersError != nil {
		return nil, fetchParametersError
	}
	return articles, nil
}

func (webClient *WebClient) Print(news []news.News) {
	webClient.output.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(webClient.output).Encode(news)
	if err != nil {
		logrus.Error("Failed to encode json: ", err)
		http.Error(webClient.output, "Failed to encode json: "+err.Error(), http.StatusInternalServerError)
	}
}

// printUsage prints the usage instructions
func (webClient *WebClient) printUsage() {
	webClient.output.Header().Set("Content-Type", "text/plain")
	_, err := fmt.Fprintln(webClient.output, "Usage of news-aggregator:"+
		"\nType --sources, and then list the news you want to retrieve information from. "+
		"The program supports such news news:\nABC, BBC, NBC, USA Today and Washington Times. \n"+
		"\nType --keywords, and then list the keywords by which you want to filter articles. \n"+
		"\nType --startDate and --endDate to filter by date. News published between the specified dates will be shown."+
		"Date format - yyyy-mm-dd"+
		"\nType --sortBy to sort by DESC/ASC."+
		"\nType --sortingBySources to sort by sources.")
	if err != nil {
		return
	}
}
