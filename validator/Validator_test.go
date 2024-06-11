package validator

import (
	"NewsAggregator/entity/article"
	"reflect"
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
		want string
	}{
		{
			name: "Check empty sources",
			args: args{
				sources: []article.Article{},
			},
			want: "Please, specify at least one news source. " +
				"The program supports such news resources:\nABC, BBC, NBC, USA " +
				"Today and Washington Times.",
		},

		{
			name: "Check not empty sources",
			args: args{
				sources: []article.Article{
					{Title: "testTitle", Description: "testDescription", Link: "testLink", Date: time.Date(2003, time.January, 20, 0, 0, 0, 0, time.UTC)}},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckSource(tt.args.sources); got != tt.want {
				t.Errorf("Actual result %v, expexted %v", got, tt.want)
			}
		})
	}
}

func TestCheckData(t *testing.T) {
	type args struct {
		startDateStr string
		endDateStr   string
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 time.Time
		want2 time.Time
	}{
		{
			name: "Check data with only start date passed",
			args: args{
				startDateStr: "2003-05-05",
			},
			want:  false,
			want1: time.Time{},
			want2: time.Time{},
		},
		{
			name: "Check data with only end date passed",
			args: args{
				endDateStr: "2003-05-05",
			},
			want:  false,
			want1: time.Time{},
			want2: time.Time{},
		},

		{
			name: "Check data with two correct date passed",
			args: args{
				startDateStr: "2003-05-01",
				endDateStr:   "2003-05-05",
			},
			want:  true,
			want1: time.Date(2003, time.May, 1, 0, 0, 0, 0, time.UTC),
			want2: time.Date(2003, time.May, 5, 0, 0, 0, 0, time.UTC),
		},

		{
			name: "Check data with two incorrect date passed",
			args: args{
				startDateStr: "2003-05-05",
				endDateStr:   "2003-05-01",
			},
			want:  false,
			want1: time.Time{},
			want2: time.Time{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := CheckData(tt.args.startDateStr, tt.args.endDateStr)
			if got != tt.want {
				t.Errorf("Actual bool var %v, expected %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Actual start date = %v, expected %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("Actual end date = %v, expected %v", got2, tt.want2)
			}
		})
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
			if got := CheckUnique(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Actual result %v, expexted %v", got, tt.want)
			}
		})
	}
}
