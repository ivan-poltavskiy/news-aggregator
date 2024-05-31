package client

import (
	"NewsAggregator/entity/article"
	"NewsAggregator/service/filter"
	"flag"
	"fmt"
	"strings"
	"time"
)

// CommandLine represents a command line client for the NewsAggregator application.
type CommandLine struct {
	aggregator   Aggregator
	sources      string
	keywords     string
	startDateStr string
	endDateStr   string
	help         bool
}

// NewCommandLine creates and initializes a new CommandLine with the provided aggregator.
func NewCommandLine(aggregator Aggregator) *CommandLine {
	cli := &CommandLine{aggregator: aggregator}
	flag.StringVar(&cli.sources, "sources", "", "Specify news sources separated by comma")
	flag.StringVar(&cli.keywords, "keywords", "", "Specify keywords to filter news articles")
	flag.StringVar(&cli.startDateStr, "startDate", "", "Specify start date (YYYY-MM-DD)")
	flag.StringVar(&cli.endDateStr, "endDate", "", "Specify end date (YYYY-MM-DD)")
	flag.BoolVar(&cli.help, "help", false, "Show help information")
	flag.Parse()
	return cli
}

// printUsage prints the usage instructions
func (cli *CommandLine) printUsage() {
	fmt.Println("Usage of NewsAggregator:" +
		"\nType --sources, and then list the resources you want to retrieve information from. \n" +
		"\nType --keywords, and then list the keywords by which you want to filter articles. \n" +
		"\nType --startDate and --endDate to filter by date. News published between the specified dates will be shown.")
}

// FetchArticles fetches articles based on the command line arguments.
func (cli *CommandLine) FetchArticles() []article.Article {
	if cli.help {
		cli.printUsage()
		return nil
	}

	if cli.sources == "" {
		fmt.Println("Please specify at least one news source using --sources flag.")
		return nil
	}

	sourceNames := strings.Split(cli.sources, ",")
	var filters []filter.Service

	filters = fetchKeywords(cli, filters)
	filters = fetchDateFilters(cli, filters)

	articles, errorMessage := cli.aggregator.Aggregate(sourceNames, filters...)
	if errorMessage != "" {
		fmt.Println(errorMessage)
	}

	return articles
}

// fetchKeywords extracts keywords from command line arguments and adds them to the filters.
func fetchKeywords(cli *CommandLine, filters []filter.Service) []filter.Service {
	if cli.keywords != "" {
		keywords := strings.Split(cli.keywords, ",")
		filters = append(filters, filter.ByKeyword{Keywords: keywords})
	}
	return filters
}

// fetchDateFilters extracts date filters from command line arguments and adds them to the filters.
func fetchDateFilters(cli *CommandLine, filters []filter.Service) []filter.Service {
	if cli.startDateStr != "" || cli.endDateStr != "" {
		if cli.startDateStr == "" || cli.endDateStr == "" {
			fmt.Println("Please specify both start date and end date or omit them")
			return nil
		}
		startDate, err := time.Parse("2006-01-02", cli.startDateStr)
		if err != nil {
			fmt.Println("Error parsing start date:", err)
			return nil
		}
		endDate, err := time.Parse("2006-01-02", cli.endDateStr)
		if err != nil {
			fmt.Println("Error parsing end date:", err)
			return nil
		}
		filters = append(filters, filter.ByDate{StartDate: startDate, EndDate: endDate})
	}
	return filters
}

// Print prints the provided articles.
func (cli *CommandLine) Print(articles []article.Article) {
	for _, article := range articles {
		fmt.Println("---------------------------------------------------")
		fmt.Println("Title of article:", article.Title)
		fmt.Println("Description:", article.Description)
		fmt.Println("Link:", article.Link)
		fmt.Println("Date:", article.Date)
	}
}
