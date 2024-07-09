package collector

import (
	"news-aggregator/entity/source"
	"reflect"
	"testing"
)

var testArticleCollector *articleCollector

func beforeEach() {
	sources := []source.Source{
		{Name: "bbc", PathToFile: "../resources/bbc-world-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "nbc", PathToFile: "../resources/nbc-news.json", SourceType: "JSON"},
	}
	testArticleCollector = &articleCollector{Sources: sources, Parsers: InitParsers()}
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
				currentSource: source.Source{Name: "bbc", PathToFile: "../resources/bbc-world-category-19-05-24.xml", SourceType: "RSS"},
				name:          "bbc",
			},
			wantQuantity: 54,
		},

		{name: "Test for nbc source",
			args: args{
				currentSource: source.Source{Name: "nbc", PathToFile: "../resources/nbc-news.json", SourceType: "JSON"},
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

func TestInitializeSource(t *testing.T) {
	tests := []struct {
		name    string
		sources []source.Source
	}{
		{
			name: "InitializeParsers with two sources",
			sources: []source.Source{
				{Name: "bbc", PathToFile: "../resources/bbc-world-category-19-05-24.xml", SourceType: "RSS"},
				{Name: "nbc", PathToFile: "../resources/nbc-news.json", SourceType: "JSON"},
			},
		},
		{
			name:    "InitializeParsers with no sources",
			sources: []source.Source{},
		},
		{
			name: "InitializeParsers with one source",
			sources: []source.Source{
				{Name: "nbc", PathToFile: "../resources/nbc-news.json", SourceType: "JSON"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testArticleCollector = &articleCollector{Sources: tt.sources, Parsers: InitParsers()}
			if !reflect.DeepEqual(testArticleCollector.Sources, tt.sources) {
				t.Errorf("Actual result = %v, expected = %v", testArticleCollector.Sources, tt.sources)
			}
		})
	}
}
