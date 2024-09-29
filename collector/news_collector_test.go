package collector

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"news-aggregator/constant"
	"news-aggregator/entity/source"
	"news-aggregator/storage"
	newsStorage "news-aggregator/storage/news"
	sourceStorage "news-aggregator/storage/source"
	"os"
	"testing"
)

var testArticleCollector *newsCollector

func beforeEach() {
	dir := os.TempDir()
	file, err := ioutil.TempFile(dir, "test-sources-*.json")
	if err != nil {
		log.Fatalf("Failed to create temp file: %v", err)
	}

	sources := []source.Source{
		{Name: "bbc", PathToFile: "../mnt/resources/testdata/bbc-world-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "nbc", PathToFile: "../mnt/resources/testdata/nbc-news.json", SourceType: "JSON"},
	}

	data, err := json.Marshal(sources)
	if err != nil {
		log.Fatalf("Failed to marshal sources: %v", err)
	}

	err = os.WriteFile(file.Name(), data, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to write to temp file: %v", err)
	}
	sourceStorage, _ := sourceStorage.NewJsonStorage(source.PathToFile(file.Name()))
	newsJsonStorage, _ := newsStorage.NewJsonStorage(source.PathToFile(constant.PathToResources))
	newStorage := storage.NewStorage(newsJsonStorage, sourceStorage)
	testArticleCollector = &newsCollector{sourceStorage: newStorage, parsers: GetDefaultParsers()}
}

func TestFindNewsByResourcesName(t *testing.T) {
	beforeEach()
	type args struct {
		sourcesNames []source.Name
	}
	tests := []struct {
		name         string
		args         args
		wantQuantity int
	}{
		{"Find articles from two sources by their names.",
			args{
				[]source.Name{
					"bbc",
					"nbc",
				},
			},
			154,
		},

		{"Find articles from one source by his name.",
			args{
				[]source.Name{"bbc"}},
			54,
		},
		{"Return zero if the source is not correct.",
			args{
				[]source.Name{
					"bbbc",
				},
			},
			0,
		},
		{"Return zero if the source was not passed on.",
			args{
				[]source.Name{},
			},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := testArticleCollector.FindNewsByResourcesName(tt.args.sourcesNames)
			if len(got) != tt.wantQuantity {
				t.Errorf("Actual result = %v, expected = %v", len(got), tt.wantQuantity)
			}
		})
	}
}

func TestFindNewsForCurrentSource(t *testing.T) {
	beforeEach()
	type args struct {
		currentSource source.Source
		name          source.Name
	}
	tests := []struct {
		name         string
		args         args
		wantQuantity int
	}{
		{name: "Test for bbc source",
			args: args{
				currentSource: source.Source{Name: "bbc", PathToFile: "../mnt/resources/testdata/bbc-world-category-19-05-24.xml", SourceType: "RSS"},
				name:          "bbc",
			},
			wantQuantity: 54,
		},

		{name: "Test for nbc source",
			args: args{
				currentSource: source.Source{Name: "nbc", PathToFile: "../mnt/resources/testdata/nbc-news.json", SourceType: "JSON"},
				name:          "nbc",
			},
			wantQuantity: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := testArticleCollector.findNewsForCurrentSource(tt.args.currentSource, tt.args.name)
			if len(got) != tt.wantQuantity {
				t.Errorf("Actual result = %v, expected = %v", len(got), tt.wantQuantity)
			}
		})
	}
}
