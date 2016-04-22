package services

import (
	"github.com/rocwong/neko"

	"github.com/go-ant/ant/core/server/models"
	"github.com/go-ant/ant/core/server/modules/middleware"
	"github.com/go-ant/ant/core/server/modules/utils"
	"github.com/go-ant/ant/core/server/modules/utils/acceptlang"
)

func UserCreate(ctx *neko.Context) {
	loginUser := middleware.Context.User
	if !loginUser.IsGranted("add-users", nil) {
		ctx.Json(models.RestApi{Error: models.ApiMsg.NoPermission})
		return
	}

	dataJson := ctx.Params.Json()
	user := models.User{
		Name:      dataJson.GetString("name"),
		Password:  dataJson.GetString("password"),
		Language:  acceptlang.ReadHeader(ctx.Req.Header).Best().Language,
		CreatedBy: loginUser.Id,
	}

	opts := &models.Options{Role: &models.Role{Id: dataJson.GetUInt32("role_id")}}
	errApi := models.CreateUser(&user, opts)

	ctx.Json(models.RestApi{Data: user, Error: errApi})
}

func UserEdit(ctx *neko.Context) {
	loginUser := middleware.Context.User
	if utils.ToUint32(ctx.Params.ByGet("user_id")) != loginUser.Id && !loginUser.IsGranted("edit-users", nil) {
		ctx.Json(models.RestApi{Error: models.ApiMsg.NoPermission})
		return
	}

	targetUser, _ := models.GetUserById(utils.ToUint32(ctx.Params.ByGet("user_id")), nil)
	if targetUser.Id == 0 {
		ctx.Json(models.RestApi{Error: models.ApiMsg.ErrUserNotFound})
		return
	}

	dataJson := ctx.Params.Json()
	targetUser.Name = dataJson.GetString("name")
	targetUser.Email = dataJson.GetString("email")
	targetUser.Location = dataJson.GetString("location")
	targetUser.Website = dataJson.GetString("website")
	targetUser.Bio = dataJson.GetString("bio")
	targetUser.Avatar = dataJson.GetString("avatar")
	targetUser.Cover = dataJson.GetString("cover")
	targetUser.UpdatedBy = loginUser.Id

	opts := &models.Options{User: loginUser, Role: &models.Role{Id: dataJson.GetUInt32("role_id")}}
	errApi := models.EditUser(targetUser, opts)

	ctx.Json(models.RestApi{Data: targetUser, Error: errApi})
}

func UserInfo(ctx *neko.Context) {
	var errApi models.ApiErr
	var user *models.User
	account := ctx.Params.ByGet("account")

	opts := &models.Options{Include: ctx.Params.ByGet("include")}

	if !utils.IsInt(account) {
		user, errApi = models.GetUserByName(account, opts)
	} else {
		if account == "0" {
			user = middleware.Context.User
			errApi = models.ApiMsg.Success
			if user == nil {
				errApi = models.ApiMsg.ErrUserNotFound
			} else if opts.Include == "role" {
				user.GetRole()
			}
		} else {
			user, errApi = models.GetUserById(utils.ToUint32(account), opts)
		}
	}

	if errApi.Code != 0 {
		ctx.Json(models.RestApi{Error: errApi})
		return
	}
	ctx.Json(models.RestApi{Data: user, Error: errApi})
}

func UserList(ctx *neko.Context) {
	loginUser := middleware.Context.User
	if !loginUser.IsGranted("browse-users", nil) {
		ctx.Json(models.RestApi{Error: models.ApiMsg.NoPermission})
		return
	}

	opts := &models.Options{
		Limit:   utils.ToUint32(ctx.Params.ByGet("limit")),
		Page:    utils.ToUint32(ctx.Params.ByGet("page")),
		Include: ctx.Params.ByGet("include"),
		GormAdp: &models.GormAdapter{},
	}

	list, recordCount, errApi := models.GetUsers(opts)
	pagination := &models.Pagination{
		PerPage:    opts.Limit,
		Page:       opts.Page,
		TotalPages: recordCount,
	}
	ctx.Json(models.RestApi{Data: list, Error: errApi, Pagination: pagination})
}

func UserChangePassword(ctx *neko.Context) {
	loginUser := middleware.Context.User
	dataJson := ctx.Params.Json()

	userId := utils.ToUint32(ctx.Params.ByGet("user_id"))
	oldPassword := dataJson.GetString("old_password")
	newPassword := dataJson.GetString("new_password")
	verifyPassword := dataJson.GetString("verify_password")

	_, errApi := models.ChangePassword(oldPassword, newPassword, verifyPassword, userId, &models.Options{User: loginUser})
	ctx.Json(models.RestApi{Error: errApi})
}

func UserLogin(ctx *neko.Context) {
	dataJson := ctx.Params.Json()
	user := &models.User{
		Name:     dataJson.GetString("name"),
		Password: dataJson.GetString("password"),
	}
	user, errApi := models.UserSignin(user, nil)
	if user == nil {
		ctx.Json(models.RestApi{Error: errApi})
		return
	}

	ctx.Session.Set(middleware.SESSION_USER_ID, user.Id)

	ctx.Json(models.RestApi{Data: user, Error: errApi})
}
