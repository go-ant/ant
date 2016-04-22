// # User Helper
// usage
// {% user %}{{user.Name}}{% enduser %}
// {% user userId %}{{user.Name}}{% enduser %}
//
// returns user data
// defaults to get owner user
package templates

import (
	"bytes"

	"gopkg.in/flosch/pongo2.v3"

	"github.com/go-ant/ant/core/server/models"
	"github.com/go-ant/ant/core/server/modules/utils"
)

type tagUserNode struct {
	key             string
	end             string
	objectEvaluator pongo2.IEvaluator
	bodyWrapper     *pongo2.NodeWrapper
	Opts            *models.Options
}

func (node *tagUserNode) Execute(ctx *pongo2.ExecutionContext, buffer *bytes.Buffer) *pongo2.Error {
	var user *models.User
	var errApi models.ApiErr
	var uid *pongo2.Value

	newCtx := pongo2.NewChildExecutionContext(ctx)
	newCtx.Private[keyNoUser] = false

	opts := &models.Options{}
	if node.objectEvaluator != nil {
		uid, _ = node.objectEvaluator.Evaluate(newCtx)
	}
	if uid != nil && uid.IsInteger() {
		user, errApi = models.GetUserById(utils.ToUint32(uid.String()), opts)
	} else {
		// default to get owner user
		opts.GormAdp = &models.GormAdapter{
			Joins: "role",
			Map:   map[string]interface{}{"roles.slug": "owner"},
		}
		user, errApi = models.GetUser(opts)
	}

	if !errApi.IsSuccess() {
		newCtx.Private[keyNoUser] = true
	} else {
		newCtx.Private[node.key] = user
	}
	node.bodyWrapper.Execute(newCtx, buffer)
	return nil
}

func tagUserParser(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
	node := &tagUserNode{
		key: "user",
		end: "enduser",
	}

	objectEvaluator, err := arguments.ParseExpression()
	if err == nil {
		node.objectEvaluator = objectEvaluator
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
