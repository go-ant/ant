package services

import (
	"io/ioutil"

	"github.com/rocwong/neko"

	"github.com/go-ant/ant/core/server/models"
)

func ThemeList(ctx *neko.Context) {
	rd, err := ioutil.ReadDir("content/themes")
	if err != nil {
		ctx.Json(models.RestApi{Error: models.ApiErr{Code: 10, Message: err.Error()}})
		return
	}

	themes := make([]string, 0)
	for _, info := range rd {
		if info.IsDir() {
			themes = append(themes, info.Name())
		}
	}

	ctx.Json(models.RestApi{Data: themes, Error: models.ApiMsg.Success})
}
