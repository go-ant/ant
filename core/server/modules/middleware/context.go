package middleware

import (
	"github.com/go-ant/ant/core/server/models"
	"github.com/go-ant/ant/core/server/modules/setting"
	"github.com/rocwong/neko"
	"path"
)

const (
	SESSION_USER_ID = "GoAnt_User_Id"
)

var (
	Context    *AssistantContext
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

type AssistantContext struct {
	User *models.User
}

// Contexter initializes a classic context for a request.
func Contexter() neko.HandlerFunc {
	return func(ctx *neko.Context) {
		if !fileFilter[path.Ext(ctx.Req.URL.Path)] {
			Context = new(AssistantContext)
			uid := ctx.Session.Get(SESSION_USER_ID)
			if uid != nil {
				Context.User, _ = models.GetUserById(uid.(uint32), nil)
			}
		}
	}
}

func RequireLogin(isApi bool) neko.HandlerFunc {
	return func(ctx *neko.Context) {
		if Context.User == nil {
			if isApi {
				ctx.Json(models.RestApi{Error: models.ApiMsg.NeedToSignIn})
			} else {
				ctx.Redirect(setting.Host.Path + "/goant/login")
			}
			ctx.Abort()
		}
	}
}
