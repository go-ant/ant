package services

import (
	"github.com/go-ant/ant/core/server/models"
	"github.com/go-ant/ant/core/server/modules/middleware"
	"github.com/go-ant/ant/core/server/modules/utils"
	"github.com/rocwong/neko"
	"strings"
)

func RoleCreate(ctx *neko.Context) {
	loginUser := middleware.Context.User
	if !loginUser.IsGranted("add-roles", nil) {
		ctx.Json(models.RestApi{Error: models.ApiMsg.NoPermission})
		return
	}

	dataJson := ctx.Params.Json()
	role := &models.Role{
		Id:          utils.ToUint32(ctx.Params.ByGet("role_id")),
		Name:        dataJson.GetString("name"),
		Description: dataJson.GetString("description"),
		CreatedBy:   loginUser.Id,
		UpdatedBy:   loginUser.Id,
	}

	opts := &models.Options{}
	opts.Permissions, _ = models.GetPermissionsByIds(strings.Split(dataJson.GetString("permissions"), ","))
	errApi := models.CreateRole(role, opts)

	ctx.Json(models.RestApi{Data: role, Error: errApi})
}

func RoleEdit(ctx *neko.Context) {
	loginUser := middleware.Context.User

	dataJson := ctx.Params.Json()
	role := &models.Role{
		Id:          utils.ToUint32(ctx.Params.ByGet("role_id")),
		Name:        dataJson.GetString("name"),
		Description: dataJson.GetString("description"),
		UpdatedBy:   loginUser.Id,
	}

	// can not edit owner-role
	if !loginUser.IsGranted("edit-roles", nil) || strings.ToLower(role.Name) == models.SiteOwner {
		ctx.Json(models.RestApi{Error: models.ApiMsg.NoPermission})
		return
	}

	opts := &models.Options{}
	opts.Permissions, _ = models.GetPermissionsByIds(strings.Split(dataJson.GetString("permissions"), ","))
	errApi := models.EditRole(role, opts)

	ctx.Json(models.RestApi{Error: errApi})
}

func RoleDelete(ctx *neko.Context) {
	loginUser := middleware.Context.User
	if !loginUser.IsGranted("delete-roles", nil) {
		ctx.Json(models.RestApi{Error: models.ApiMsg.NoPermission})
		return
	}

	opts := &models.Options{}
	errApi := models.DeleteRole(utils.ToUint32(ctx.Params.ByGet("role_id")), opts)

	ctx.Json(models.RestApi{Error: errApi})
}

func RoleInfo(ctx *neko.Context) {
	id := ctx.Params.ByGet("role_id")
	opts := &models.Options{Include: ctx.Params.ByGet("include")}

	role, errApi := models.GetRoleById(utils.ToUint32(id), opts)

	if role == nil {
		ctx.Json(models.RestApi{Error: models.ApiMsg.ErrRoleNotFound})
		return
	}
	// can not get owner-role data
	if role.Slug == models.SiteOwner {
		ctx.Json(models.RestApi{Error: models.ApiMsg.NoPermission})
		return
	}

	ctx.Json(models.RestApi{Data: role, Error: errApi})
}

func RoleList(ctx *neko.Context) {
	opts := &models.Options{
		GormAdp: &models.GormAdapter{},
	}
	opts.GormAdp.Query = "slug != 'owner'"
	opts.GormAdp.OrderBy = "id asc"
	list, errApi := models.GetRoles(opts)

	ctx.Json(models.RestApi{Data: list, Error: errApi})
}
