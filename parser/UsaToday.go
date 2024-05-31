package parser

import (
	"NewsAggregator/entity"
	"NewsAggregator/entity/article"
	"NewsAggregator/entity/source"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"strings"
	"time"
)

// UsaToday reads and parses an USAToday`s file specified by the path and returns a slice of articles.
type UsaToday struct {
}

func (htmlParser UsaToday) ParseSource(path source.PathToFile) []article.Article {
	file, err := os.Open(string(path))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		log.Fatal(err)
	}

	const layout = "January 2, 2006"
	const outputLayout = "2006-01-02"
	baseURL := "https://www.usatoday.com"

	var articles []article.Article
	doc.Find("main.gnt_cw div.gnt_m_flm a.gnt_m_flm_a").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		description, _ := s.Attr("data-c-br")
		link, _ := s.Attr("href")

		if !strings.HasPrefix(link, "http") {
			link = baseURL + link
		}

		dateStr, _ := s.Find("div.gnt_m_flm_sbt").Attr("data-c-dt")
		parsedDate, err := time.Parse(layout, dateStr)
		if err != nil {
			log.Println("Error parsing date:", err)
			parsedDate = time.Time{}
		}

		formattedDateStr := parsedDate.Format(outputLayout)
		formattedDate, err := time.Parse(outputLayout, formattedDateStr)
		if err != nil {
			log.Println("Error formatting date:", err)
		}

		articles = append(articles, article.Article{
			Id:          entity.Id(i + 1),
			Title:       article.Title(strings.TrimSpace(title)),
			Description: article.Description(strings.TrimSpace(description)),
			Link:        article.Link(strings.TrimSpace(link)),
			Date:        formattedDate,
		})
	})

	return articles
}
