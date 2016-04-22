// # Post Helper
// usage
// {% posts postNum %}{{user.Name}}{% endposts %}
// {% posts postNum "{'Limit': '15', 'Columns': ['id'], 'Include': 'user'}" %}{{user.Name}}{% endposts %}
//
// return posts data
package templates

import (
	"bytes"
	"encoding/json"
	"math"
	"strings"
	"time"

	"gopkg.in/flosch/pongo2.v3"

	"github.com/go-ant/ant/core/server/models"
	"github.com/go-ant/ant/core/server/modules/utils"
)

type tagPostsNode struct {
	key             string
	end             string
	objectEvaluator pongo2.IEvaluator
	bodyWrapper     *pongo2.NodeWrapper
	Parms           *tagPostsParams
	Opts            *models.Options
}

type tagPostsParams struct {
	Limit      uint32
	Columns    []string
	OrderByAsc bool
	Include    string
}

func (node *tagPostsNode) Execute(ctx *pongo2.ExecutionContext, buffer *bytes.Buffer) *pongo2.Error {
	newCtx := pongo2.NewChildExecutionContext(ctx)
	newCtx.Private[keyNoPosts] = false

	node.Opts.GormAdp.Query = "status = ? and page = false and published_at < ?"
	node.Opts.GormAdp.Args = []interface{}{models.PostStatusPublished, time.Now()}

	if node.objectEvaluator != nil {
		page, _ := node.objectEvaluator.Evaluate(newCtx)
		node.Opts.Page = utils.ToUint32(page.String())
	}

	if node.Parms.Limit == 0 {
		node.Opts.Limit = models.GetAppSetting().PostsPerPage
	}

	page := &page{}
	list, recordCount, errApi := models.GetPosts(node.Opts)
	page.Count = uint32(math.Ceil(float64(recordCount) / float64(node.Opts.Limit)))

	if page.Count == 0 || node.Opts.Page <= 1 {
		page.Prev = 0
	} else {
		page.Prev = node.Opts.Page - 1
	}

	if page.Count == 0 || node.Opts.Page == page.Count {
		page.Next = 0
	} else {
		page.Next = node.Opts.Page + 1
	}

	newCtx.Private[keyPostCount] = recordCount
	newCtx.Private["page"] = page

	if !errApi.IsSuccess() || len(list) == 0 {
		newCtx.Private[keyNoPosts] = true
		node.bodyWrapper.Execute(newCtx, buffer)
		return nil
	}

	newCtx.Private[node.key] = list
	node.bodyWrapper.Execute(newCtx, buffer)

	return nil
}

func tagPostsParser(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
	node := &tagPostsNode{
		key: "posts",
		end: "endposts",
	}

	objectEvaluator, err := arguments.ParseExpression()
	if err == nil {
		node.objectEvaluator = objectEvaluator
	}

	formatToken := arguments.MatchType(pongo2.TokenString)
	if formatToken != nil {
		json.Unmarshal([]byte(strings.Replace(formatToken.Val, "'", "\"", -1)), &node.Parms)
	}
	if node.Parms == nil {
		node.Parms = &tagPostsParams{Limit: models.GetAppSetting().PostsPerPage}
	}

	node.Opts = &models.Options{
		Limit: node.Parms.Limit,
		Page:  1,
		GormAdp: &models.GormAdapter{
			OrderBy: "published_at desc, id desc",
		},
	}

	if !utils.IsEmpty(node.Parms.Include) {
		node.Opts.Include = node.Parms.Include
	}

	if len(node.Parms.Columns) > 0 {
		node.Opts.GormAdp.Columns = node.Parms.Columns
		if node.Opts.IsInclude("user") && !utils.StringInSlice("author_id", node.Opts.GormAdp.Columns) {
			node.Opts.GormAdp.Columns = append(node.Opts.GormAdp.Columns, "author_id")
		}
	} else {
		node.Opts.GormAdp.Columns = []string{"id", "title", "slug", "html", "cover", "published_at", "author_id"}
	}

	if node.Parms.OrderByAsc {
		node.Opts.GormAdp.OrderBy = "published_at asc, id desc"
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
