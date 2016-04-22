package services

import (
	"github.com/rocwong/neko"

	"github.com/go-ant/ant/core/server/models"
	"github.com/go-ant/ant/core/server/modules/cache"
	"github.com/go-ant/ant/core/server/modules/middleware"
)

func SettingEdit(ctx *neko.Context) {
	loginUser := middleware.Context.User
	if !loginUser.IsGranted("edit-settings", nil) {
		ctx.Json(models.RestApi{Error: models.ApiMsg.NoPermission})
		return
	}

	errApi := models.ApiMsg.SaveFail
	var jsonSetting *models.JsonSetting
	if ctx.Params.BindJSON(&jsonSetting) == nil {
		opts := &models.Options{User: loginUser}
		errApi = models.EditSetting(jsonSetting, opts)

		if errApi.IsSuccess() {
			cache.Set(models.CacheKeyAppSettings, jsonSetting)
		}

	}
	ctx.Json(models.RestApi{Error: errApi})
}

func SettingList(ctx *neko.Context) {
	loginUser := middleware.Context.User
	if !loginUser.IsGranted("edit-settings", nil) {
		ctx.Json(models.RestApi{Error: models.ApiMsg.NoPermission})
		return
	}
	appSetting := models.GetAppSetting()
	ctx.Json(models.RestApi{Data: appSetting, Error: models.ApiMsg.Success})
}
