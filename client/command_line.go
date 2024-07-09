package client

import (
	"flag"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/reiver/go-porterstemmer"
	"news-aggregator/constant"
	"news-aggregator/entity/news"
	"github.com/sirupsen/logrus"
	"news-aggregator/entity/article"
	"news-aggregator/filter"
	"news-aggregator/sorter"
	"news-aggregator/validator"
	"os"
	"regexp"
	"strings"
	"text/template"
)

// commandLineClient represents a command line client for the news-aggregator application.
type commandLineClient struct {
	aggregator       Aggregator
	sources          []string
	keywords         string
	startDateStr     string
	endDateStr       string
	sortBy           string
	sortingBySources bool
	help             bool
	DateSorter       DateSorter
	filters          []filter.ArticleFilter
}

// NewCommandLine creates and initializes a new commandLineClient with the provided aggregator.
func NewCommandLine(aggregator Aggregator) Client {
	cli := &commandLineClient{aggregator: aggregator}
	cli.DateSorter = DateSorter{}
	var sourcesStr string
	flag.StringVar(&sourcesStr, "sources", "", "Specify news sources separated by comma")
	flag.StringVar(&cli.keywords, "keywords", "", "Specify keywords to filter collector articles")
	flag.StringVar(&cli.startDateStr, "startDate", "", "Specify start date (YYYY-MM-DD)")
	flag.StringVar(&cli.endDateStr, "endDate", "", "Specify end date (YYYY-MM-DD)")
	flag.StringVar(&cli.sortBy, "sortBy", "", "Specify sort by DESC/ASC.")
	flag.BoolVar(&cli.sortingBySources, "sortingBySources", false, "Enable sorting articles by sources")
	flag.BoolVar(&cli.help, "help", false, "Show help information")
	flag.Parse()

	cli.sources = checkUnique(strings.Split(sourcesStr, ","))
	cli.filters = buildKeywordFilter(cli.keywords, cli.filters)
	var err error
	cli.filters, err = buildDateFilters(cli.startDateStr, cli.endDateStr, cli.filters)
	if err != nil {
		logrus.Error("Command line client: Date filter error: ", err)
	}

	logrus.Info("Command line client: Initialized with sources: ", cli.sources, " and filters: ", cli.filters)
	return cli
}

// FetchNews fetches news based on the command line arguments.
func (cli *commandLineClient) FetchNews() ([]news.News, error) {
	if cli.help {
		cli.printUsage()
		return nil, nil
	}

	filters, fetchParametersError := cli.fetchParameters()
	uniqueSources := checkUnique(strings.Split(cli.sources, ","))
	if fetchParametersError != nil {
		return nil, fetchParametersError
	}
	logrus.Info("Command line client: Fetching articles with sources: ", cli.sources, " and filters: ", cli.filters)
	news, err := cli.aggregator.Aggregate(uniqueSources, filters...)
	if err != nil {
		logrus.Error("Command line client: Aggregation error: ", err)
		return nil, err
	}
	news, fetchParametersError = Sorter.SortNews(sorter.DateSorter{}, news, cli.sortBy)
	if fetchParametersError != nil {
		logrus.Error("Command line client: Date sorting error: ", fetchParametersError)
		return nil, fetchParametersError
	}
	logrus.Info("Command line client: Articles fetched and sorted by date with sortBy: ", cli.sortBy)
	return news, nil
}

// Print outputs the transferred news.
func (cli *commandLineClient) Print(newsForOutput []news.News) {
	funcMap := sprig.FuncMap()
	funcMap["emphasise"] = func(keywords, text string) string {
		if keywords == "" {
			return text
		} else {
			for _, keyword := range strings.Split(keywords, ",") {
				stemString := porterstemmer.StemString(strings.ToLower(keyword))
				re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(stemString))
				text = re.ReplaceAllString(text, "//"+stemString+"//")
			}
			return text
		}
	}

	tmpl, err := template.New("news").Funcs(funcMap).ParseFiles("client/OutputTemplate.tmpl")
	if err != nil {
		logrus.Fatal("Command line client: Template parsing error: ", err)
	}

	type newsData struct {
		News             news.News
		Keywords         string
		SortingBySources bool
	}

	var data []newsData
	for _, n := range newsForOutput {
		data = append(data, newsData{
			News:             n,
			Keywords:         cli.keywords,
			SortingBySources: cli.sortingBySources,
		})
	}
	outputData := struct {
		Filters          []string
		Count            int
		News             []newsData
		NewsBySource     map[string][]newsData
		SortingBySources bool
	}{
		Filters:          []string{cli.keywords, cli.startDateStr, cli.endDateStr},
		Count:            len(newsForOutput),
		News:             data,
		SortingBySources: cli.sortingBySources,
	}

	if cli.sortingBySources {
		outputData.NewsBySource = make(map[string][]newsData)
		for _, n := range newsForOutput {
			sourceName := string(n.SourceName)
			outputData.NewsBySource[sourceName] = append(outputData.NewsBySource[sourceName], newsData{
				News:             n,
				Keywords:         cli.keywords,
				SortingBySources: cli.sortingBySources,
			})
		}
	}

	err = tmpl.ExecuteTemplate(os.Stdout, "news", outputData)
	logrus.Info("Command line client: Printing articles with count: ", outputData.Count)
	err = tmpl.ExecuteTemplate(os.Stdout, "articles", outputData)
	if err != nil {
		logrus.Fatal("Command line client: Template execution error: ", err)
	}
}

// printUsage prints the usage instructions
func (cli *commandLineClient) printUsage() {
	fmt.Println("Usage of news-aggregator:" +
		"\nType --sources, and then list the resources you want to retrieve information from. " +
		"The program supports such news resources:\nABC, BBC, NBC, USA Today and Washington Times. \n" +
		"\nType --keywords, and then list the keywords by which you want to filter articles. \n" +
		"\nType --startDate and --endDate to filter by date. News published between the specified dates will be shown." +
		"Date format - yyyy-mm-dd" + "" +
		"Type --sortBy to sort by DESC/ASC." + "Type --sortingBySources to sort by sources.")
}
