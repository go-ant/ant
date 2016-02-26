// # Code Injection Helper
// usage
// {% ant_head %}
// {% ant_foot %}
//
// outputs scripts and other assets at the top/bottom of a goant theme
package templates

import (
	"bytes"
	"github.com/go-ant/ant/core/server/models"
	"gopkg.in/flosch/pongo2.v3"
	"strings"
)

// header code injection
type tagHeadCodeInjectionNode struct {
	Code string
}

func (node *tagHeadCodeInjectionNode) Execute(ctx *pongo2.ExecutionContext, buffer *bytes.Buffer) *pongo2.Error {
	appSetting := models.GetAppSetting()
	buffer.WriteString(node.Code + "\n" + appSetting.AntHead)
	return nil
}

func tagHeadCodeInjectionParser(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
	appSetting := models.GetAppSetting()
	headCode := make([]string, 0)
	headCode = append(headCode, `<meta name="referrer" content="origin" />`)
	headCode = append(headCode, `<meta name="generator" content="GoAnt `+appSetting.Version+`" />`)

	node := &tagHeadCodeInjectionNode{
		Code: strings.Join(headCode, "\n"),
	}
	return node, nil
}

// footer code injection
type tagFootCodeInjectionNode struct {
	Code string
}

func (node *tagFootCodeInjectionNode) Execute(ctx *pongo2.ExecutionContext, buffer *bytes.Buffer) *pongo2.Error {
	appSetting := models.GetAppSetting()
	buffer.WriteString(node.Code + "\n" + appSetting.AntFoot)
	return nil
}

func tagFootCodeInjectionParser(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
	node := &tagFootCodeInjectionNode{}
	return node, nil
}
