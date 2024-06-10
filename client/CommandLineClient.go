package client

import (
	"NewsAggregator/entity/article"
	"NewsAggregator/filter"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
)

// CommandLineClient represents a command line client for the NewsAggregator application.
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
	fmt.Println("Usage of NewsAggregator:" +
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
		filtersForTemplate = append(filtersForTemplate, "Keywords: "+strings.Join(uniqueKeywords, ", "))
		filters = append(filters, filter.ByKeyword{Keywords: uniqueKeywords})
	}
	return filters
}

// fetchDateFilters extracts date filters from command line arguments and adds them to the filters.
func fetchDateFilters(cli *CommandLineClient, filters []filter.ArticleFilter) []filter.ArticleFilter {
	isValid, startDate, endDate := CheckData(cli.startDateStr, cli.endDateStr)
	if isValid {
		filtersForTemplate = append(filtersForTemplate, fmt.Sprintf("Date filter - %s to %s", cli.startDateStr, cli.endDateStr))
		filters = append(filters, filter.ByDate{StartDate: startDate, EndDate: endDate})
	}
	return filters
}

var filtersForTemplate = []string{}

func GetFilters() string {
	return strings.Join(filtersForTemplate, ", ")
}

// Print prints the provided articles using a template.
func (cli *CommandLineClient) Print(articles []article.Article) {
	funcMap := template.FuncMap{
		"emphasise": func(keywords, text string) string {
			for _, keyword := range strings.Split(keywords, ",") {
				re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(keyword))
				text = re.ReplaceAllString(text, "**"+keyword+"**")
			}
			return text
		},
	}

	tmpl, err := template.New("articles").Funcs(funcMap).ParseFiles("client/OutputTemplate.tmpl")
	if err != nil {
		panic(err)
	}

	type articleData struct {
		Article  article.Article
		Keywords string
	}

	var data []articleData
	for _, art := range articles {
		data = append(data, articleData{
			Article:  art,
			Keywords: cli.keywords,
		})
	}

	outputData := struct {
		Filters  string
		Count    int
		Articles []articleData
	}{
		Filters:  GetFilters(),
		Count:    len(articles),
		Articles: data,
	}

	err = tmpl.ExecuteTemplate(os.Stdout, "articles", outputData)
	if err != nil {
		panic(err)
	}
}
