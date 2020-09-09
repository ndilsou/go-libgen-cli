package nodeutil

import (
	"fmt"

	"golang.org/x/net/html"
)

// Table is a container for a 2D html <table>. Make access to element nodes easier
//  Expects each row to have the same number of columns.
type Table struct {
	// Name of the columns of the table
	Headers []string
	// Map of the <td> nodes in the row
	Rows []Row
}

//Row is a container for a html tr with easy access to the tr
type Row map[string]*html.Node

// Shape returns the number of rows and columns in the table.
func (t *Table) Shape() (int, int) { return len(t.Rows), len(t.Headers) }

// NewTable creates a new Table
func NewTable(t *html.Node) *Table {
	h := parseHeader(t)
	r := parseRows(t, h)
	return &Table{Headers: h, Rows: r}
}

func parseHeader(n *html.Node) []string {
	var headers []string

	tr, _ := FindOne(n, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == "tr"
	})
	tds := FindLimited(tr, func(n *html.Node) bool { return n.Type == html.ElementNode && n.Data == "td" }, 2)

	for _, td := range tds {
		n, ok := FindOne(td, func(n *html.Node) bool { return n.Type == html.TextNode })
		if !ok {
			fmt.Println("Failed to find a column name")
		}
		headers = append(headers, n.Data)
	}
	return headers
}

func parseRows(n *html.Node, header []string) []Row {
	trs := FindLimited(n, func(n *html.Node) bool { return n.Type == html.ElementNode && n.Data == "tr" }, 3)
	var rs []Row
	for _, tr := range trs[1:] {
		r := make(map[string]*html.Node)
		tds := FindLimited(tr, func(n *html.Node) bool { return n.Type == html.ElementNode && n.Data == "td" }, 2)
		for i, td := range tds {
			r[header[i]] = td
		}
		rs = append(rs, r)
	}
	return rs
}
