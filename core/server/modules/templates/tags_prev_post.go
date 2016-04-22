// ### Previous Post helper
// usage
// {% prev_post %}<a href ="{{prev_post.GetURL()}}>previous post</a>{% endprev_post %}
package templates

import (
	"bytes"

	"gopkg.in/flosch/pongo2.v3"

	"github.com/go-ant/ant/core/server/models"
)

type tagPrevPostNode struct {
	key             string
	end             string
	objectEvaluator pongo2.IEvaluator
	bodyWrapper     *pongo2.NodeWrapper
	Opts            *models.Options
}

func (node *tagPrevPostNode) Execute(ctx *pongo2.ExecutionContext, buffer *bytes.Buffer) *pongo2.Error {
	if ctx.Public["post"] != nil {
		newCtx := pongo2.NewChildExecutionContext(ctx)
		post := newCtx.Public["post"].(*models.Post)

		node.Opts.GormAdp.Query = "status = ? and page = false and published_at < ?"
		node.Opts.GormAdp.Args = []interface{}{models.PostStatusPublished, post.PublishedAt}
		nextPost, errApi := models.GetPost(node.Opts)
		if errApi.IsSuccess() {
			newCtx.Private[node.key] = nextPost
			node.bodyWrapper.Execute(newCtx, buffer)
		}
	}

	return nil
}

func tagPrevPostParser(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
	node := &tagPrevPostNode{
		key: "prev_post",
		end: "endprev_post",
	}

	node.Opts = &models.Options{
		Limit: 1,
		GormAdp: &models.GormAdapter{
			Columns: []string{"id", "title", "slug", "html", "cover", "published_at"},
			OrderBy: "published_at desc, id desc",
		},
	}

	// Body wrapping
	wrapper, endargs, err := doc.WrapUntilTag(node.end)
	if err != nil {
		return nil, err
	}
	node.bodyWrapper = wrapper

	if endargs.Count() > 0 {
		return nil, endargs.Error("Arguments not allowed here.", nil)
	}

	return node, nil
}
