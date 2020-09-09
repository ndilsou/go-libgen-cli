package nodeutil

import (
	"os"

	"golang.org/x/net/html"
)

//PrintNode displays a html.Node to Stdout
func Print(n *html.Node) {
	html.Render(os.Stdout, n)
}
