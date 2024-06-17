package html

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"news_aggregator/constant"
	"news_aggregator/entity/article"
	"news_aggregator/entity/source"
	"os"
	"regexp"
	"strings"
	"time"
)

// UsaToday reads and parses an USAToday`s file specified by the path and returns a slice of articles.
type UsaToday struct{}

func (htmlParser UsaToday) ParseSource(path source.PathToFile, name source.Name) (articles []article.Article, parseError error) {
	file, err := os.Open(string(path))
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			parseError = fmt.Errorf("cannot close HTML file: %w", cerr)
		}
	}()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return nil, err
	}

	baseURL := "https://www.usatoday.com"

	doc.Find("main.gnt_cw div.gnt_m_flm a.gnt_m_flm_a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		title := s.Text()
		description, _ := s.Attr("data-c-br")
		link, _ := s.Attr("href")

		if !strings.HasPrefix(link, "http") {
			link = baseURL + link
		}

		date, _ := s.Find("div.gnt_m_flm_sbt").Attr("data-c-dt")
		var parsedDate time.Time

		if date != "" {
			re := regexp.MustCompile("[A-Za-z]+\\s\\d{1,2}")
			datePart := re.FindString(date)
			if datePart != "" {
				datePart = fmt.Sprintf("%s %d", datePart, time.Now().Year())
				parsedDate, err = time.Parse("January 2 2006", datePart)
				if err != nil {
					parseError = err
					return false
				}
			}
		}

		formattedDateStr := parsedDate.Format(constant.DateOutputLayout)
		formattedDate, err := time.Parse(constant.DateOutputLayout, formattedDateStr)
		if err != nil {
			parseError = err
			return false
		}

		if formattedDate.Year() < 2000 {
			formattedDateStr = time.Now().Format(constant.DateOutputLayout)
			formattedDate, err = time.Parse(constant.DateOutputLayout, formattedDateStr)
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
