package main

import (
	"fmt"
	"log"
	"ndilsou/libgen/pkg/libgen"
	"ndilsou/libgen/pkg/nodeutil"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	query := strings.Join(os.Args[1:], "+")
	doc, err := libgen.Scrape("libgen.li", query, 1)
	if err != nil {
		log.Fatal(err)
	}
	books, err := libgen.ListBooks(doc)
	if err != nil {
		log.Fatal(err)
	}
	for _, b := range books {
		fmt.Println(b)
	}
}

func isLinkTable(n *html.Node) bool {
	if n.Type == html.ElementNode && n.Data == "table" {
		if nodeutil.HasAttr(n.Attr, "class", "c") {
			return true
		}
	}
	return false
}

func isBookLink(n *html.Node) bool {
	if n.Type == html.ElementNode && n.Data == "tr" {
		return true
	}
	return false
}
