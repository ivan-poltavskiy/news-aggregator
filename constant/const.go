package constant

import "news-aggregator/entity/source"

const DateOutputLayout = "2006-01-02"

var BbcSource = source.Source{Name: "bbc", PathToFile: "resources/bbc-world-category-19-05-24.xml", SourceType: "RSS"}
var NbcSource = source.Source{Name: "nbc", PathToFile: "resources/nbc-news.json", SourceType: "JSON"}
var AbcSource = source.Source{Name: "abc", PathToFile: "resources/abcnews-international-category-19-05-24.xml", SourceType: "RSS"}
var WashingtonSource = source.Source{Name: "washington", PathToFile: "resources/washingtontimes-world-category-19-05-24.xml", SourceType: "RSS"}
var UsaTodaySource = source.Source{Name: "usatoday", PathToFile: "resources/usatoday-world-news.html", SourceType: "UsaToday"}
