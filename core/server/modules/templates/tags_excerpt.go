// # Excerpt Helper
// usage
// {% excerpt %}
// {% excerpt 50 %}
// {% excerpt 50 "words" %}
// {% excerpt 50 "characters" %}
//
// attempts to remove all html from the string, and then shortens the result according to the provided option.
// defaults to words="50"
package templates

import (
	"bytes"
	"strings"

	"gopkg.in/flosch/pongo2.v3"

	"github.com/go-ant/ant/core/server/models"
	"github.com/go-ant/ant/core/server/modules/utils"
)

type tagExcerptNode struct {
	excerptType     string
	excerptLength   int
	objectEvaluator pongo2.IEvaluator
}

func (node *tagExcerptNode) Execute(ctx *pongo2.ExecutionContext, buffer *bytes.Buffer) *pongo2.Error {

	var post *models.Post
	if ctx.Private["post"] != nil {
		post = ctx.Private["post"].(*pongo2.Value).Interface().(*models.Post)
	}

	if post == nil && ctx.Public["post"] != nil {
		post = ctx.Public["post"].(*models.Post)
	}

	if post != nil {
		excerpt := ""
		if node.excerptType != "words" {
			excerpt = utils.StripTagsFromHtml(post.Html)
			runes := []rune(excerpt)
			if len(runes) > node.excerptLength {
				excerpt = string(runes[:node.excerptLength])
			}
		} else {
			excerpt = utils.StripTagsFromHtml(post.Html)
			words := strings.Fields(excerpt)
			if len(words) > node.excerptLength {
				excerpt = strings.Join(words[:node.excerptLength], " ")
			}
		}
		buffer.WriteString(excerpt)
	}

	return nil
}

func tagExcerptParser(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
	node := &tagExcerptNode{
		excerptType:   "words",
		excerptLength: 50,
	}

	formatToken := arguments.MatchType(pongo2.TokenNumber)

	if formatToken != nil {
		node.excerptLength = utils.ToInt(formatToken.Val)
	}

	formatToken = arguments.MatchType(pongo2.TokenString)
	if formatToken != nil {
		node.excerptType = formatToken.Val
	}

	if node.excerptLength == 0 {
		node.excerptLength = 50
	}

	return node, nil
}
