package client

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"news-aggregator/entity/article"
	"news-aggregator/filter"
	"strings"
)

type WebClient struct {
	aggregator       Aggregator
	Sources          []string
	sortBy           string
	sortingBySources bool
	help             bool
	DateSorter       DateSorter
	filters          []filter.ArticleFilter
	output           http.ResponseWriter
}

// NewWebClient creates and initializes a new web client with the provided aggregator.
func NewWebClient(r http.Request, w http.ResponseWriter, aggregator Aggregator) Client {

	queryParams := r.URL.Query()
	webClient := &WebClient{aggregator: aggregator}
	webClient.Sources = checkUnique(strings.Split(queryParams.Get("sources"), ","))
	webClient.sortBy = queryParams.Get("sortBy")
	webClient.sortingBySources = queryParams.Get("sortingBySources") == "true"
	webClient.help = queryParams.Get("help") == "true"
	webClient.DateSorter = DateSorter{}
	webClient.filters = buildKeywordFilter(queryParams.Get("keywords"), webClient.filters)
	filters, err := buildDateFilters(queryParams.Get("startDate"), queryParams.Get("endDate"), webClient.filters)
	if err != nil {
		fmt.Println(err)
	}
	webClient.filters = filters
	webClient.output = w
	logrus.Info("New web client initialized")
	return webClient
}

// FetchArticles retrieves articles based on arguments provided as params.
func (webClient *WebClient) FetchArticles() ([]article.Article, error) {
	if webClient.help {
		webClient.printUsage()
		return nil, nil
	}

	articles, err := webClient.aggregator.Aggregate(webClient.Sources, webClient.filters...)
	if err != nil {
		return nil, err
	}
	logrus.Info("Web client: articles aggregate successfully. Length: ", len(articles))

	articles, fetchParametersError := webClient.DateSorter.SortArticle(articles, webClient.sortBy)
	if fetchParametersError != nil {
		return nil, fetchParametersError
	}
	return articles, nil
}

func (webClient *WebClient) Print(articles []article.Article) {
	webClient.output.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(webClient.output).Encode(articles)
	if err != nil {
		logrus.Error("Failed to encode json: ", err)
		http.Error(webClient.output, "Failed to encode json: "+err.Error(), http.StatusInternalServerError)
	}
}

// printUsage prints the usage instructions
func (webClient *WebClient) printUsage() {
	fmt.Println("Usage of news-aggregator:" +
		"\nType --sources, and then list the resources you want to retrieve information from. " +
		"The program supports such news resources:\nABC, BBC, NBC, USA Today and Washington Times. \n" +
		"\nType --keywords, and then list the keywords by which you want to filter articles. \n" +
		"\nType --startDate and --endDate to filter by date. News published between the specified dates will be shown." +
		"Date format - yyyy-mm-dd" + "" +
		"Type --sortBy to sort by DESC/ASC." + "Type --sortingBySources to sort by sources.")
}
