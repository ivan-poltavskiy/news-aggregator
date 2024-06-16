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
		startDate string
		endDate   string
	}
	tests := []struct {
		name        string
		args        args
		wantError   bool
		wantIsValid bool
	}{
		{
			name: "Check data with only start date passed",
			args: args{
				startDate: "2003-05-05",
			},
			wantError:   true,
			wantIsValid: false,
		},
		{
			name: "Check data with only end date passed",
			args: args{
				endDate: "2003-05-05",
			},
			wantError:   true,
			wantIsValid: false,
		},

		{
			name: "Check data with two correct date passed",
			args: args{
				startDate: "2003-05-01",
				endDate:   "2003-05-05",
			},
			wantError:   false,
			wantIsValid: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotError, gotBool := ValidateDate(tt.args.startDate, tt.args.endDate)

			if gotError != nil != tt.wantError && gotBool != tt.wantIsValid {
				t.Errorf("Actual error: %v, expected %v", gotError != nil, tt.wantIsValid)
			}
		})
	}
}
