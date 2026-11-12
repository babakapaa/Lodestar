package notes

import (
	"bytes"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

func pageFields(body []byte, f SourceFields) SourceFields {
	meta := map[string][]string{}
	var title string
	z := html.NewTokenizer(bytes.NewReader(body))
	inTitle := false
	for {
		switch z.Next() {
		case html.ErrorToken:
			return applyPageMeta(f, meta, title)
		case html.StartTagToken, html.SelfClosingTagToken:
			name, hasAttr := z.TagName()
			switch string(name) {
			case "meta":
				key, value := metaPair(z, hasAttr)
				if key != "" && value != "" {
					meta[key] = append(meta[key], value)
				}
			case "title":
				inTitle = title == ""
			case "body":
				return applyPageMeta(f, meta, title)
			}
		case html.TextToken:
			if inTitle {
				title = strings.TrimSpace(string(z.Text()))
				inTitle = false
			}
		case html.EndTagToken:
			inTitle = false
		}
	}
}

func metaPair(z *html.Tokenizer, more bool) (string, string) {
	var key, value string
	for more {
		var k, v []byte
		k, v, more = z.TagAttr()
		switch string(k) {
		case "name", "property":
			key = strings.ToLower(string(v))
		case "content":
			value = strings.TrimSpace(string(v))
		}
	}
	return key, value
}

func applyPageMeta(f SourceFields, meta map[string][]string, title string) SourceFields {
	pick := func(keys ...string) string {
		for _, k := range keys {
			if v := meta[k]; len(v) > 0 {
				return v[0]
			}
		}
		return ""
	}
	f.Title = firstNonEmpty(pick("citation_title", "dc.title", "og:title"), title)
	f.Authors = meta["citation_author"]
	if len(f.Authors) == 0 {
		f.Authors = slices.Concat(meta["dc.creator"], meta["author"])
	}
	f.Container = pick("citation_journal_title", "citation_conference_title", "og:site_name")
	f.Publisher = pick("citation_publisher", "dc.publisher")
	f.Volume = pick("citation_volume")
	f.Issue = pick("citation_issue")
	if fp := pick("citation_firstpage"); fp != "" {
		f.Pages = fp
		if lp := pick("citation_lastpage"); lp != "" {
			f.Pages += "–" + lp
		}
	}
	if date := pick("citation_publication_date", "citation_date", "dc.date", "article:published_time"); len(date) >= 4 {
		f.Year = date[:4]
	}
	if doi := pick("citation_doi", "dc.identifier"); strings.HasPrefix(NormalizeDOI(doi), "10.") {
		f.DOI = NormalizeDOI(doi)
	}
	f.ISBN = pick("citation_isbn")
	if pick("citation_journal_title") != "" {
		f.Kind = "article"
	}
	return f
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
