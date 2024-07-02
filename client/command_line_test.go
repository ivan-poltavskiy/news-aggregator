package client

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"io"
	"news-aggregator/aggregator/mock_aggregator"
	"news-aggregator/entity/article"
	"news-aggregator/filter"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestCommandLineClient_FetchArticles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAggregator := mock_aggregator.NewMockAggregator(ctrl)

	type fields struct {
		aggregator   Aggregator
		sources      []string
		keywords     string
		startDateStr string
		endDateStr   string
		help         bool
	}
	tests := []struct {
		name   string
		fields fields
		setup  func()
		want   []article.Article
	}{
		{
			name: "Test with articles",
			fields: fields{
				aggregator:   mockAggregator,
				sources:      []string{"source1", "source2"},
				keywords:     "test",
				startDateStr: "2023-01-01",
				endDateStr:   "2023-12-31",
				help:         false,
			},
			setup: func() {
				mockAggregator.EXPECT().
					Aggregate([]string{"source1", "source2"}, gomock.Any()).
					Return([]article.Article{
						{Title: "Test Title", Description: "Test Description", Link: "http://test.com", Date: time.Date(2023, time.May, 1, 0, 0, 0, 0, time.UTC)},
					}, nil)
			},
			want: []article.Article{
				{Title: "Test Title", Description: "Test Description", Link: "http://test.com", Date: time.Date(2023, time.May, 1, 0, 0, 0, 0, time.UTC)},
			},
		},
		{
			name: "Test with error message",
			fields: fields{
				aggregator:   mockAggregator,
				sources:      []string{""},
				keywords:     "test",
				startDateStr: "2023-01-01",
				endDateStr:   "2023-12-31",
				help:         false,
			},
			setup: func() {
				mockAggregator.EXPECT().
					Aggregate([]string{""}, gomock.Any())
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			cli := &commandLineClient{
				aggregator:   tt.fields.aggregator,
				sources:      tt.fields.sources,
				keywords:     tt.fields.keywords,
				startDateStr: tt.fields.startDateStr,
				endDateStr:   tt.fields.endDateStr,
				help:         tt.fields.help,
			}
			if got, _ := cli.FetchArticles(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Actual result = %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestFetchKeywords(t *testing.T) {
	cli := &commandLineClient{keywords: "keyword1,keyword2"}
	var filters []filter.ArticleFilter
	filters = buildKeywordFilter(cli.keywords, filters)

	expectedFilters := []filter.ArticleFilter{
		filter.ByKeyword{Keywords: []string{"keyword1", "keyword2"}},
	}

	if !reflect.DeepEqual(filters, expectedFilters) {
		t.Errorf("buildKeywordFilter() failed, got: %v, want: %v", filters, expectedFilters)
	}
}

func TestFetchDateFilters(t *testing.T) {
	cli := &commandLineClient{startDateStr: "2023-01-01", endDateStr: "2023-12-31"}
	var filters []filter.ArticleFilter
	filters, _ = buildDateFilters(cli.startDateStr, cli.endDateStr, filters)

	startDate := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC)
	expectedFilters := []filter.ArticleFilter{
		filter.ByDate{StartDate: startDate, EndDate: endDate},
	}

	if !reflect.DeepEqual(filters, expectedFilters) {
		t.Errorf("buildDateFilters() failed, got: %v, want: %v", filters, expectedFilters)
	}
}

func TestCommandLineClient_printUsage(t *testing.T) {
	cli := &commandLineClient{}
	expectedOutput := "Usage of news-aggregator:" +
		"\nType --sources, and then list the resources you want to retrieve information from. " +
		"The program supports such news resources:\nABC, BBC, NBC, USA Today and Washington Times. \n" +
		"\nType --keywords, and then list the keywords by which you want to filter articles. \n" +
		"\nType --startDate and --endDate to filter by date. News published between the specified dates will be shown." +
		"Date format - yyyy-mm-dd" + "" +
		"Type --sortBy to sort by DESC/ASC." + "Type --sortingBySources to sort by sources."

	var output bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cli.printUsage()

	w.Close()
	os.Stdout = old
	io.Copy(&output, r)

	if strings.TrimSpace(output.String()) != strings.TrimSpace(expectedOutput) {
		t.Errorf("Expected:\n%s\nGot:\n%s", expectedOutput, output.String())
	}
}

func TestCheckUnique(t *testing.T) {
	type args struct {
		input []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Check unique with 3 identical values",
			args: args{
				input: []string{"Check", "Check", "Check"},
			},
			want: []string{"Check"},
		},

		{
			name: "Check unique with 5 identical values",
			args: args{
				input: []string{"Check", "Random", "Check", "Random", "Check"},
			},
			want: []string{"Check", "Random"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkUnique(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Actual result %v, expected %v", got, tt.want)
			}
		})
	}
}
