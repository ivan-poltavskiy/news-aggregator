package client

import (
	"errors"
	"flag"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/reiver/go-porterstemmer"
	"news-aggregator/constant"
	"news-aggregator/entity/news"
	"news-aggregator/filter"
	"news-aggregator/sorter"
	"news-aggregator/validator"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"
)

// commandLineClient represents a command line client for the news-aggregator application.
type commandLineClient struct {
	aggregator       Aggregator
	sources          string
	keywords         string
	startDateStr     string
	endDateStr       string
	sortBy           string
	sortingBySources bool
	help             bool
}

// NewCommandLine creates and initializes a new commandLineClient with the provided aggregator.
func NewCommandLine(aggregator Aggregator) Client {
	cli := &commandLineClient{aggregator: aggregator}
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

	news, err := cli.aggregator.Aggregate(uniqueSources, filters...)
	if err != nil {
		return nil, err
	}
	news, fetchParametersError = Sorter.SortNews(sorter.DateSorter{}, news, cli.sortBy)
	if fetchParametersError != nil {
		return nil, fetchParametersError
	}
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
		panic(err)
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
	if err != nil {
		panic(err)
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

// fetchParameters extracts and validates command line parameters,
// including sources and filters, and returns them for use in news fetching.
func (cli *commandLineClient) fetchParameters() ([]filter.NewsFilter, error) {

	var filters []filter.NewsFilter

	filters = buildKeywordFilter(cli, filters)
	filters, err := buildDateFilters(cli, filters)
	if err != nil {
		return nil, err
	}
	return filters, nil
}

// buildKeywordFilter extracts keywords from command line arguments and adds them to the filters.
func buildKeywordFilter(cli *commandLineClient, filters []filter.NewsFilter) []filter.NewsFilter {
	if cli.keywords != "" {
		keywords := strings.Split(cli.keywords, ",")
		uniqueKeywords := checkUnique(keywords)
		filters = append(filters, filter.ByKeyword{Keywords: uniqueKeywords})
	}
	return filters
}

// buildDateFilters extracts date filters from command line arguments and adds them to the filters.
func buildDateFilters(cli *commandLineClient, filters []filter.NewsFilter) ([]filter.NewsFilter, error) {

	validationErr, isValid := validator.ValidateDate(cli.startDateStr, cli.endDateStr)

	if validationErr != nil {
		return nil, validationErr
	}
	if isValid {

		startDate, err := time.Parse(constant.DateOutputLayout, cli.startDateStr)

		if err != nil {
			return nil, errors.New("Invalid start date: " + cli.startDateStr)
		}

		endDate, err := time.Parse(constant.DateOutputLayout, cli.endDateStr)

		if err != nil {
			return nil, errors.New("Invalid end date: " + cli.endDateStr)
		}

		return append(filters, filter.ByDate{StartDate: startDate, EndDate: endDate}), nil
	}
	return filters, nil
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
