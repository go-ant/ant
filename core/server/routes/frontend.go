package routes

import (
	"github.com/rocwong/neko"

	"github.com/go-ant/ant/core/server/controllers"
	"github.com/go-ant/ant/core/server/modules/setting"
)

func frontendRoutes(app *neko.Engine) {

	app.GET(setting.INSTALLER_URL, controllers.Installer)
	app.GET(pageUrl("/"), ignSignIn, controllers.Index)
	app.GET(pageUrl("/posts/page/:page"), ignSignIn, controllers.Index)
	app.GET(pageUrl("/page/:slug"), ignSignIn, controllers.Page)
	app.GET(pageUrl("/post/:slug"), ignSignIn, controllers.Post)

	app.GET(pageUrl("/assets/*filepath"), ignSignIn, controllers.AssetsHandler)

	app.Group(pageUrl("/goant"), func(router *neko.RouterGroup) {
		router.GET("/login", reqSignOut, controllers.Login)
		router.GET("/logout", controllers.Logout)

		router.GET("/", reqSignIn, controllers.Admin)
	})
}
