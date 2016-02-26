package routes

import (
	"github.com/go-ant/ant/core/server/modules/setting"
	"github.com/rocwong/neko"
	"path"
)

func pageUrl(url string) string {
	return path.Join(setting.Host.Path, url)
}

func Create(app *neko.Engine) {
	frontendRoutes(app)
	servicesRoutes(app)
}
