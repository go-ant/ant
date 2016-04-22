package routes

import (
	"path"

	"github.com/rocwong/neko"

	"github.com/go-ant/ant/core/server/modules/setting"
)

var (
	ignSignIn  = auth(&authOptions{})
	reqSignIn  = auth(&authOptions{SignInRequired: true})
	reqSignOut = auth(&authOptions{SignOutRequired: true})
)

func Create(app *neko.Engine) {
	frontendRoutes(app)
	servicesRoutes(app)
}

func pageUrl(url string) string {
	return path.Join(setting.Host.Path, url)
}
