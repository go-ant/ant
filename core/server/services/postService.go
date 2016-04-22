package services

import (
	"github.com/rocwong/neko"

	"github.com/go-ant/ant/core/server/models"
	"github.com/go-ant/ant/core/server/modules/middleware"
	"github.com/go-ant/ant/core/server/modules/utils"
)

func PostCreate(ctx *neko.Context) {
	loginUser := middleware.Context.User
	if !loginUser.IsGranted("add-posts", nil) && !loginUser.IsGranted("edit-all-posts", nil) {
		ctx.Json(models.RestApi{Error: models.ApiMsg.NoPermission})
		return
	}

	dataJson := ctx.Params.Json()
	post := &models.Post{
		Id:              utils.ToUint32(ctx.Params.ByGet("post_id")),
		Title:           dataJson.GetString("title"),
		Slug:            dataJson.GetString("slug"),
		Markdown:        dataJson.GetString("markdown"),
		Cover:           dataJson.GetString("cover"),
		Language:        dataJson.GetString("language"),
		Page:            utils.ToBool(dataJson.Get("page")),
		Featured:        utils.ToBool(dataJson.Get("featured")),
		Status:          dataJson.GetString("status"),
		MetaTitle:       dataJson.GetString("meta_title"),
		MetaDescription: dataJson.GetString("meta_description"),
		AuthorId:        utils.ToUint32(dataJson.GetString("author_id")),
		PublishedAt:     utils.ToTime(dataJson.Get("published_at")),
		PublishedBy:     loginUser.Id,
		CreatedBy:       loginUser.Id,
		UpdatedBy:       loginUser.Id,
	}
	post.AuthorId = loginUser.Id
	post.Language = loginUser.Language

	opts := &models.Options{}
	errApi := models.CreatePost(post, opts)

	ctx.Json(models.RestApi{Data: post, Error: errApi})
}

func PostEdit(ctx *neko.Context) {
	loginUser := middleware.Context.User
	if !loginUser.IsGranted("edit-posts", nil) && !loginUser.IsGranted("edit-all-posts", nil) {
		ctx.Json(models.RestApi{Error: models.ApiMsg.NoPermission})
		return
	}

	dataJson := ctx.Params.Json()
	post := &models.Post{
		Id:              utils.ToUint32(ctx.Params.ByGet("post_id")),
		Title:           dataJson.GetString("title"),
		Slug:            dataJson.GetString("slug"),
		Markdown:        dataJson.GetString("markdown"),
		Cover:           dataJson.GetString("cover"),
		Language:        dataJson.GetString("language"),
		Page:            utils.ToBool(dataJson.Get("page")),
		Featured:        utils.ToBool(dataJson.Get("featured")),
		Status:          dataJson.GetString("status"),
		MetaTitle:       dataJson.GetString("meta_title"),
		MetaDescription: dataJson.GetString("meta_description"),
		AuthorId:        utils.ToUint32(dataJson.GetString("author_id")),
		PublishedAt:     utils.ToTime(dataJson.Get("published_at")),
		PublishedBy:     loginUser.Id,
		UpdatedBy:       loginUser.Id,
	}
	post.AuthorId = loginUser.Id
	post.Language = loginUser.Language

	opts := &models.Options{}
	errApi := models.EditPost(post, opts)

	ctx.Json(models.RestApi{Data: post, Error: errApi})
}

func PostDelete(ctx *neko.Context) {
	loginUser := middleware.Context.User
	if !loginUser.IsGranted("delete-posts", nil) && !loginUser.IsGranted("edit-all-posts", nil) {
		ctx.Json(models.RestApi{Error: models.ApiMsg.NoPermission})
		return
	}

	opts := &models.Options{}
	errApi := models.DeletePost(utils.ToUint32(ctx.Params.ByGet("post_id")), opts)
	ctx.Json(models.RestApi{Error: errApi})
}

func PostInfo(ctx *neko.Context) {
	loginUser := middleware.Context.User
	id := ctx.Params.ByGet("post_id")

	opts := &models.Options{Include: ctx.Params.ByGet("include")}
	post, errApi := models.GetPostById(utils.ToUint32(id), opts)

	if post == nil {
		ctx.Json(models.RestApi{Error: models.ApiMsg.ErrPostNotFound})
		return
	}

	if post.CreatedBy != loginUser.Id && !loginUser.IsGranted("edit-all-posts", nil) {
		ctx.Json(models.RestApi{Error: models.ApiMsg.NoPermission})
		return
	}

	ctx.Json(models.RestApi{Data: post, Error: errApi})
}

func PostList(ctx *neko.Context) {
	loginUser := middleware.Context.User

	opts := &models.Options{
		Limit:   utils.ToUint32(ctx.Params.ByGet("limit")),
		Page:    utils.ToUint32(ctx.Params.ByGet("page")),
		GormAdp: &models.GormAdapter{},
	}

	if loginUser.IsGranted("edit-all-posts", nil) {

	} else if loginUser.IsGranted("browse-posts", nil) {
		opts.GormAdp.Map = map[string]interface{}{"created_by": loginUser.Id}
	} else {
		ctx.Json(models.RestApi{Error: models.ApiMsg.NoPermission})
		return
	}

	opts.GormAdp.OrderBy = "published_at desc"
	opts.GormAdp.Columns = []string{"id", "title", "slug", "cover", "language", "page", "featured", "status", "published_at"}

	list, recordCount, errApi := models.GetPosts(opts)
	pagination := &models.Pagination{
		PerPage:    opts.Limit,
		Page:       opts.Page,
		TotalPages: recordCount,
	}
	ctx.Json(models.RestApi{Data: list, Error: errApi, Pagination: pagination})
}
