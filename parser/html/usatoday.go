package html

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"news_aggregator/entity/article"
	"news_aggregator/entity/source"
	"os"
	"regexp"
	"strings"
	"time"
)

// UsaToday reads and parses an USAToday`s file specified by the path and returns a slice of articles.
type UsaToday struct {
}

func (htmlParser UsaToday) ParseSource(path source.PathToFile, name source.Name) ([]article.Article, error) {
	file, err := os.Open(string(path))
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return nil, err
	}

	const outputLayout = "2006-01-02"
	baseURL := "https://www.usatoday.com"

	var articles []article.Article
	var parseError error

	doc.Find("main.gnt_cw div.gnt_m_flm a.gnt_m_flm_a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		title := s.Text()
		description, _ := s.Attr("data-c-br")
		link, _ := s.Attr("href")

		if !strings.HasPrefix(link, "http") {
			link = baseURL + link
		}

		dateStr, _ := s.Find("div.gnt_m_flm_sbt").Attr("data-c-dt")
		var parsedDate time.Time
		var err error

		if dateStr != "" {
			re := regexp.MustCompile("[A-Za-z]+\\s\\d{1,2}")
			datePart := re.FindString(dateStr)
			if datePart != "" {
				datePart = fmt.Sprintf("%s %d", datePart, time.Now().Year())
				parsedDate, err = time.Parse("January 2 2006", datePart)
				if err != nil {
					parseError = err
					return false
				}
			}
		}

		formattedDateStr := parsedDate.Format(outputLayout)
		formattedDate, err := time.Parse(outputLayout, formattedDateStr)
		if err != nil {
			parseError = err
			return false
		}

		if formattedDate.Year() < 2000 {
			formattedDateStr = time.Now().Format(outputLayout)
			formattedDate, err = time.Parse(outputLayout, formattedDateStr)
			if err != nil {
				parseError = err
				return false
			}
		}

		articles = append(articles, article.Article{
			Title:       article.Title(strings.TrimSpace(title)),
			Description: article.Description(strings.TrimSpace(description)),
			Link:        article.Link(strings.TrimSpace(link)),
			Date:        formattedDate,
			SourceName:  name,
		})

		return true
	})

	if parseError != nil {
		return nil, parseError
	}

	return articles, nil
}
