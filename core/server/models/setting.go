package models

import (
	"encoding/json"
	"github.com/go-ant/ant/core/server/modules/cache"
	"github.com/go-ant/ant/core/server/modules/utils"
	"strings"
	"time"
)

type JsonSetting struct {
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Logo         string        `json:"logo"`
	Cover        string        `json:"cover"`
	Language     string        `json:"language"`
	PostsPerPage uint32        `json:"posts_per_page"`
	ForceI18n    bool          `json:"force_i18n"`
	ActiveApps   []string      `json:"active_apps"`
	InstallApps  []string      `json:"install_apps"`
	AntHead      string        `json:"ant_head"`
	AntFoot      string        `json:"ant_foot"`
	Navigation   []*Navigation `json:"navigation"`
	Permalink    string        `json:"permalink"`
	ActiveTheme  string        `json:"active_theme"`
	IsPrivate    bool          `json:"is_private"`
	Password     string        `json:"password"`
	Version      string        `json:"version"`
}

func (c *JsonSetting) NavigationSerialize() string {
	if len(c.Navigation) > 0 {
		var nav = make([]string, 0)
		for _, j := range c.Navigation {
			b, err := json.Marshal(j)
			if err == nil {
				nav = append(nav, string(b))
			}
		}
		return strings.Join(nav, " ")
	}
	return ""
}

type Setting struct {
	Id        uint32    `gorm:"primary_key"`
	Key       string    `sql:"not null;type:varchar(50);"`
	Value     string    `sql:"type:text"`
	Type      string    `sql:"type:varchar(20);"`
	CreatedAt time.Time `sql:"not null;type:datetime;"`
	CreatedBy uint32    `sql:"not null;type:bigint unsigned;"`
	UpdatedAt time.Time `sql:"type:datetime;"`
	UpdatedBy uint32    `sql:"type:bigint unsigned;"`
}

type Navigation struct {
	Label string `json:"label"`
	Url   string `json:"url"`
}

const CacheKeyAppSettings string = "cache_app_settings"

func CreateSettings(settings []*Setting) ApiErr {
	var recordCount uint32
	db.Model(&Setting{}).Count(&recordCount)
	if recordCount == 0 {
		tx := db.Begin()
		for _, setting := range settings {
			tx.Create(setting)
		}
		tx.Commit()
	}
	return ApiMsg.Success
}

func EditSetting(setting *JsonSetting, opts *Options) ApiErr {
	tx := db.Begin()

	tx.Model(&Setting{}).Where("`key` = ?", "title").UpdateColumns(map[string]interface{}{"value": setting.Title, "updated_by": opts.User.Id})
	tx.Model(&Setting{}).Where("`key` = ?", "description").UpdateColumns(map[string]interface{}{"value": setting.Description, "updated_by": opts.User.Id})
	tx.Model(&Setting{}).Where("`key` = ?", "logo").UpdateColumns(map[string]interface{}{"value": setting.Logo, "updated_by": opts.User.Id})
	tx.Model(&Setting{}).Where("`key` = ?", "cover").UpdateColumns(map[string]interface{}{"value": setting.Cover, "updated_by": opts.User.Id})
	tx.Model(&Setting{}).Where("`key` = ?", "language").UpdateColumns(map[string]interface{}{"value": setting.Language, "updated_by": opts.User.Id})
	tx.Model(&Setting{}).Where("`key` = ?", "posts_per_page").UpdateColumns(map[string]interface{}{"value": setting.PostsPerPage, "updated_by": opts.User.Id})
	tx.Model(&Setting{}).Where("`key` = ?", "force_i18n").UpdateColumns(map[string]interface{}{"value": setting.ForceI18n, "updated_by": opts.User.Id})
	tx.Model(&Setting{}).Where("`key` = ?", "ant_head").UpdateColumns(map[string]interface{}{"value": setting.AntHead, "updated_by": opts.User.Id})
	tx.Model(&Setting{}).Where("`key` = ?", "ant_foot").UpdateColumns(map[string]interface{}{"value": setting.AntFoot, "updated_by": opts.User.Id})
	tx.Model(&Setting{}).Where("`key` = ?", "navigation").UpdateColumns(map[string]interface{}{"value": setting.NavigationSerialize(), "updated_by": opts.User.Id})
	tx.Model(&Setting{}).Where("`key` = ?", "permalink").UpdateColumns(map[string]interface{}{"value": setting.Permalink, "updated_by": opts.User.Id})
	tx.Model(&Setting{}).Where("`key` = ?", "active_theme").UpdateColumns(map[string]interface{}{"value": setting.ActiveTheme, "updated_by": opts.User.Id})
	tx.Model(&Setting{}).Where("`key` = ?", "is_private").UpdateColumns(map[string]interface{}{"value": setting.IsPrivate, "updated_by": opts.User.Id})
	tx.Model(&Setting{}).Where("`key` = ?", "password").UpdateColumns(map[string]interface{}{"value": setting.Password, "updated_by": opts.User.Id})

	tx.Commit()

	return ApiMsg.Saved
}

func GetSetting(opts *Options) (*JsonSetting, ApiErr) {
	dbInner := initDb(opts)
	setting := make([]*Setting, 0)
	if err := dbInner.Find(&setting).Error; err != nil {
		return nil, UnknowError(err.Error())
	}

	jsonSetting := new(JsonSetting)
	for _, n := range setting {
		switch strings.ToLower(n.Key) {
		case "version":
			jsonSetting.Version = n.Value
		case "title":
			jsonSetting.Title = n.Value
		case "description":
			jsonSetting.Description = n.Value
		case "logo":
			jsonSetting.Logo = n.Value
		case "cover":
			jsonSetting.Cover = n.Value
		case "language":
			jsonSetting.Language = n.Value
		case "posts_per_page":
			jsonSetting.PostsPerPage = utils.ToUint32(n.Value)
		case "force_i18n":
			jsonSetting.ForceI18n = utils.ToBool(n.Value)
		case "active_apps":
			jsonSetting.ActiveApps = strings.Fields(n.Value)
		case "install_apps":
			jsonSetting.InstallApps = strings.Fields(n.Value)
		case "ant_head":
			jsonSetting.AntHead = n.Value
		case "ant_foot":
			jsonSetting.AntFoot = n.Value
		case "navigation":
			jsonSetting.Navigation = make([]*Navigation, 0)
			for _, jsonString := range strings.Fields(n.Value) {
				nav := &Navigation{}
				json.Unmarshal([]byte(jsonString), nav)
				jsonSetting.Navigation = append(jsonSetting.Navigation, nav)
			}
		case "permalink":
			jsonSetting.Permalink = n.Value
		case "active_theme":
			jsonSetting.ActiveTheme = n.Value
		case "is_private":
			jsonSetting.IsPrivate = utils.ToBool(n.Value)
		case "password":
			jsonSetting.Password = n.Value
		}
	}

	return jsonSetting, ApiMsg.Success
}

// GetAppSetting get setting from cache
func GetAppSetting() *JsonSetting {
	appSetting, found := cache.Get(CacheKeyAppSettings)
	if !found {
		opts := &Options{}
		appSetting, _ := GetSetting(opts)
		if appSetting != nil {
			cache.Set(CacheKeyAppSettings, appSetting)
		}
		return appSetting
	}
	return appSetting.(*JsonSetting)
}
