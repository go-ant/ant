// # Body Class Helper
// usage
// {% body_class %}
//
// output classes for the body element
package templates

import (
	"bytes"
	"gopkg.in/flosch/pongo2.v3"
	"strings"
)

type tagBodyClassNode struct{}

func (node *tagBodyClassNode) Execute(ctx *pongo2.ExecutionContext, buffer *bytes.Buffer) *pongo2.Error {
	classes := []string{}
	tpl := ctx.Public["tpl"].(string)
	switch tpl {
	case "home":
		classes = append(classes, "tpl-home")
		if ctx.Public["pageNum"].(uint32) > 1 {
			classes = append(classes, "paged")
		}
	case "post":
		classes = append(classes, "tpl-post")
	case "page":
		classes = append(classes, "tpl-page")
	case "404":
		classes = append(classes, "tpl-404")
	default:

	}
	buffer.WriteString(strings.Join(classes, " "))
	return nil
}

func tagBodyClassParser(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
	node := &tagBodyClassNode{}

	return node, nil
}
