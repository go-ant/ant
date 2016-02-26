// # Post Class Helper
// usage
// {% post_class %}
//
// output classes for the `post/page` page
package templates

import (
	"bytes"
	"github.com/go-ant/ant/core/server/models"
	"gopkg.in/flosch/pongo2.v3"
	"strings"
)

type tagPostClassNode struct{}

func (node *tagPostClassNode) Execute(ctx *pongo2.ExecutionContext, buffer *bytes.Buffer) *pongo2.Error {

	if ctx.Public["post"] != nil {
		classes := []string{"post"}
		post := ctx.Public["post"].(*models.Post)

		if post.Featured {
			classes = append(classes, "featured")
		}

		if post.Page {
			classes = append(classes, "page")
		}

		buffer.WriteString(strings.Join(classes, " "))
	}

	return nil
}

func tagPostClassParser(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
	node := &tagPostClassNode{}

	return node, nil
}
