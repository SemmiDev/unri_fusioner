package sinta

import (
	"errors"
	"fmt"
	"strings"
	uf "unri_fusioner"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

var ErrAuthorNotFound = errors.New("author not found")

func (s *Sinta) ScrapeAuthorProfile(authorID int) (*Author, error) {
	var responseErr error
	url := fmt.Sprintf("%s/authors/profile/%d", s.SintaDomain, authorID)

	c := colly.NewCollector(func(c *colly.Collector) {
		c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"
	})

	author := Author{
		ID: authorID,
		MetaProfile: MetaProfile{
			SintaID: authorID,
		},
		URL: url,
	}

	c.OnHTML("html", func(e *colly.HTMLElement) {
		e.ForEach("h3 > a", func(_ int, el *colly.HTMLElement) {
			author.Name = el.Text
		})

		if author.Name == "" {
			responseErr = ErrAuthorNotFound
			return
		}

		e.ForEach(".meta-profile a:nth-child(1)", func(_ int, el *colly.HTMLElement) {
			affiliationID := uf.SplitAndGet(el.Attr("href"), "/", uf.SplitOpt{
				Last: true,
			})
			affiliationIDInt := uf.CastToInt(affiliationID)

			author.MetaProfile.Affiliation.ID = affiliationIDInt
			author.MetaProfile.Affiliation.Name = strings.TrimSpace(el.Text)
			author.MetaProfile.Affiliation.URL = el.Attr("href")
		})

		e.ForEach(".meta-profile a:nth-child(3)", func(_ int, el *colly.HTMLElement) {
			author.MetaProfile.Department.Name = strings.TrimSpace(el.Text)
		})

		var subjects []string
		e.DOM.Find(".subject-list a").Each(func(_ int, s *goquery.Selection) {
			subjects = append(subjects, s.Text())
		})

		author.Subjects = subjects

		e.DOM.Find(".stat-profile .pr-num").Each(func(i int, s *goquery.Selection) {
			value := uf.CastToInt(s.Text())

			switch i {
			case 0:
				author.Score.Overall = value
			case 1:
				author.Score.ThreeYears = value
			case 2:
				author.Score.Affiliation = value
			case 3:
				author.Score.AffiliationThreeYears = value
			}
		})

		statistics := Statistic{}

		e.DOM.Find(".stat-table > tbody > tr").Each(func(i int, s *goquery.Selection) {
			var rowData []any
			s.Find("td").Each(func(j int, t *goquery.Selection) {
				rowData = append(rowData, t.Text())
			})

			scopus := uf.CastToInt(rowData[1].(string))
			gScholar := uf.CastToInt(rowData[2].(string))
			wos := uf.CastToInt(rowData[3].(string))

			indexer := Indexer{
				Scopus:   scopus,
				GScholar: gScholar,
				WOS:      wos,
			}

			switch rowData[0].(string) {
			case "Article":
				statistics.Articles = indexer
			case "Citation":
				statistics.Citations = indexer
			case "Cited Document":
				statistics.CitedDocs = indexer
			case "H-Index":
				statistics.HIndex = indexer
			case "i10-Index":
				statistics.I10Index = indexer
			case "G-Index":
				statistics.GIndex = indexer
			}
		})

		author.Statistic = statistics
	})

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	if responseErr != nil {
		return nil, responseErr
	}

	return &author, nil
}
