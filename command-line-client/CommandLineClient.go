package command_line_client

import (
	"NewsAggregator/entity/source"
	. "NewsAggregator/filterService"
	. "NewsAggregator/initialization-data"
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

	InitializeSource()

	sourceNames := strings.Split(cli.sources, ",")
	uniqueNames := make(map[string]struct{})
	var sourceNameObjects []source.Name

	for _, name := range sourceNames {
		if _, ok := uniqueNames[name]; !ok {
			uniqueNames[name] = struct{}{}
			sourceNameObjects = append(sourceNameObjects, source.Name(name))
		}
	}

	articles, errorMessage := FindNewsForAllResources(sourceNameObjects)

	if errorMessage != "" {
		fmt.Println(errorMessage)
	}

	if cli.keywords != "" {
		articles = FilterNewsByKeyword(cli.keywords, articles)
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

		articles = FilterByDate(startDate, endDate, articles)
	}

	for _, article := range articles {
		fmt.Println("---------------------------------------------------")
		fmt.Println("Title of article:", article.Title)
		fmt.Println("Description:", article.Description)
		fmt.Println("Link:", article.Link)
		fmt.Println("Date:", article.Date)
	}
}
