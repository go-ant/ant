package main

import (
	"github.com/go-ant/ant/core/server/modules/cache"
	_ "github.com/go-ant/ant/core/server/modules/capabilities"
	"github.com/go-ant/ant/core/server/modules/middleware"
	"github.com/go-ant/ant/core/server/modules/setting"
	"github.com/go-ant/ant/core/server/modules/startup"
	_ "github.com/go-ant/ant/core/server/modules/templates"
	"github.com/go-ant/ant/core/server/routes"
	"github.com/neko-contrib/gzip"
	"github.com/neko-contrib/pongo2"
	"github.com/neko-contrib/sessions"
	"github.com/rocwong/neko"
	"path"
)

func main() {
	app := neko.Classic("GoAnt")

	if setting.EnableGzip {
		app.Use(gzip.Gzip(gzip.DefaultCompression))
	}
	// backend assets path
	app.Static(path.Join(setting.Host.Path, "/goant/assets"), "core/built/assets")
	app.Static(path.Join(setting.Host.Path, setting.API.FilesPath), setting.API.UploadFolder)

	app.Use(middleware.Installer())
	app.Use(pongo2.Renderer(
		pongo2.Options{
			Extension: ".html",
			BaseDir:   "content/themes",
			MultiDir: map[string]string{
				"backend": "core/built",
			},
		}),
	)
	app.Use(sessions.Sessions("goant", sessions.NewCookieStore([]byte("goant"))))
	app.Use(middleware.Contexter())
	app.Use(middleware.UnknowPage())

	cache.NewCache(cache.Options{Store: cache.MemoryStore})

	// create site routers
	routes.Create(app)

	startup.Run()

	app.Run(setting.Host.Addr + ":" + setting.Host.Port)

}
