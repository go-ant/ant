// ### Page URL Helper
// usage
// {% page_url %}
//
// output the url for the page specified in the current object context
package templates

import (
	"bytes"
	"github.com/go-ant/ant/core/server/modules/setting"
	"gopkg.in/flosch/pongo2.v3"
	"path"
)

type tagPageUrlNode struct {
	path            string
	objectEvaluator pongo2.IEvaluator
}

func (node *tagPageUrlNode) Execute(ctx *pongo2.ExecutionContext, buffer *bytes.Buffer) *pongo2.Error {

	if node.objectEvaluator != nil {
		pageId, _ := node.objectEvaluator.Evaluate(ctx)
		buffer.WriteString(path.Join(node.path, pageId.String()))
	}
	return nil
}

func tagPageUrlParser(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
	node := &tagPageUrlNode{
		path: path.Join(setting.Host.Path, "posts/page"),
	}

	objectEvaluator, err := arguments.ParseExpression()
	if err == nil {
		node.objectEvaluator = objectEvaluator
	}

	return node, nil
}
