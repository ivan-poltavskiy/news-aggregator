package aggregator

import (
	"github.com/golang/mock/gomock"
	aggregator "news-aggregator/aggregator/mock_aggregator"
	"news-aggregator/constant"
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
	"news-aggregator/filter"
	"reflect"
	"testing"
	"time"
)

func TestNews_Aggregate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	constant.PathToStorage = "../" + constant.PathToStorage

	mockCollector := aggregator.NewMockCollector(ctrl)
	type args struct {
		sources []string
		filters []filter.NewsFilter
	}
	tests := []struct {
		name         string
		args         args
		setup        func()
		wantQuantity int
		wantErr      bool
	}{
		{
			name: "Test with date filter from bbc and nbc sources",
			args: args{
				sources: []string{"bbc", "nbc"},
				filters: []filter.NewsFilter{filter.ByDate{
					StartDate: parseDate("2024-05-17"),
					EndDate:   parseDate("2024-05-19"),
				}},
			},
			setup: func() {
				mockCollector.EXPECT().FindNewsByResourcesName([]source.Name{"bbc", "nbc"}).
					Return([]news.News{
						{Title: "Test Title", Description: "Test Description", Link: "http://test.com", Date: time.Date(2024, time.May, 18, 0, 0, 0, 0, time.UTC)},
					}, nil)
			},
			wantQuantity: 1,
			wantErr:      false,
		},
		{
			name: "Test with keyword filter from bbc source with rss type",
			args: args{
				sources: []string{"bbc"},
				filters: []filter.NewsFilter{filter.ByKeyword{
					Keywords: []string{"Trump"},
				}},
			},
			setup: func() {
				mockCollector.EXPECT().FindNewsByResourcesName([]source.Name{"bbc"}).
					Return([]news.News{
						{Title: "Trump News 1", Description: "Description 1", Link: "http://test1.com", Date: time.Now()},
						{Title: "Trump News 2", Description: "Description 2", Link: "http://test2.com", Date: time.Now()},
					}, nil)
			},
			wantQuantity: 2,
			wantErr:      false,
		},
		{
			name: "Test with keyword filter from Usa Today source with html type",
			args: args{
				sources: []string{"usatoday"},
				filters: []filter.NewsFilter{filter.ByKeyword{
					Keywords: []string{"ukr"},
				}},
			},
			setup: func() {
				mockCollector.EXPECT().FindNewsByResourcesName([]source.Name{"usatoday"}).
					Return([]news.News{
						{Title: "Ukraine News 1", Description: "Description 1", Link: "http://test1.com", Date: time.Now()},
						{Title: "Ukraine News 2", Description: "Description 2", Link: "http://test2.com", Date: time.Now()},
						{Title: "Ukraine News 3", Description: "Description 3", Link: "http://test3.com", Date: time.Now()},
						{Title: "Ukraine News 4", Description: "Description 4", Link: "http://test4.com", Date: time.Now()},
					}, nil)
			},
			wantQuantity: 4,
			wantErr:      false,
		},
		{
			name: "Test with keyword filter from NBC source with JSON type",
			args: args{
				sources: []string{"nbc"},
				filters: []filter.NewsFilter{filter.ByKeyword{
					Keywords: []string{"ukr"},
				}},
			},
			setup: func() {
				mockCollector.EXPECT().FindNewsByResourcesName([]source.Name{"nbc"}).
					Return([]news.News{
						{Title: "Ukraine News from NBC", Description: "Description", Link: "http://test.com", Date: time.Now()},
					}, nil)
			},
			wantQuantity: 1,
			wantErr:      false,
		},
		{
			name: "Test with non-existent sources",
			args: args{
				sources: []string{"source1", "source2"},
				filters: nil,
			},
			setup:        func() {},
			wantQuantity: 0,
			wantErr:      true,
		},
		{
			name: "Test with empty sources",
			args: args{
				sources: nil,
				filters: nil,
			},
			setup:        func() {},
			wantQuantity: 0,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			na := New(mockCollector)
			got, err := na.Aggregate(tt.args.sources, tt.args.filters...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Aggregate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
