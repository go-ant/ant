package services

import (
	"github.com/rocwong/neko"

	"github.com/go-ant/ant/core/server/models"
)

func PermissionList(ctx *neko.Context) {
	list, err := models.GetPermissions(nil)
	ctx.Json(models.RestApi{Data: list, Error: err})
}
