package command_line_client

import (
	"NewsAggregator/entity/source"
	. "NewsAggregator/filter"
	. "NewsAggregator/news"
	"flag"
	"fmt"
	"strings"
	"time"
)

type CommandLineClient struct {
	sources      string
	keywords     string
	startDateStr string
	endDateStr   string
	help         bool
}

func NewCommandLineClient() *CommandLineClient {
	cli := &CommandLineClient{}
	flag.StringVar(&cli.sources, "sources", "", "Specify news sources separated by comma")
	flag.StringVar(&cli.keywords, "keywords", "", "Specify keywords to filter news articles")
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
		"\nType --keywords, and then list the keywords by which you want to filter articles. \n" +
		"\nType --startDate and --endDate to filter by date. News published between the specified dates will be shown.")
}

// Run processes the CLI commands
func (cli *CommandLineClient) Run() {
	if cli.help {
		cli.PrintUsage()
		return
	}

	if cli.sources == "" {
		fmt.Println("Please specify at least one news source using --sources flag.")
		return
	}

	sourceNames := strings.Split(cli.sources, ",")
	uniqueSourceNames := filterUnique(sourceNames)
	var sourceNameObjects []source.Name

	for _, name := range uniqueSourceNames {
		sourceNameObjects = append(sourceNameObjects, source.Name(name))
	}

	articles, errorMessage := FindNewsByResources(sourceNameObjects)

	if errorMessage != "" {
		fmt.Println(errorMessage)
	}

	if cli.keywords != "" {
		keywordList := strings.Split(cli.keywords, ",")
		uniqueKeywords := filterUnique(keywordList)
		articles = NewsByKeywords(uniqueKeywords, articles)
	}

	var startDate, endDate time.Time
	var err error
	if cli.startDateStr != "" || cli.endDateStr != "" {
		if cli.startDateStr == "" || cli.endDateStr == "" {
			fmt.Println("Please specify both start date and end date or omit them")
			return
		}
		startDate, err = time.Parse("2006-01-02", cli.startDateStr)
		if err != nil {
			fmt.Println("Error parsing start date:", err)
			return
		}
		endDate, err = time.Parse("2006-01-02", cli.endDateStr)
		if err != nil {
			fmt.Println("Error parsing end date:", err)
			return
		}

		articles = ByDate(startDate, endDate, articles)
	}

	for _, article := range articles {
		fmt.Println("---------------------------------------------------")
		fmt.Println("Title of article:", article.Title)
		fmt.Println("Description:", article.Description)
		fmt.Println("Link:", article.Link)
		fmt.Println("Date:", article.Date)
	}
}

// filterUnique returns a slice containing only unique strings from the input slice.
func filterUnique(input []string) []string {
	uniqueMap := make(map[string]struct{})
	var uniqueList []string
	for _, item := range input {
		if _, ok := uniqueMap[item]; !ok {
			uniqueMap[item] = struct{}{}
			uniqueList = append(uniqueList, item)
		}
	}
	return uniqueList
}
