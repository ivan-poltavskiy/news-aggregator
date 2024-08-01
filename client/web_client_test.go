package client

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"news-aggregator/entity/news"
	"news-aggregator/filter"
	"news-aggregator/mocks"
	"news-aggregator/sorter"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestWebClient_FetchNews(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAggregator := mocks.NewMockAggregator(ctrl)

	type fields struct {
		aggregator       Aggregator
		Sources          []string
		sortBy           string
		sortingBySources bool
		help             bool
		DateSorter       sorter.DateSorter
		filters          []filter.NewsFilter
		output           http.ResponseWriter
	}
	tests := []struct {
		name    string
		fields  fields
		setup   func()
		want    []news.News
		wantErr bool
	}{
		{
			name: "Fetch with valid filters",
			fields: fields{
				aggregator:       mockAggregator,
				Sources:          []string{"source1", "source2"},
				sortBy:           "desc",
				sortingBySources: true,
				filters: []filter.NewsFilter{
					filter.ByKeyword{Keywords: []string{"test"}},
					filter.ByDate{StartDate: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC)},
				},
				output: httptest.NewRecorder(),
			},
			setup: func() {
				mockAggregator.EXPECT().
					Aggregate([]string{"source1", "source2"}, gomock.Any()).
					Return([]news.News{
						{Title: "Test Title", Description: "Test Description", Link: "http://test.com", Date: time.Date(2023, time.May, 1, 0, 0, 0, 0, time.UTC)},
					}, nil)
			},
			want: []news.News{
				{Title: "Test Title", Description: "Test Description", Link: "http://test.com", Date: time.Date(2023, time.May, 1, 0, 0, 0, 0, time.UTC)},
			},
			wantErr: false,
		},
		{
			name: "Fetch with error",
			fields: fields{
				aggregator:       mockAggregator,
				Sources:          []string{""},
				sortBy:           "desc",
				sortingBySources: true,
				filters: []filter.NewsFilter{
					filter.ByKeyword{Keywords: []string{"test"}},
					filter.ByDate{StartDate: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC)},
				},
				output: httptest.NewRecorder(),
			},
			setup: func() {
				mockAggregator.EXPECT().
					Aggregate([]string{""}, gomock.Any()).
					Return(nil, fmt.Errorf("aggregation error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			webClient := &WebClient{
				aggregator:       tt.fields.aggregator,
				Sources:          tt.fields.Sources,
				sortBy:           tt.fields.sortBy,
				sortingBySources: tt.fields.sortingBySources,
				help:             tt.fields.help,
				DateSorter:       tt.fields.DateSorter,
				filters:          tt.fields.filters,
				output:           tt.fields.output,
			}
			got, err := webClient.FetchNews()
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchNews() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchNews() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebClient_Print(t *testing.T) {
	type fields struct {
		output http.ResponseWriter
	}
	type args struct {
		news []news.News
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Valid Print",
			fields: fields{
				output: httptest.NewRecorder(),
			},
			args: args{
				news: []news.News{
					{Title: "Test Title", Description: "Test Description", Link: "http://test.com", Date: time.Date(2023, time.May, 1, 0, 0, 0, 0, time.UTC)},
				},
			},
			want: `[{"title":"Test Title","description":"Test Description","url":"http://test.com","publishedAt":"2023-05-01T00:00:00Z","SourceName":""}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			webClient := &WebClient{
				output: tt.fields.output,
			}
			webClient.Print(tt.args.news)
			result := webClient.output.(*httptest.ResponseRecorder).Body.String()
			if strings.TrimSpace(result) != tt.want {
				t.Errorf("Print() got = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestWebClient_printUsage(t *testing.T) {
	type fields struct {
		output http.ResponseWriter
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Valid printUsage",
			fields: fields{
				output: httptest.NewRecorder(),
			},
			want: "Usage of news-aggregator:\n" +
				"Type --sources, and then list the news you want to retrieve information from. " +
				"The program supports such news news:\nABC, BBC, NBC, USA Today and Washington Times. \n" +
				"\nType --keywords, and then list the keywords by which you want to filter articles. \n" +
				"\nType --startDate and --endDate to filter by date. News published between the specified dates will be shown." +
				"Date format - yyyy-mm-dd" +
				"\nType --sortBy to sort by DESC/ASC." +
				"\nType --sortingBySources to sort by sources.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			webClient := &WebClient{
				output: tt.fields.output,
			}
			webClient.printUsage()
			result := webClient.output.(*httptest.ResponseRecorder).Body.String()
			if strings.TrimSpace(result) != tt.want {
				t.Errorf("printUsage() got = %v, want %v", result, tt.want)
			}
		})
	}
}
