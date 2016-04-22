package services

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/rocwong/neko"

	"github.com/go-ant/ant/core/server/models"
	"github.com/go-ant/ant/core/server/modules/middleware"
	"github.com/go-ant/ant/core/server/modules/setting"
	"github.com/go-ant/ant/core/server/modules/utils"
	"github.com/go-ant/ant/core/server/modules/utils/uuid"
)

func Upload(ctx *neko.Context) {
	loginUser := middleware.Context.User

	if loginUser == nil {
		ctx.Json(models.RestApi{Error: models.ApiMsg.NoPermission})
		return
	}

	err := ctx.Req.ParseMultipartForm(int64(setting.API.UploadMaxSize))
	if err != nil {
		ctx.Json(models.RestApi{Error: models.ApiErr{Code: 10, Message: err.Error()}})
		return
	}

	file, fileInfo, err := ctx.Req.FormFile("file")
	if err != nil {
		ctx.Json(models.RestApi{Error: models.ApiErr{Code: 10, Message: err.Error()}})
		return
	}
	defer file.Close()

	// Check file type
	fileExt := filepath.Ext(fileInfo.Filename)
	if !utils.StringInSlice(fileExt, setting.API.UploadExtensions) ||
		!utils.StringInSlice(fileInfo.Header.Get("content-type"), setting.API.UploadContentTypes) {
		ctx.Json(models.RestApi{Error: models.ApiMsg.ErrFileNotSupported})
		return
	}

	fileSize, err := file.Seek(0, 2)
	if err != nil {
		ctx.Json(models.RestApi{Error: models.ApiErr{Code: 10, Message: err.Error()}})
		return
	}

	// Check file size
	if int64(setting.API.UploadMaxSize) < fileSize {
		ctx.Json(models.RestApi{Error: models.ApiMsg.ErrFileTooLarge})
		return
	}

	savePath := time.Now().Format("/2006/01")
	fileName := uuid.NewV4().String() + fileExt
	fullPath := path.Join(setting.API.UploadFolder, savePath)

	err = os.MkdirAll(fullPath, os.ModePerm)
	if err != nil {
		ctx.Json(models.RestApi{Error: models.ApiErr{Code: 10, Message: err.Error()}})
		return
	}

	dst, err := os.Create(path.Join(fullPath, fileName))
	if err != nil {
		ctx.Json(models.RestApi{Error: models.ApiErr{Code: 10, Message: err.Error()}})
		return
	}
	defer dst.Close()

	file.Seek(0, 0)
	if _, err = io.Copy(dst, file); err != nil {
		ctx.Json(models.RestApi{Error: models.ApiErr{Code: 10, Message: err.Error()}})
		return
	}

	fileFullPath := path.Join(setting.Host.Path, setting.API.FilesPath, savePath, fileName)
	ctx.Json(models.RestApi{Data: neko.JSON{"file": fileFullPath}, Error: models.ApiMsg.Success})
}
