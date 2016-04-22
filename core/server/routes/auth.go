package routes

import (
	"path"
	"strings"

	"github.com/rocwong/neko"

	"github.com/go-ant/ant/core/server/models"
	"github.com/go-ant/ant/core/server/modules/middleware"
	"github.com/go-ant/ant/core/server/modules/setting"
)

type authOptions struct {
	IsApi           bool
	SignInRequired  bool
	SignOutRequired bool
}

func auth(options *authOptions) neko.HandlerFunc {
	return func(ctx *neko.Context) {
		// cannot view any page before installation.
		if !setting.InstallLock {
			ctx.Redirect(setting.INSTALLER_URL)
			ctx.Abort()
			return
		}

		if options.SignOutRequired && middleware.Context.IsSigned && ctx.Req.RequestURI != "/" {
			ctx.Redirect(setting.Host.Path)
			return
		}

		if options.SignInRequired && !middleware.Context.IsSigned {
			if isApiPath(ctx.Req.URL.Path) {
				ctx.Json(models.RestApi{Error: models.ApiMsg.NeedToSignIn})
			} else {
				ctx.Redirect(setting.Host.Path + "/goant/login")
			}
			ctx.Abort()
			return
		}
	}
}

func isApiPath(url string) bool {
	return strings.HasPrefix(url, path.Join(setting.Host.Path, "/api/"))
}