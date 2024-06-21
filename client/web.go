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

func NewWebClient(r http.Request, aggregator Aggregator) Client {

	queryParams := r.URL.Query()

	webClient := &WebClient{aggregator: aggregator}
	webClient.Sources = queryParams.Get("sources")
	webClient.keywords = queryParams.Get("keywords")
	webClient.startDateStr = queryParams.Get("startDate")
	webClient.endDateStr = queryParams.Get("endDate")
	webClient.sortBy = queryParams.Get("sortBy")
	webClient.sortingBySources = queryParams.Get("sortingBySources") == "true"
	return webClient
}

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

//// buildKeywordFilter extracts keywords from command line arguments and adds them to the filters.
//func buildKeywordFilter(cli *WebClient, filters []filter.ArticleFilter) []filter.ArticleFilter {
//	if cli.keywords != "" {
//		keywords := strings.Split(cli.keywords, ",")
//		uniqueKeywords := checkUnique(keywords)
//		filters = append(filters, filter.ByKeyword{Keywords: uniqueKeywords})
//	}
//	return filters
//}

//// buildDateFilters extracts date filters from command line arguments and adds them to the filters.
//func buildDateFilters(cli *WebClient, filters []filter.ArticleFilter) ([]filter.ArticleFilter, error) {
//
//	validationErr, isValid := validator.ValidateDate(cli.startDateStr, cli.endDateStr)
//
//	if validationErr != nil {
//		return nil, validationErr
//	}
//	if isValid {
//
//		startDate, err := time.Parse(constant.DateOutputLayout, cli.startDateStr)
//
//		if err != nil {
//			return nil, errors.New("Invalid start date: " + cli.startDateStr)
//		}
//
//		endDate, err := time.Parse(constant.DateOutputLayout, cli.endDateStr)
//
//		if err != nil {
//			return nil, errors.New("Invalid end date: " + cli.endDateStr)
//		}
//
//		return append(filters, filter.ByDate{StartDate: startDate, EndDate: endDate}), nil
//	}
//	return filters, nil
//}
//
//// CheckUnique returns a slice containing only unique strings from the input slice.
//func checkUnique(input []string) []string {
//	uniqueMap := make(map[string]struct{})
//	var uniqueList []string
//	for _, item := range input {
//		if _, ok := uniqueMap[item]; !ok {
//			uniqueMap[item] = struct{}{}
//			uniqueList = append(uniqueList, item)
//		}
//	}
//	return uniqueList
//}
