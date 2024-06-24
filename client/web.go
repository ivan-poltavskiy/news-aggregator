package client

import (
	"fmt"
	"net/http"
	"news-aggregator/entity/article"
	"news-aggregator/filter"
	"strings"
)

type WebClient struct {
	aggregator       Aggregator
	Sources          string
	keywords         string
	startDateStr     string
	endDateStr       string
	sortBy           string
	sortingBySources bool
	help             bool
}

// NewWebClient creates and initializes a new web client with the provided aggregator.
func NewWebClient(r http.Request, aggregator Aggregator) Client {

	queryParams := r.URL.Query()

	webClient := &WebClient{aggregator: aggregator}
	webClient.Sources = queryParams.Get("sources")
	webClient.keywords = queryParams.Get("keywords")
	webClient.startDateStr = queryParams.Get("startDate")
	webClient.endDateStr = queryParams.Get("endDate")
	webClient.sortBy = queryParams.Get("sortBy")
	webClient.sortingBySources = queryParams.Get("sortingBySources") == "true"
	webClient.help = queryParams.Get("help") == "true"
	return webClient
}

// FetchArticles retrieves articles based on arguments provided as params.
func (webClient *WebClient) FetchArticles() ([]article.Article, error) {
	if webClient.help {
		webClient.printUsage()
		return nil, nil
	}
	filters, uniqueSources, fetchParametersError := webClient.fetchParameters()
	if fetchParametersError != nil {
		return nil, fetchParametersError
	}

	articles, err := webClient.aggregator.Aggregate(uniqueSources, filters...)
	if err != nil {
		return nil, err
	}

	articles, fetchParametersError = DateSorter{}.SortArticle(articles, webClient.sortBy)
	if fetchParametersError != nil {
		return nil, fetchParametersError
	}
	return articles, nil
}

func (webClient *WebClient) Print(articles []article.Article) {

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

// fetchParameters extracts and validates command line parameters,
// including sources and filters, and returns them for use in article fetching.
func (webClient *WebClient) fetchParameters() ([]filter.ArticleFilter, []string, error) {
	sourceNames := strings.Split(webClient.Sources, ",")
	var filters []filter.ArticleFilter

	filters = buildKeywordFilter(webClient.keywords, filters)
	filters, err := buildDateFilters(webClient.startDateStr, webClient.endDateStr, filters)
	if err != nil {
		return nil, nil, err
	}
	uniqueSources := checkUnique(sourceNames)
	return filters, uniqueSources, nil
}
