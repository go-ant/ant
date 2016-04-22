package middleware

import (
	"github.com/rocwong/neko"

	"github.com/go-ant/ant/core/server/modules/setting"
)

func UnknowPage() neko.HandlerFunc {
	return func(ctx *neko.Context) {
		ctx.Next()
		if !ctx.Writer.Written() {
			switch ctx.Writer.Status() {
			case 404:
				ctx.Render("#backend/404", neko.JSON{"Home": setting.Host.Path, "Code": 404}, 404)
			}
		}
	}
}
