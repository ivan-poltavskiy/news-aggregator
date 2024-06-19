package filter

import (
	"fmt"
	"github.com/kljensen/snowball"
	"news_aggregator/entity/article"
	"strings"
)

// ByKeyword filters the slice of articles by provided keyword and returns
// the slice of matching articles.
type ByKeyword struct {
	Keywords []string
}

// Filter filters the incoming collector list from different sources by keywords.
func (f ByKeyword) Filter(articles []article.Article) []article.Article {
	var matchingArticles []article.Article
	for _, keyword := range f.Keywords {
		stemmedKeyword, err := snowball.Stem(keyword, "english", true)
		if err != nil {
			fmt.Printf("Error stemming keyword %s: %v\n", keyword, err)
			continue
		}
		matchingArticles = append(matchingArticles, filterNewsByKeyword(stemmedKeyword, articles)...)
	}
	return matchingArticles
}

// filterNewsByKeyword filters the incoming collector list by keyword and returns the filtered list.
func filterNewsByKeyword(keyword string, articles []article.Article) []article.Article {
	var matchingArticles []article.Article

	for _, a := range articles {
		titleStemmed, err := snowball.Stem(strings.ToLower(string(a.Title)), "english", true)
		if err != nil {
			fmt.Printf("Error stemming title: %v\n", err)
			continue
		}
		descriptionStemmed, err := snowball.Stem(strings.ToLower(string(a.Description)), "english", true)
		if err != nil {
			fmt.Printf("Error stemming description: %v\n", err)
			continue
		}

		if strings.Contains(titleStemmed, keyword) || strings.Contains(descriptionStemmed, keyword) {
			matchingArticles = append(matchingArticles, a)
		}
	}

	if len(matchingArticles) == 0 {
		fmt.Println("No matches found for this keyword.")
	}

	return matchingArticles
}
