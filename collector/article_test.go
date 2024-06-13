package collector

import (
	"news_aggregator/entity/article"
	"news_aggregator/entity/source"
	"news_aggregator/parser"
	"reflect"
	"testing"
)

func beforeEach() {
	InitializeSource([]source.Source{
		{Name: "bbc", PathToFile: "../resources/bbc-world-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "nbc", PathToFile: "../resources/nbc-news.json", SourceType: "JSON"},
	})
	parser.Initialize()
}

func TestFindByResourcesName(t *testing.T) {
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
			got, _ := FindByResourcesName(tt.args.sourcesNames)
			if len(got) != tt.wantQuantity {
				t.Errorf("Actual result = %v, expected = %v", len(got), tt.wantQuantity)
			}
		})
	}
}

func Test_findForCurrentSource(t *testing.T) {
	beforeEach()
	type args struct {
		currentSource source.Source
		name          source.Name
		allArticles   []article.Article
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
				allArticles:   []article.Article{},
			},
			wantQuantity: 54,
		},

		{name: "Test for nbc source",
			args: args{
				currentSource: source.Source{Name: "nbc", PathToFile: "../resources/nbc-news.json", SourceType: "JSON"},
				name:          "nbc",
				allArticles:   []article.Article{},
			},
			wantQuantity: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := findForCurrentSource(tt.args.currentSource, tt.args.name, tt.args.allArticles); !reflect.DeepEqual(len(got), tt.wantQuantity) {
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
			name: "Initialize with two sources",
			sources: []source.Source{
				{Name: "bbc", PathToFile: "../resources/bbc-world-category-19-05-24.xml", SourceType: "RSS"},
				{Name: "nbc", PathToFile: "../resources/nbc-news.json", SourceType: "JSON"},
			},
		},
		{
			name:    "Initialize with no sources",
			sources: []source.Source{},
		},
		{
			name: "Initialize with one source",
			sources: []source.Source{
				{Name: "nbc", PathToFile: "../resources/nbc-news.json", SourceType: "JSON"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitializeSource(tt.sources)
			if !reflect.DeepEqual(Sources, tt.sources) {
				t.Errorf("Actual result = %v, expected = %v", Sources, tt.sources)
			}
		})
	}
}
