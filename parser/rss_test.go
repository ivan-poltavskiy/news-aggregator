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
			},
			want: []article.Article{
				{
					Title:       "Test Article 1",
					Description: "Description 1",
					Link:        "http://example.com/1",
					Date:        time.Date(2024, time.June, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					Title:       "Test Article 2",
					Description: "Description 2",
					Link:        "http://example.com/2",
					Date:        time.Date(2024, time.June, 2, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rss := Rss{}
			if got := rss.ParseSource(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSource() = %v, want %v", got, tt.want)
			}
		})
	}
}
