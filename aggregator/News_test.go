package aggregator

import (
	"NewsAggregator/collector"
	"NewsAggregator/entity/source"
	"NewsAggregator/filter"
	"NewsAggregator/parser"
	"reflect"
	"testing"
	"time"
)

func beforeEach() {
	collector.InitializeSource([]source.Source{
		{Name: "bbc", PathToFile: "../resources/bbc-world-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "nbc", PathToFile: "../resources/nbc-news.json", SourceType: "JSON"},
	})
	parser.InitializeParserMap()
}

func TestNews_Aggregate(t *testing.T) {
	beforeEach()
	type args struct {
		sources []string
		filters []filter.ArticleFilter
	}
	tests := []struct {
		name         string
		args         args
		wantQuantity int
	}{
		{
			name: "Test with date filter from two sources",
			args: args{
				sources: []string{"bbc", "nbc"},
				filters: []filter.ArticleFilter{filter.ByDate{
					StartDate: parseDate("2024-05-17"),
					EndDate:   parseDate("2024-05-19"),
				}},
			},
			wantQuantity: 89,
		},
		{
			name: "Test with keyword filter from one source",
			args: args{
				sources: []string{"bbc"},
				filters: []filter.ArticleFilter{filter.ByKeyword{
					Keywords: []string{"Trump"},
				}},
			},
			wantQuantity: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			na := &News{}
			got, _ := na.Aggregate(tt.args.sources, tt.args.filters...)
			if !reflect.DeepEqual(len(got), tt.wantQuantity) {
				t.Errorf("Aggregate() got = %v, wantQuantity %v", len(got), tt.wantQuantity)
			}
		})
	}
}

func parseDate(dateStr string) time.Time {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		panic(err)
	}
	return date
}
