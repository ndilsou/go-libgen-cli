package libgen

import (
	"fmt"
	"log"
	"ndilsou/libgen/pkg/nodeutil"
	"net/http"
	"strconv"

	"golang.org/x/net/html"
)

//Scrape fetches a page from libgen and return the body
func Scrape(domain, query string, page int) (*html.Node, error) {
	u := fmt.Sprintf("http://%s/search.php?req=%s&page=%d", domain, query, page)
	log.Println(u)
	r, err := http.Get(u)
	if err != nil {
		return nil, err
	} else if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Scrape request failed with status %d and reason %s", r.StatusCode, r.Header.Get("Reason"))
	}

	doc, err := html.Parse(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}

	return doc, nil
}

//Book is a container for the useful details about a Libgen book listing
type Book struct {
	ID                                                                int
	Authors, Title, Publisher, Language, Extension, Size, Pages, Year string
}

//ListBooks provides a slice of all the on a page.
func ListBooks(doc *html.Node) ([]*Book, error) {
	t, ok := nodeutil.FindOne(doc, isLinkTable)
	if !ok {
		return nil, fmt.Errorf("ListBooks found no table in the page")
	}
	table := nodeutil.NewTable(t)

	books := []*Book{}
	for _, r := range table.Rows {
		b, err := unmarshallBook(&r)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

func isLinkTable(n *html.Node) bool {
	if n.Type == html.ElementNode && n.Data == "table" {
		if nodeutil.HasAttr(n.Attr, "class", "c") {
			return true
		}
	}
	return false
}

func unmarshallBook(r *nodeutil.Row) (*Book, error) {
	i, err := getTextCol(r, "ID", 4)
	if err != nil {
		return nil, err
	}
	id, err := strconv.Atoi(i)
	if err != nil {
		return nil, fmt.Errorf("unmarshallBook failed to parse ID: %s", err)
	}

	auth, err := getTextCol(r, "Author(s)", 4)
	if err != nil {
		return nil, err
	}

	title, err := getTextCol(r, "Title", 4)
	if err != nil {
		return nil, err
	}

	pub, err := getTextCol(r, "Publisher", 4)
	if err != nil {
		return nil, err
	}

	year, err := getTextCol(r, "Year", 4)
	if err != nil {
		return nil, err
	}

	pages, err := getTextCol(r, "Pages", 4)
	if err != nil {
		return nil, err
	}

	lan, _ := getTextCol(r, "Language", 4)
	if err != nil {
		return nil, err
	}

	size, err := getTextCol(r, "Size", 4)
	if err != nil {
		return nil, err
	}

	ext, err := getTextCol(r, "Extension", 4)
	if err != nil {
		return nil, err
	}

	b := &Book{
		ID:        id,
		Authors:   auth,
		Title:     title,
		Publisher: pub,
		Year:      year,
		Pages:     pages,
		Language:  lan,
		Size:      size,
		Extension: ext,
	}
	return b, nil
}

func getTextCol(r *nodeutil.Row, n string, d int) (string, error) {
	c, prs := (*r)[n]
	if !prs {
		return "", fmt.Errorf("unmarshallBook: missing %s in row", n)
	}
	f, prs := nodeutil.FindOneLimited(c, func(n *html.Node) bool {
		return n.Type == html.TextNode
	}, d)
	if !prs {
		return "", fmt.Errorf("unmarshallBook: failed to extract %s", n)
	}
	return f.Data, nil
}

func getTextColDebug(r *nodeutil.Row, n string, d int) (string, error) {
	c, prs := (*r)[n]
	if !prs {
		return "", fmt.Errorf("unmarshallBook: missing %s in row", n)
	}
	f, prs := nodeutil.FindOne(c, func(n *html.Node) bool {
		res := n.Type == html.TextNode
		log.Println(n.Type, html.TextNode, n.Data, res)
		return res
	})
	if !prs {
		return "", fmt.Errorf("unmarshallBook: failed to extract %s", n)
	}
	return f.Data, nil

}
