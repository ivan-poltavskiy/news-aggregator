package filter

import (
	"news-aggregator/constant"
	"news-aggregator/entity/article"
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
		articles []article.Article
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []article.Article
	}{
		{
			name: "Articles within date range",
			fields: fields{
				StartDate: parseDate("2023-01-01"),
				EndDate:   parseDate("2023-12-31"),
			},
			args: args{
				articles: []article.Article{
					{Title: "Article 1", Date: parseDate("2023-03-15")},
					{Title: "Article 2", Date: parseDate("2023-06-10")},
					{Title: "Article 3", Date: parseDate("2024-01-01")},
				},
			},
			want: []article.Article{
				{Title: "Article 1", Date: parseDate("2023-03-15")},
				{Title: "Article 2", Date: parseDate("2023-06-10")},
			},
		},
		{
			name: "Empty article list",
			fields: fields{
				StartDate: parseDate("2023-01-01"),
				EndDate:   parseDate("2023-12-31"),
			},
			args: args{
				articles: []article.Article{},
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
