package middleware

import (
	"path"

	"github.com/rocwong/neko"

	"github.com/go-ant/ant/core/server/models"
)

const (
	SESSION_USER_ID = "GoAnt_User_Id"
)

var (
	Context    *ContextHelper
	fileFilter map[string]bool = map[string]bool{
		".css":  true,
		".js":   true,
		".jpg":  true,
		".jpeg": true,
		".gif":  true,
		".png":  true,
		".ico":  true,
		".json": true,
		".eot":  true,
		".svg":  true,
		".ttf":  true,
		".woff": true,
	}
)

type ContextHelper struct {
	User     *models.User
	IsSigned bool
}

// Contexter initializes a classic context for a request.
func Contexter() neko.HandlerFunc {
	return func(ctx *neko.Context) {
		if !fileFilter[path.Ext(ctx.Req.URL.Path)] && models.HasEngine {
			Context = new(ContextHelper)
			uid := ctx.Session.Get(SESSION_USER_ID)
			if uid != nil {
				Context.User, _ = models.GetUserById(uid.(uint32), nil)
				Context.IsSigned = true
			}
		}
	}
}
