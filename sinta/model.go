package sinta

type Affiliation struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Department struct {
	Name string `json:"name"`
}

type MetaProfile struct {
	SintaID     int         `json:"sintaID"`
	Affiliation Affiliation `json:"affiliation"`
	Department  Department  `json:"department"`
}

type Indexer struct {
	Scopus   int `json:"scopus"`
	GScholar int `json:"GScholar"`
	WOS      int `json:"WOS"`
}

type Statistic struct {
	Articles  Indexer `json:"articles"`
	Citations Indexer `json:"citations"`
	CitedDocs Indexer `json:"citedDocs"`
	HIndex    Indexer `json:"HIndex"`
	I10Index  Indexer `json:"i10Index"`
	GIndex    Indexer `json:"GIndex"`
}

type Score struct {
	Overall               int `json:"overall"`
	ThreeYears            int `json:"threeYears"`
	Affiliation           int `json:"affiliation"`
	AffiliationThreeYears int `json:"affiliationThreeYears"`
}

type Author struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	URL         string      `json:"url"`
	MetaProfile MetaProfile `json:"metaProfile"`
	Subjects    []string    `json:"subjects"`
	Score       Score       `json:"score"`
	Statistic   Statistic   `json:"statistic"`
}
