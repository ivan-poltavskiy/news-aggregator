package parser

import (
	"news_aggregator/entity/article"
	"news_aggregator/entity/source"
	"reflect"
	"testing"
	"time"
)

func TestRss_ParseSource(t *testing.T) {
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
			name: "Parse valid RSS file",
			args: args{
				path: "../resources/testdata/test_rss.xml",
				name: "testrss",
			},
			want: []article.Article{
				{
					Title:       "Test Article 1",
					Description: "Description 1",
					Link:        "http://example.com/1",
					Date:        time.Date(2024, time.June, 1, 0, 0, 0, 0, time.UTC),
					SourceName:  "testrss"},
				{
					Title:       "Test Article 2",
					Description: "Description 2",
					Link:        "http://example.com/2",
					Date:        time.Date(2024, time.June, 2, 0, 0, 0, 0, time.UTC),
					SourceName:  "testrss",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rss := Rss{}
			if got, _ := rss.ParseSource(tt.args.path, tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSource() = %v, want %v", got, tt.want)
			}
		})
	}
}
