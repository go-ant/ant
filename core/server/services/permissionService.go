package services

import (
	"github.com/go-ant/ant/core/server/models"
	"github.com/rocwong/neko"
)

func PermissionList(ctx *neko.Context) {
	list, err := models.GetPermissions(nil)
	ctx.Json(models.RestApi{Data: list, Error: err})
}
