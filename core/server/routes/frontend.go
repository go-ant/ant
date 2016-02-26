package routes

import (
	"github.com/go-ant/ant/core/server/controllers"
	"github.com/go-ant/ant/core/server/modules/middleware"
	"github.com/rocwong/neko"
)

func frontendRoutes(app *neko.Engine) {

	app.GET(pageUrl("/"), controllers.Index)
	app.GET(pageUrl("/posts/page/:page"), controllers.Index)
	app.GET(pageUrl("/page/:slug"), controllers.Page)
	app.GET(pageUrl("/post/:slug"), controllers.Post)

	app.GET(pageUrl("/assets/*filepath"), controllers.AssetsHandler)

	app.Group(pageUrl("/goant"), func(router *neko.RouterGroup) {
		router.GET("/setup", controllers.Installer)
		router.GET("/login", controllers.Login)
		router.GET("/logout", controllers.Logout)

		router.GET("/", middleware.RequireLogin(false), controllers.Admin)
	})
}
