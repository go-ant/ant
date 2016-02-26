package routes

import (
	"github.com/go-ant/ant/core/server/modules/middleware"
	"github.com/go-ant/ant/core/server/services"
	"github.com/rocwong/neko"
)

func servicesRoutes(app *neko.Engine) {

	app.Group(pageUrl("/api"), func(router *neko.RouterGroup) {

		router.Group("/backend", func(router *neko.RouterGroup) {
			// posts services
			router.Group("/posts", func(router *neko.RouterGroup) {
				router.POST("", services.PostCreate)
				router.PUT("/:post_id", services.PostEdit)
				router.GET("", services.PostList)
				router.GET("/:post_id", services.PostInfo)
				router.DELETE("/:post_id", services.PostDelete)
			})

			// users services
			router.Group("/users", func(router *neko.RouterGroup) {
				router.GET("", services.UserList)
				router.GET("/info/:account", services.UserInfo)

				router.POST("/info", services.UserCreate)
				router.PUT("/info/:user_id", services.UserEdit)
				router.PUT("/password/:user_id", services.UserChangePassword)
			})

			// roles services
			router.Group("/roles", func(router *neko.RouterGroup) {
				router.POST("", services.RoleCreate)
				router.PUT("/:role_id", services.RoleEdit)
				router.GET("", services.RoleList)
				router.GET("/:role_id", services.RoleInfo)
				router.DELETE("/:role_id", services.RoleDelete)
			})

			// permissions services
			router.Group("/permissions", func(router *neko.RouterGroup) {
				router.GET("", services.PermissionList)
			})

			// themes services
			router.Group("/themes", func(router *neko.RouterGroup) {
				router.GET("", services.ThemeList)
			})

			// settings services
			router.Group("/settings", func(router *neko.RouterGroup) {
				router.PUT("", services.SettingEdit)
				router.GET("", services.SettingList)
			})

			// upload services
			router.POST("/upload", services.Upload)

			router.GET("/app", services.InitialApp)

		}, middleware.RequireLogin(true))

		router.POST("/signin", services.UserLogin)
		router.POST("/setup", services.AppSetup)

	})

}
