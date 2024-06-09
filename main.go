package main

import (
	"NewsAggregator/aggregator"
	"NewsAggregator/client"
	"NewsAggregator/collector"
	"NewsAggregator/entity/source"
	"NewsAggregator/parser"
	"os"
	"strings"
	"text/template"
)

func main() {

	collector.InitializeSource([]source.Source{
		{Name: "bbc", PathToFile: "resources/bbc-world-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "nbc", PathToFile: "resources/nbc-news.json", SourceType: "JSON"},
		{Name: "abc", PathToFile: "resources/abcnews-international-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "washington", PathToFile: "resources/washingtontimes-world-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "usatoday", PathToFile: "resources/usatoday-world-news.html", SourceType: "Html"},
	})
	parser.InitializeParserMap()

	newsAggregator := aggregator.New()
	cli := client.NewCommandLine(newsAggregator)
	articles := cli.FetchArticles()

	funcMap := template.FuncMap{
		"highlight": func(text, keywords string) string {
			for _, keyword := range strings.Split(keywords, ",") {
				text = strings.ReplaceAll(text, keyword, ""+keyword+"")
			}
			return text
		}}
	tmpl, err := template.New("articles").Funcs(funcMap).ParseFiles("client/OutputTemplate.tmpl")
	if err != nil {
		panic(err)
	}
	err = tmpl.ExecuteTemplate(os.Stdout, "articles", articles)
	if err != nil {
		panic(err)
	}
}
