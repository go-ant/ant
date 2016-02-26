// # Asset Helper
// usage
// {% asset "css/screen.css" %}
// {% asset "css/screen.css" goant=true %}
//
// output the path to the specified asset. the goant flag output the asset path for the goant admin
package templates

import (
	"bytes"
	"github.com/go-ant/ant/core/server/models"
	"github.com/go-ant/ant/core/server/modules/setting"
	"gopkg.in/flosch/pongo2.v3"
	"path"
)

type tagAssetNode struct {
	filePath string
}

func (node *tagAssetNode) Execute(ctx *pongo2.ExecutionContext, buffer *bytes.Buffer) *pongo2.Error {
	buffer.WriteString(node.filePath)
	return nil
}

func tagAssetParser(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
	node := &tagAssetNode{}

	formatToken := arguments.MatchType(pongo2.TokenString)
	if formatToken != nil {
		appSetting := models.GetAppSetting()
		node.filePath = setting.Host.Path + "/assets/"
		node.filePath = path.Join(node.filePath, formatToken.Val)

		if arguments.Remaining() > 0 {
			keyToken := arguments.MatchType(pongo2.TokenIdentifier)
			if keyToken != nil || keyToken.Val == "goant" {
				if arguments.Match(pongo2.TokenSymbol, "=") != nil {
					valueToken := arguments.MatchType(pongo2.TokenKeyword)
					if valueToken != nil || valueToken.Val == "true" {
						node.filePath = path.Join(setting.Host.Path, formatToken.Val)
					}
				}
			}
		}
		if path.Ext(node.filePath) == "" {
			node.filePath += "/"
		} else if appSetting != nil {
			node.filePath = node.filePath + "?v=" + appSetting.Version
		} else {
			node.filePath = node.filePath
		}
	}

	return node, nil
}
