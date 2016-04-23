package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"time"

	"github.com/ahmetalpbalkan/go-linq"

	"github.com/go-ant/ant/core/server/models"
	"github.com/go-ant/ant/core/server/modules/capabilities"
	"github.com/go-ant/ant/core/server/modules/setting"
	"github.com/go-ant/ant/core/server/modules/utils"
)

const (
	permissions_json_path = "content/data/fixtures/permissions.json"
	roles_json_path       = "content/data/fixtures/roles.json"
	posts_json_path       = "content/data/fixtures/posts.json"
)

// GlobalInit initialize global configuration
func GlobalInit() {
	setting.NewContext()

	if setting.InstallLock {
		if !models.HasEngine {
			if err := models.NewEngine(); err != nil {

			}
		}
		// cache all roles permission
		capabilities.RemoveAll()
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

// DoImport import base data
func DoImport() {
	dataPermissons := GetPermissions()
	dataRoles := GetRoles()

	appUrl := "/"
	if !utils.IsEmpty(setting.Host.Path) {
		appUrl = setting.Host.Path
	}
	strNav, _ := json.Marshal(models.Navigation{Label: "Home", Url: appUrl})


	// import app settings
	appSettings := make([]*models.Setting, 0)
	appSettings = append(appSettings, &models.Setting{Key: "app_url", Value: appUrl, Type: "blog", CreatedBy: 1})
	appSettings = append(appSettings, &models.Setting{Key: "title", Value: "GoAnt", Type: "blog", CreatedBy: 1})
	appSettings = append(appSettings, &models.Setting{Key: "description", Value: "A blog engine written in Go", Type: "blog", CreatedBy: 1})
	appSettings = append(appSettings, &models.Setting{Key: "logo", Value: "", Type: "blog", CreatedBy: 1})
	appSettings = append(appSettings, &models.Setting{Key: "cover", Value: setting.Host.Path + "/goant/assets/css/images/cover.jpg", Type: "blog", CreatedBy: 1})
	appSettings = append(appSettings, &models.Setting{Key: "language", Value: "en", Type: "blog", CreatedBy: 1})
	appSettings = append(appSettings, &models.Setting{Key: "posts_per_page", Value: "15", Type: "blog", CreatedBy: 1})
	appSettings = append(appSettings, &models.Setting{Key: "force_i18n", Value: "1", Type: "blog", CreatedBy: 1})
	appSettings = append(appSettings, &models.Setting{Key: "active_apps", Value: "", Type: "blog", CreatedBy: 1})
	appSettings = append(appSettings, &models.Setting{Key: "install_apps", Value: "", Type: "blog", CreatedBy: 1})
	appSettings = append(appSettings, &models.Setting{Key: "is_private", Value: "0", Type: "private", CreatedBy: 1})
	appSettings = append(appSettings, &models.Setting{Key: "password", Value: "", Type: "private", CreatedBy: 1})
	appSettings = append(appSettings, &models.Setting{Key: "active_theme", Value: "journey", Type: "theme", CreatedBy: 1})
	appSettings = append(appSettings, &models.Setting{Key: "ant_head", Value: "", Type: "blog", CreatedBy: 1})
	appSettings = append(appSettings, &models.Setting{Key: "ant_foot", Value: "", Type: "blog", CreatedBy: 1})
	appSettings = append(appSettings, &models.Setting{Key: "navigation", Value: string(strNav), Type: "blog", CreatedBy: 1})
	appSettings = append(appSettings, &models.Setting{Key: "permalink", Value: "/:slug/", Type: "blog", CreatedBy: 1})
	appSettings = append(appSettings, &models.Setting{Key: "version", Value: models.APP_VERSION, Type: "blog", CreatedBy: 1})
	models.CreateSettings(appSettings)

	// import permisson data
	permissons := make([]*models.Permission, 0, len(dataPermissons))
	for key, obj := range dataPermissons {
		for _, perm := range obj {
			perm := &models.Permission{
				Name:       perm["name"],
				Slug:       perm["slug"],
				ObjectType: key,
				IsCore:     true,
			}
			permissons = append(permissons, perm)
			models.CreatePermission(perm, nil)
		}
	}

	// import role data
	var roleKeys []string
	for k := range dataRoles {
		roleKeys = append(roleKeys, k)
	}
	sort.Strings(roleKeys)

	for _, k := range roleKeys {
		rolePerm := dataRoles[k]

		permsFilter := make([]linq.T, 0)
		for key, perm := range rolePerm.Permissions {
			if len(perm) == 1 && perm[0] == "all" {
				filters, _ := linq.From(permissons).Where(func(s linq.T) (bool, error) {
					return s.(*models.Permission).ObjectType == key, nil
				}).Results()
				permsFilter = append(permsFilter, filters...)
			} else {
				filters, _ := linq.From(permissons).Where(func(s linq.T) (bool, error) {
					return utils.StringInSlice(s.(*models.Permission).Slug, perm), nil
				}).Results()
				permsFilter = append(permsFilter, filters...)
			}
		}

		perms := make([]*models.Permission, 0)
		for _, perm := range permsFilter {
			perms = append(perms, perm.(*models.Permission))
		}
		role := &models.Role{
			Name:        rolePerm.Name,
			Description: rolePerm.Description,
			Slug:        strings.ToLower(rolePerm.Slug),
		}

		models.CreateRole(role, &models.Options{Permissions: perms})
	}
}

// ImportPosts import post data
func ImportPosts(user *models.User) {
	dataPosts := GetPosts()
	var postKeys []string
	for k := range dataPosts {
		postKeys = append(postKeys, k)
	}
	sort.Strings(postKeys)

	for _, k := range postKeys {
		post := &models.Post{
			Title:       dataPosts[k].Title,
			Slug:        dataPosts[k].Slug,
			Markdown:    dataPosts[k].Markdown,
			AuthorId:    user.Id,
			Language:    user.Language,
			Status:      models.PostStatusPublished,
			PublishedAt: time.Now(),
			PublishedBy: user.Id,
			CreatedBy:   user.Id,
		}
		models.CreatePost(post, nil)
	}
}

// GetPermissions load permission data from json file
func GetPermissions() map[string][]map[string]string {
	jsonPermissions := map[string][]map[string]string{}

	bytes, err := ioutil.ReadFile(permissions_json_path)
	if err != nil {
		panic(fmt.Errorf("Err: Failed to read file 'permission.json'"))
	}
	if err = json.Unmarshal(bytes, &jsonPermissions); err != nil {
		panic(fmt.Errorf("Err: %v", err))
	}
	return jsonPermissions
}

// GetRoles load role data from json file
func GetRoles() map[string]rolePermissions {
	jsonRoles := map[string]rolePermissions{}

	bytes, err := ioutil.ReadFile(roles_json_path)
	if err != nil {
		panic(fmt.Errorf("Err: Failed to read file 'roles.json'"))
	}
	if err = json.Unmarshal(bytes, &jsonRoles); err != nil {
		panic(fmt.Errorf("Err: %v", err))
	}

	return jsonRoles
}

func GetPosts() map[string]post {
	jsonPosts := map[string]post{}

	bytes, err := ioutil.ReadFile(posts_json_path)
	if err != nil {
		panic(fmt.Errorf("Err: Failed to read file 'posts.json'"))
	}
	if err = json.Unmarshal(bytes, &jsonPosts); err != nil {
		panic(fmt.Errorf("Err: %v", err))
	}

	return jsonPosts
}
