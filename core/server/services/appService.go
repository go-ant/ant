package services

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"time"

	"github.com/rocwong/neko"
	"gopkg.in/ini.v1"

	"github.com/go-ant/ant/core/server/data"
	"github.com/go-ant/ant/core/server/models"
	"github.com/go-ant/ant/core/server/modules/middleware"
	"github.com/go-ant/ant/core/server/modules/setting"
	"github.com/go-ant/ant/core/server/modules/utils"
	"github.com/go-ant/ant/core/server/modules/utils/uuid"
)

// base64ImgUpload gets image from base64, return image path
func base64ImgUpload(data string) (string, error) {

	regBase64, _ := regexp.Compile("^data:[\\w]+/[\\w]+;base64,")
	strBase64 := regBase64.ReplaceAllString(data, "")
	byteBase64, _ := base64.StdEncoding.DecodeString(strBase64)

	img, err := png.Decode(bytes.NewReader(byteBase64))
	if err != nil {
		return "", err
	}

	savePath := time.Now().Format("/2006/01/")
	fileName := uuid.NewV4().String() + ".jpg"
	fullPath := setting.API.UploadFolder + savePath

	err = os.MkdirAll(fullPath, os.ModePerm)
	if err != nil {
		return "", err
	}
	out, err := os.Create(fullPath + fileName)

	if err != nil {
		return "", err
	}
	defer out.Close()

	if err = png.Encode(out, img); err != nil {
		return "", err
	}

	return path.Join(setting.Host.Path, setting.API.FilesPath, savePath, fileName), nil
}

func AppInstall(ctx *neko.Context) {
	dataJson := ctx.Params.Json()
	appTitle := dataJson.GetString("title")
	user := models.User{
		Name:     dataJson.GetString("name"),
		Password: dataJson.GetString("password"),
		Avatar:   dataJson.GetString("avatar"),
		Language: dataJson.GetString("language"),
	}
	if utils.IsEmpty(user.Name) {
		ctx.Json(models.RestApi{Error: models.ApiMsg.ErrUserNameCanNotBeEmpty})
		return
	}

	if user.Avatar != "" {
		avatarPath, err := base64ImgUpload(user.Avatar)
		if err != nil {
			ctx.Json(models.RestApi{Error: models.UnknowError(err.Error())})
			return
		}
		user.Avatar = avatarPath
	}

	var errApi models.ApiErr
	if setting.InstallLock {
		errApi = models.ApiMsg.NoPermission
	} else {
		if !models.HasEngine {
			if err := models.NewEngine(); err != nil {
				ctx.Json(models.RestApi{Error: models.ApiMsg.DatabaseFailed})
				return
			}
		}

		// create tables
		if errApi := models.InitialDatabase(); !errApi.IsSuccess() {
			ctx.Json(models.RestApi{Error: errApi})
			return
		}

		// initial base data
		data.DoImport()

		// create owner
		opts := &models.Options{
			Role: &models.Role{Slug: models.SiteOwner},
		}
		ownerRole, _ := models.GetRole(opts)
		if ownerRole == nil {
			errApi = models.ApiMsg.SaveFail
		} else {
			opts.Role = ownerRole
			errApi = models.CreateUser(&user, opts)
		}

		// update app setting
		jsonSetting := models.GetAppSetting()
		jsonSetting.Title = appTitle
		jsonSetting.Language = user.Language
		if utils.IsEmpty(jsonSetting.Title) {
			jsonSetting.Title = user.Name + "'s blog"
		}
		opts = &models.Options{User: &user}
		models.EditSetting(jsonSetting, opts)

		// import post data
		data.ImportPosts(&user)

		// save custom settings
		if !utils.IsFile(setting.CustomConf) {
			os.MkdirAll(filepath.Dir(setting.CustomConf), os.ModePerm)
		}

		cfg, _ := ini.Load(setting.CustomConf)
		cfg.Section("security").Key("install_lock").SetValue("true")
		cfg.SaveTo(setting.CustomConf)

		data.GlobalInit()

	}

	ctx.Json(models.RestApi{Data: nil, Error: errApi})
}

func InitialApp(ctx *neko.Context) {
	loginUser := middleware.Context.User

	if loginUser == nil {
		ctx.Json(models.RestApi{Error: models.ApiMsg.NoPermission})
		return
	}
	menus := getMenus(loginUser)

	// get user permissions
	role := loginUser.GetRole()
	role.GetPermissions()
	permissions := make([]string, 0, len(role.Permissions))
	for _, perm := range role.Permissions {
		permissions = append(permissions, perm.Slug)
	}

	ctx.Json(models.RestApi{
		Data: neko.JSON{
			"menus":       menus,
			"user":        loginUser,
			"permissions": permissions,
			"language":    models.GetAppSetting().Language,
		},
		Error: models.ApiMsg.Success,
	})
}

func getMenus(loginUser *models.User) []models.Menu {

	menus := make([]models.Menu, 0)

	//	menus = append(menus, models.Menu{
	//		Name:  "nav.dashboard",
	//		Label: "dashboard",
	//		Icon:  "dashboard",
	//		URL:   "dashboard",
	//	})

	if loginUser.IsGranted("browse-posts", nil) || loginUser.IsGranted("edit-all-posts", nil) {
		subMenus := make([]models.Menu, 0)
		subMenus = append(subMenus, models.Menu{Name: "nav.post_list", Label: "post", URL: "posts"})

		if loginUser.IsGranted("add-posts", nil) {
			subMenus = append(subMenus, models.Menu{Name: "nav.post_add", Label: "add post", URL: "posts/add"})
		}

		//		subMenus = append(subMenus, models.Menu{Name: "nav.post_category", Label: "post category", URL: "posts/categories"})
		//		subMenus = append(subMenus, models.Menu{Name: "nav.post_tag", Label: "post tag", URL: "posts/tags"})

		menus = append(menus, models.Menu{
			Name:  "nav.post_manage",
			Label: "post namage",
			Icon:  "create",
			List:  subMenus,
		})
	}

	//	if loginUser.IsGranted("switch-themes", nil) {
	//		subMenus := make([]models.Menu, 0)
	//		subMenus = append(subMenus, models.Menu{Name: "nav.theme", Label: "theme", URL: "appearance/themes"})
	//
	//		menus = append(menus, models.Menu{
	//			Name:  "nav.appearance",
	//			Label: "appearance",
	//			Icon:  "web",
	//			List:  subMenus,
	//		})
	//
	//	}

	if loginUser.IsGranted("browse-users", nil) || loginUser.IsGranted("browse-roles", nil) {
		subMenus := make([]models.Menu, 0)

		if loginUser.IsGranted("browse-users", nil) {
			subMenus = append(subMenus, models.Menu{Name: "nav.user_list", Label: "user list", URL: "users"})
		}

		if loginUser.IsGranted("add-users", nil) {
			subMenus = append(subMenus, models.Menu{Name: "nav.user_add", Label: "add user", URL: "users/add"})
		}

		if loginUser.IsGranted("browse-roles", nil) {
			subMenus = append(subMenus, models.Menu{Name: "nav.role_list", Label: "role list", URL: "roles"})
		}

		menus = append(menus, models.Menu{
			Name:  "nav.user_manage",
			Label: "user manage",
			Icon:  "group",
			List:  subMenus,
		})
	}

	if loginUser.IsGranted("edit-settings", nil) {
		subMenus := make([]models.Menu, 0)
		subMenus = append(subMenus, models.Menu{Name: "nav.general", Label: "general", URL: "settings/general"})
		subMenus = append(subMenus, models.Menu{Name: "nav.navigation", Label: "navigation", URL: "settings/navigation"})

		menus = append(menus, models.Menu{
			Name:  "nav.setting",
			Label: "settings",
			Icon:  "settings",
			List:  subMenus,
		})
	}

	return menus
}
