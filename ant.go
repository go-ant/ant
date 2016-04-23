package main

import (
	"fmt"
	"path"

	"github.com/neko-contrib/gzip"
	"github.com/neko-contrib/pongo2"
	"github.com/neko-contrib/sessions"
	"github.com/rocwong/neko"

	"github.com/go-ant/ant/core/server/data"
	"github.com/go-ant/ant/core/server/modules/cache"
	_ "github.com/go-ant/ant/core/server/modules/capabilities"
	"github.com/go-ant/ant/core/server/modules/middleware"
	"github.com/go-ant/ant/core/server/modules/setting"
	"github.com/go-ant/ant/core/server/modules/startup"
	_ "github.com/go-ant/ant/core/server/modules/templates"
	"github.com/go-ant/ant/core/server/routes"
)

func main() {
	app := neko.Classic("GoAnt")
	data.GlobalInit()
	if setting.EnableGzip {
		app.Use(gzip.Gzip(gzip.DefaultCompression))
	}
	// backend assets path
	app.Static(path.Join(setting.Host.Path, "/goant/assets"), "built/assets")
	app.Static(path.Join(setting.Host.Path, setting.API.FilesPath), setting.API.UploadFolder)

	app.Use(pongo2.Renderer(
		pongo2.Options{
			Extension: ".html",
			BaseDir:   "content/themes",
			MultiDir: map[string]string{
				"backend": "built",
			},
		}),
	)
	app.Use(sessions.Sessions(setting.Session.Name, sessions.NewCookieStore([]byte(setting.Session.Key))))
	app.Use(middleware.UnknowPage())
	app.Use(middleware.Contexter())

	cache.NewCache(cache.Options{Store: cache.MemoryStore})

	// create site routers
	routes.Create(app)

	startup.Run()
	app.Run(fmt.Sprintf("%s:%s", setting.Host.Addr, setting.Host.Port))

}
