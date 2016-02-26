package middleware

import (
	"github.com/go-ant/ant/core/server/models"
	"github.com/go-ant/ant/core/server/modules/capabilities"
	"github.com/go-ant/ant/core/server/modules/setting"
	"github.com/rocwong/neko"
	"path"
	"strings"
)

var (
	INSTALLER_URL = path.Join(setting.Host.Path, "/goant/setup")
	IsInstalled   = false
)

// Installer initializes site data
func Installer() neko.HandlerFunc {
	return func(ctx *neko.Context) {
		if !IsInstalled {
			if users, _, _ := models.GetUsers(nil); len(users) == 0 {
				if "get" == strings.ToLower(ctx.Req.Method) && ctx.Req.URL.Path != INSTALLER_URL {
					ctx.Redirect(INSTALLER_URL)
					ctx.Abort()
				}
			} else {
				IsInstalled = true

				// cache all roles permission
				roles, _ := models.GetRoles(nil)
				for _, role := range roles {
					role.GetPermissions()
					permissionsToAdd := make([]string, 0, len(role.Permissions))
					for _, perm := range role.Permissions {
						permissionsToAdd = append(permissionsToAdd, perm.Slug)
					}
					capabilities.SetRole(role.Slug, permissionsToAdd)
				}
			}
		}
	}
}
