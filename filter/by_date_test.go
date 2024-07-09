package filter

import (
	"news-aggregator/constant"
	"news-aggregator/entity/news"
	"reflect"
	"testing"
	"time"
)

func TestByDate_Filter(t *testing.T) {
	type fields struct {
		StartDate time.Time
		EndDate   time.Time
	}
	type args struct {
		articles []news.News
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []news.News
	}{
		{
			name: "News within date range",
			fields: fields{
				StartDate: parseDate("2023-01-01"),
				EndDate:   parseDate("2023-12-31"),
			},
			args: args{
				articles: []news.News{
					{Title: "News 1", Date: parseDate("2023-03-15")},
					{Title: "News 2", Date: parseDate("2023-06-10")},
					{Title: "News 3", Date: parseDate("2024-01-01")},
				},
			},
			want: []news.News{
				{Title: "News 1", Date: parseDate("2023-03-15")},
				{Title: "News 2", Date: parseDate("2023-06-10")},
			},
		},
		{
			name: "Empty news list",
			fields: fields{
				StartDate: parseDate("2023-01-01"),
				EndDate:   parseDate("2023-12-31"),
			},
			args: args{
				articles: []news.News{},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := ByDate{
				StartDate: tt.fields.StartDate,
				EndDate:   tt.fields.EndDate,
			}
			if got := f.Filter(tt.args.articles); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Actual result: = %v Expexted: %v", got, tt.want)
			}
		})
	}
}

func parseDate(dateStr string) time.Time {
	date, err := time.Parse(constant.DateOutputLayout, dateStr)
	if err != nil {
		panic(err)
	}
	return date
}
