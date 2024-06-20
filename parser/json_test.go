package parser

import (
	"news-aggregator/entity/article"
	"news-aggregator/entity/source"
	"reflect"
	"testing"
	"time"
)

func TestJson_ParseSource(t *testing.T) {
	type args struct {
		path source.PathToFile
		name source.Name
	}
	tests := []struct {
		name string
		args args
		want []article.Article
	}{
		{
			name: "Parse valid JSON file",
			args: args{
				path: "../resources/testdata/json_articles.json",
				name: "testjson",
			},
			want: []article.Article{
				{Title: "Test Article 1", Description: "Description 1", Link: "http://example.com/1", Date: parseDate("2024-06-01"), SourceName: "testjson"},
				{Title: "Test Article 2", Description: "Description 2", Link: "http://example.com/2", Date: parseDate("2024-06-02"), SourceName: "testjson"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonFile := Json{}
			if got, _ := jsonFile.Parse(tt.args.path, tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
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
