package client

import (
	"flag"
	"fmt"
	"news_aggregator/entity/article"
	"news_aggregator/entity/source"
	"news_aggregator/filter"
	"news_aggregator/validator"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"time"
)

// CommandLineClient represents a command line client for the news-aggregator application.
type CommandLineClient struct {
	aggregator       Aggregator
	sources          string
	keywords         string
	startDateStr     string
	endDateStr       string
	sortBy           string
	sortingBySources bool
	help             bool
}

var filtersForTemplate []string

// NewCommandLine creates and initializes a new CommandLineClient with the provided aggregator.
func NewCommandLine(aggregator Aggregator) *CommandLineClient {
	cli := &CommandLineClient{aggregator: aggregator}
	flag.StringVar(&cli.sources, "sources", "", "Specify news sources separated by comma")
	flag.StringVar(&cli.keywords, "keywords", "", "Specify keywords to filter collector articles")
	flag.StringVar(&cli.startDateStr, "startDate", "", "Specify start date (YYYY-MM-DD)")
	flag.StringVar(&cli.endDateStr, "endDate", "", "Specify end date (YYYY-MM-DD)")
	flag.StringVar(&cli.sortBy, "sortBy", "", "Specify sort by DESC/ASC.")
	flag.BoolVar(&cli.sortingBySources, "sortingBySources", false, "Enable sorting articles by sources")
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
		"Date format - yyyy-mm-dd" + "" +
		"Type --sortBy to sort by DESC/ASC." + "Type --sortingBySources to sort by sources.")
}

// FetchArticles fetches articles based on the command line arguments.
func (cli *CommandLineClient) FetchArticles() ([]article.Article, error) {
	if cli.help {
		cli.printUsage()
		return nil, nil
	}

	filters, uniqueSources, fetchParametrsError := fetchParameters(cli)
	if fetchParametrsError != nil {
		return nil, fetchParametrsError
	}

	articles, err := cli.aggregator.Aggregate(uniqueSources, filters...)
	if err != nil {
		return nil, err
	}
	cli.sortedByDate(articles)
	return articles, nil
}

// Print prints the provided articles using a template.
func (cli *CommandLineClient) Print(articles []article.Article) {
	funcMap := template.FuncMap{
		"emphasise": func(keywords, text string) string {
			if keywords == "" {
				return text
			} else {
				for _, keyword := range strings.Split(keywords, ",") {
					re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(keyword))
					text = re.ReplaceAllString(text, "//"+keyword+"//")
				}
				return text
			}
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
		Filters          string
		Count            int
		Articles         []articleData
		ArticlesBySource map[source.Name][]articleData
		SortingBySources bool
	}{
		Filters:          strings.Join(filtersForTemplate, ", "),
		Count:            len(articles),
		Articles:         data,
		SortingBySources: cli.sortingBySources,
	}

	if cli.sortingBySources {
		outputData.ArticlesBySource = make(map[source.Name][]articleData)
		for _, art := range articles {
			sourceName := art.SourceName
			outputData.ArticlesBySource[sourceName] = append(outputData.ArticlesBySource[sourceName], articleData{
				Article:  art,
				Keywords: cli.keywords,
			})
		}
	}

	err = tmpl.ExecuteTemplate(os.Stdout, "articles", outputData)
	if err != nil {
		panic(err)
	}
}

// sortedByDate sorts news by ASC or DESC.
func (cli *CommandLineClient) sortedByDate(articles []article.Article) {

	if strings.ToLower(cli.sortBy) == "asc" {
		sort.Slice(articles, func(i, j int) bool {
			return articles[i].Date.Before(articles[j].Date)
		})
	} else if strings.ToLower(cli.sortBy) == "desc" {
		sort.Slice(articles, func(i, j int) bool {
			return articles[i].Date.After(articles[j].Date)
		})
	}
}

// fetchParameters extracts and validates command line parameters,
// including sources and filters, and returns them for use in article fetching.
func fetchParameters(cli *CommandLineClient) ([]filter.ArticleFilter, []string, error) {
	sourceNames := strings.Split(cli.sources, ",")
	var filters []filter.ArticleFilter

	filters = buildKeywordFilter(cli, filters)
	filters, err := buildDateFilters(cli, filters)
	if err != nil {
		return nil, nil, err
	}
	uniqueSources := checkUnique(sourceNames)
	return filters, uniqueSources, nil
}

// buildKeywordFilter extracts keywords from command line arguments and adds them to the filters.
func buildKeywordFilter(cli *CommandLineClient, filters []filter.ArticleFilter) []filter.ArticleFilter {
	if cli.keywords != "" {
		keywords := strings.Split(cli.keywords, ",")
		uniqueKeywords := checkUnique(keywords)
		filtersForTemplate = append(filtersForTemplate, "Keywords: "+strings.Join(uniqueKeywords, ", "))
		filters = append(filters, filter.ByKeyword{Keywords: uniqueKeywords})
	}
	return filters
}

// buildDateFilters extracts date filters from command line arguments and adds them to the filters.
func buildDateFilters(cli *CommandLineClient, filters []filter.ArticleFilter) ([]filter.ArticleFilter, error) {

	startDate, _ := time.Parse("2006-01-02", cli.startDateStr)

	endDate, _ := time.Parse("2006-01-02", cli.endDateStr)

	validationErr := validator.ValidateDate(startDate, endDate)

	if validationErr == nil {
		filtersForTemplate = append(filtersForTemplate, fmt.Sprintf("Date filter - %s to %s", cli.startDateStr, cli.endDateStr))
		filters = append(filters, filter.ByDate{StartDate: startDate, EndDate: endDate})
		return filters, nil
	}
	return nil, validationErr
}

// CheckUnique returns a slice containing only unique strings from the input slice.
func checkUnique(input []string) []string {
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
