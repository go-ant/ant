package controllers

import (
	"github.com/go-ant/ant/core/server/models"
	"github.com/go-ant/ant/core/server/modules/middleware"
	"github.com/go-ant/ant/core/server/modules/setting"
	"github.com/go-ant/ant/core/server/modules/utils"
	"github.com/rocwong/neko"
	"net/http"
	"path/filepath"
)

// themeView returns view path
func themeViews(view string) string {
	appSetting := models.GetAppSetting()
	return appSetting.ActiveTheme + "/" + view
}

func Index(ctx *neko.Context) {
	appSetting := models.GetAppSetting()

	pageNum := utils.ToUint32(ctx.Params.ByGet("page"))
	if pageNum == 0 {
		pageNum = 1
	}

	ctx.Render(themeViews("index"), neko.JSON{
		"tpl":     "home",
		"app":     appSetting,
		"meta":    map[string]string{"Title": appSetting.Title, "Description": appSetting.Description},
		"pageNum": pageNum,
	})
}

func Post(ctx *neko.Context) {
	slug := ctx.Params.ByGet("slug")
	appSetting := models.GetAppSetting()
	opts := &models.Options{
		Include: "author",
		GormAdp: &models.GormAdapter{},
	}
	if appSetting.Permalink == "/:slug/" {
		opts.GormAdp.Query = "slug=?"
		opts.GormAdp.Args = []interface{}{slug}
	}

	post, errApi := models.GetPost(opts)

	if post == nil || errApi.Code != 0 || post.Id == 0 || post.Status == models.PostStatusDraft || post.Page {
		ctx.Render(themeViews("404"), neko.JSON{
			"tpl": "404",
			"app": appSetting,
		})
		return
	}

	ctx.Render(themeViews("post"), neko.JSON{
		"tpl":  "post",
		"app":  appSetting,
		"meta": post.GetPageMeta(),
		"post": post,
	})
}

func Page(ctx *neko.Context) {
	slug := ctx.Params.ByGet("slug")
	appSetting := models.GetAppSetting()
	opts := &models.Options{
		GormAdp: &models.GormAdapter{},
	}
	if appSetting.Permalink == "/:slug/" {
		opts.GormAdp.Query = "slug=?"
		opts.GormAdp.Args = []interface{}{slug}
	}

	post, errApi := models.GetPost(opts)

	if post == nil || errApi.Code != 0 || post.Id == 0 || post.Status == models.PostStatusDraft || !post.Page {
		ctx.Render(themeViews("404"), neko.JSON{
			"tpl": "404",
			"app": appSetting,
		})
		return
	}

	ctx.Render(themeViews("page"), neko.JSON{
		"tpl":  "page",
		"app":  appSetting,
		"meta": post.GetPageMeta(),
		"post": post,
	})
}

func Login(ctx *neko.Context) {
	reUrl := ctx.Params.ByGet("reurl")
	if middleware.Context.User != nil {
		if reUrl == "" {
			ctx.Redirect("/goant/#/")
		} else {
			ctx.Redirect(reUrl)
		}
		return
	}
	ctx.Render("#backend/login", neko.JSON{})
}

func Logout(ctx *neko.Context) {
	ctx.Session.Delete(middleware.SESSION_USER_ID)
	ctx.Redirect(setting.Host.Path)
}

func Admin(ctx *neko.Context) {
	ctx.Render("#backend/index", neko.JSON{"PageTitle": middleware.Context.User.Name})
}

func Installer(ctx *neko.Context) {
	if middleware.IsInstalled && ctx.Req.URL.Path == middleware.INSTALLER_URL {
		ctx.Redirect(setting.Host.Path)
		return
	}
	ctx.Render("#backend/setup", nil)
}

// AssetsHandler site theme assets handler
func AssetsHandler(ctx *neko.Context) {
	appSetting := models.GetAppSetting()
	http.ServeFile(ctx.Writer, ctx.Req, filepath.Join("content/themes", appSetting.ActiveTheme, "assets", ctx.Params.ByGet("filepath")))
}
