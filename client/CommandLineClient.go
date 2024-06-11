package client

import (
	"NewsAggregator/entity/article"
	"NewsAggregator/filter"
	"flag"
	"fmt"
	"strings"
)

// CommandLineClient represents a command line client for the news-aggregator application.
type CommandLineClient struct {
	aggregator   Aggregator
	sources      string
	keywords     string
	startDateStr string
	endDateStr   string
	help         bool
}

// NewCommandLine creates and initializes a new CommandLineClient with the provided aggregator.
func NewCommandLine(aggregator Aggregator) *CommandLineClient {
	cli := &CommandLineClient{aggregator: aggregator}
	flag.StringVar(&cli.sources, "sources", "", "Specify news sources separated by comma")
	flag.StringVar(&cli.keywords, "keywords", "", "Specify keywords to filter collector articles")
	flag.StringVar(&cli.startDateStr, "startDate", "", "Specify start date (YYYY-MM-DD)")
	flag.StringVar(&cli.endDateStr, "endDate", "", "Specify end date (YYYY-MM-DD)")
	flag.BoolVar(&cli.help, "help", false, "Show help information")
	flag.Parse()
	return cli
}

// printUsage prints the usage instructions
func (cli *CommandLineClient) printUsage() {
	fmt.Println("Usage of news-aggregator:" +
		"\nType --sources, and then list the resources you want to retrieve information from. " +
		"The program supports such news resources:\nABC, BBC, NBC, USA Today and Washington Times. \n" +
		"\nType --keywords, and then list the keywords by which you want to filter articles. \n" +
		"\nType --startDate and --endDate to filter by date. News published between the specified dates will be shown." +
		"Date format - yyyy-mm-dd")
}

// FetchArticles fetches articles based on the command line arguments.
func (cli *CommandLineClient) FetchArticles() []article.Article {
	if cli.help {
		cli.printUsage()
		return nil
	}

	filters, uniqueSources := fetchParameters(cli)

	articles, errorMessage := cli.aggregator.Aggregate(uniqueSources, filters...)
	if errorMessage != "" {
		fmt.Println(errorMessage)
	}

	return articles
}

// fetchParameters extracts and validates command line parameters,
// including sources and filters, and returns them for use in article fetching.
func fetchParameters(cli *CommandLineClient) ([]filter.ArticleFilter, []string) {
	sourceNames := strings.Split(cli.sources, ",")
	var filters []filter.ArticleFilter

	filters = fetchKeywords(cli, filters)
	filters = fetchDateFilters(cli, filters)
	uniqueSources := CheckUnique(sourceNames)
	return filters, uniqueSources
}

// fetchKeywords extracts keywords from command line arguments and adds them to the filters.
func fetchKeywords(cli *CommandLineClient, filters []filter.ArticleFilter) []filter.ArticleFilter {
	if cli.keywords != "" {
		keywords := strings.Split(cli.keywords, ",")
		uniqueKeywords := CheckUnique(keywords)
		filters = append(filters, filter.ByKeyword{Keywords: uniqueKeywords})
	}
	return filters
}

// fetchDateFilters extracts date filters from command line arguments and adds them to the filters.
func fetchDateFilters(cli *CommandLineClient, filters []filter.ArticleFilter) []filter.ArticleFilter {
	isValid, startDate, endDate := CheckData(cli.startDateStr, cli.endDateStr)
	if isValid {
		filters = append(filters, filter.ByDate{StartDate: startDate, EndDate: endDate})
	}
	return filters
}

// Print prints the provided articles.
func (cli *CommandLineClient) Print(articles []article.Article) {
	for _, article := range articles {
		fmt.Println("---------------------------------------------------")
		fmt.Println("Title:", article.Title)
		fmt.Println("Description:", article.Description)
		fmt.Println("Link:", article.Link)
		fmt.Println("Date:", article.Date)
	}
}
