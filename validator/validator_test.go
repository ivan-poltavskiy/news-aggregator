package validator

import (
	"news_aggregator/entity/article"
	"testing"
	"time"
)

func TestCheckSource(t *testing.T) {
	type args struct {
		sources []article.Article
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Check empty sources",
			args: args{
				sources: []article.Article{},
			},
			want: false,
		},

		{
			name: "Check not empty sources",
			args: args{
				sources: []article.Article{
					{Title: "testTitle", Description: "testDescription", Link: "testLink", Date: time.Date(2003, time.January, 20, 0, 0, 0, 0, time.UTC)}},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateSource(tt.args.sources); got != tt.want {
				t.Errorf("Actual result %v, expexted %v", got, tt.want)
			}
		})
	}
}

func TestValidateDate(t *testing.T) {
	type args struct {
		startDate time.Time
		endDate   time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Check data with only start date passed",
			args: args{
				startDate: parseDate("2003-05-05"),
			},
			want: true,
		},
		{
			name: "Check data with only end date passed",
			args: args{
				endDate: parseDate("2003-05-05"),
			},
			want: true,
		},

		{
			name: "Check data with two correct date passed",
			args: args{
				startDate: parseDate("2003-05-01"),
				endDate:   parseDate("2003-05-05"),
			},
			want: false,
		},

		{
			name: "Check data with two incorrect date passed",
			args: args{
				startDate: parseDate("2003-05-05"),
				endDate:   parseDate("2003-05-01"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateDate(tt.args.startDate, tt.args.endDate)

			if (got == nil) == tt.want {
				t.Errorf("Actual: %v, expected %v", got, tt.want)
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
