package aggregator

import (
	"news_aggregator/collector"
	"news_aggregator/entity/source"
	"news_aggregator/filter"
	"news_aggregator/parser"
	"reflect"
	"testing"
	"time"
)

func beforeEach() {
	collector.InitializeSource([]source.Source{
		{Name: "bbc", PathToFile: "../resources/bbc-world-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "nbc", PathToFile: "../resources/nbc-news.json", SourceType: "JSON"},
		{Name: "usatoday", PathToFile: "../resources/usatoday-world-news.html", SourceType: "Html"},
	})
	parser.Initialize()
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
			name: "Test with date filter from bbc and nbc sources",
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
			name: "Test with keyword filter from bbc source with rss type",
			args: args{
				sources: []string{"bbc"},
				filters: []filter.ArticleFilter{filter.ByKeyword{
					Keywords: []string{"Trump"},
				}},
			},
			wantQuantity: 2,
		},
		{
			name: "Test with keyword filter from Usa Today source with html type",
			args: args{
				sources: []string{"usatoday"},
				filters: []filter.ArticleFilter{filter.ByKeyword{
					Keywords: []string{"ukr"},
				}},
			},
			wantQuantity: 4,
		},
		{
			name: "Test with keyword filter from NBC source with JSON type",
			args: args{
				sources: []string{"nbc"},
				filters: []filter.ArticleFilter{filter.ByKeyword{
					Keywords: []string{"ukr"},
				}},
			},
			wantQuantity: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			na := &news{}
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
