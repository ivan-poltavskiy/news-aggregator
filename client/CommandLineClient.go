package client

import (
	"NewsAggregator/aggregator"
	"NewsAggregator/entity/article"
	"NewsAggregator/filter_service"
	"flag"
	"fmt"
	"strings"
	"time"
)

type CommandLineClient struct {
	aggregator   aggregator.Aggregator
	sources      string
	keywords     string
	startDateStr string
	endDateStr   string
	help         bool
}

func NewCommandLineClient(aggregator aggregator.Aggregator) *CommandLineClient {
	cli := &CommandLineClient{aggregator: aggregator}
	flag.StringVar(&cli.sources, "sources", "", "Specify news sources separated by comma")
	flag.StringVar(&cli.keywords, "keywords", "", "Specify keywords to filter_service news articles")
	flag.StringVar(&cli.startDateStr, "startDate", "", "Specify start date (YYYY-MM-DD)")
	flag.StringVar(&cli.endDateStr, "endDate", "", "Specify end date (YYYY-MM-DD)")
	flag.BoolVar(&cli.help, "help", false, "Show help information")
	flag.Parse()
	return cli
}

// PrintUsage prints the usage instructions
func (cli *CommandLineClient) PrintUsage() {
	fmt.Println("Usage of NewsAggregator:" +
		"\nType --sources, and then list the resources you want to retrieve information from. \n" +
		"\nType --keywords, and then list the keywords by which you want to filter_service articles. \n" +
		"\nType --startDate and --endDate to filter_service by date. News published between the specified dates will be shown.")
}

func (cli *CommandLineClient) FetchArticles() []article.Article {
	if cli.help {
		cli.PrintUsage()
		return nil
	}

	if cli.sources == "" {
		fmt.Println("Please specify at least one news source using --sources flag.")
		return nil
	}

	sourceNames := strings.Split(cli.sources, ",")
	var filters []filter_service.FilterService

	if cli.keywords != "" {
		keywords := strings.Split(cli.keywords, ",")
		filters = append(filters, filter_service.ByKeyword{Keywords: keywords})
	}

	var startDate, endDate time.Time
	var err error
	if cli.startDateStr != "" || cli.endDateStr != "" {
		if cli.startDateStr == "" || cli.endDateStr == "" {
			fmt.Println("Please specify both start date and end date or omit them")
			return nil
		}
		startDate, err = time.Parse("2006-01-02", cli.startDateStr)
		if err != nil {
			fmt.Println("Error parsing start date:", err)
			return nil
		}
		endDate, err = time.Parse("2006-01-02", cli.endDateStr)
		if err != nil {
			fmt.Println("Error parsing end date:", err)
			return nil
		}
		filters = append(filters, filter_service.ByDate{StartDate: startDate, EndDate: endDate})
	}

	articles, errorMessage := cli.aggregator.Aggregate(sourceNames, filters...)
	if errorMessage != "" {
		fmt.Println(errorMessage)
	}

	return articles
}

func (cli *CommandLineClient) Print(articles []article.Article) {
	for _, article := range articles {
		fmt.Println("---------------------------------------------------")
		fmt.Println("Title of article:", article.Title)
		fmt.Println("Description:", article.Description)
		fmt.Println("Link:", article.Link)
		fmt.Println("Date:", article.Date)
	}
}
