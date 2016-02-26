package capabilities

import (
	"github.com/go-ant/ant/core/server/modules/capabilities/rbac"
)

var allRole *rbac.Rbac

func init() {
	allRole = rbac.New()
}

// SetRole add/edit a role
func SetRole(name string, perms []string) {
	if allRole.Get(name) == nil {
		allRole.Add(name, perms, nil)
	} else {
		allRole.Set(name, perms, nil)
	}
}

// IsGranted check if given name has permission
func IsGranted(name, permission string, assert rbac.AssertionFunc) bool {
	return allRole.IsGranted(name, permission, assert)
}
