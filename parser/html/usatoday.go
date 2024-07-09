package html

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"news-aggregator/constant"
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
	"os"
	"regexp"
	"strings"
	"time"
)

// UsaToday reads and parses an USAToday`s file specified by the path and returns a slice of news.
type UsaToday struct{}

const (
	newsLinkSelector         = "main.gnt_cw div.gnt_m_flm a.gnt_m_flm_a"
	newsDescriptionAttribute = "data-c-br"
	newsDateSelector         = "div.gnt_m_flm_sbt"
)

func (htmlParser UsaToday) Parse(path source.PathToFile, name source.Name) (newsArticles []news.News, parseError error) {
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

	doc.Find(newsLinkSelector).EachWithBreak(func(i int, s *goquery.Selection) bool {
		title := s.Text()
		description, _ := s.Attr(newsDescriptionAttribute)
		link, _ := s.Attr("href")

		if !strings.HasPrefix(link, "http") {
			link = baseURL + link
		}

		date, _ := s.Find(newsDateSelector).Attr("data-c-dt")
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

		newsDate := parsedDate.Format(constant.DateOutputLayout)
		formattedNewsDate, err := time.Parse(constant.DateOutputLayout, newsDate)
		if err != nil {
			parseError = err
			return false
		}

		if formattedNewsDate.Year() < 2000 {
			newsDate = time.Now().Format(constant.DateOutputLayout)
			formattedNewsDate, err = time.Parse(constant.DateOutputLayout, newsDate)
			if err != nil {
				parseError = err
				return false
			}
		}

		newsArticles = append(newsArticles, news.News{
			Title:       news.Title(strings.TrimSpace(title)),
			Description: news.Description(strings.TrimSpace(description)),
			Link:        news.Link(strings.TrimSpace(link)),
			Date:        formattedNewsDate,
			SourceName:  name,
		})

		return true
	})

	if parseError != nil {
		return nil, parseError
	}

	return newsArticles, nil
}
