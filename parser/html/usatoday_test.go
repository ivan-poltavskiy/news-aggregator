package html

import (
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
	"reflect"
	"testing"
	"time"
)

func TestUsaToday_ParseSource(t *testing.T) {
	type args struct {
		path source.PathToFile
		name source.Name
	}
	tests := []struct {
		name string
		args args
		want []news.News
	}{
		{
			name: "Parse valid HTML file",
			args: args{
				path: "../../mnt/resources/testdata/test_usatoday.html",
				name: "testusatoday",
			},
			want: []news.News{
				{
					Title:       "Test News 1",
					Description: "Description 1",
					Link:        "https://www.usatoday.com/story/1",
					Date:        time.Date(time.Now().Year(), time.June, 1, 0, 0, 0, 0, time.UTC),
					SourceName:  "testusatoday",
				},
				{
					Title:       "Test News 2",
					Description: "Description 2",
					Link:        "https://www.usatoday.com/story/2",
					Date:        time.Date(time.Now().Year(), time.June, 2, 0, 0, 0, 0, time.UTC),
					SourceName:  "testusatoday",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			htmlParser := UsaToday{}
			if got, _ := htmlParser.Parse(tt.args.path, tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
