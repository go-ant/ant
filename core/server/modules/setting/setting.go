package setting

import (
	"fmt"
	"path"
	"strings"

	"gopkg.in/ini.v1"

	"github.com/go-ant/ant/core/server/modules/utils"
)

const (
	CONFIG_FILE = "config/app.ini"
	CUSTOM_PATH = "custom/"
)

var (
	EnableGzip  bool
	InstallLock bool

	Session struct {
		Name  string
		Key   string
		Store string
	}

	Host struct {
		Path string
		Addr string
		Port string
	}

	API struct {
		UploadFolder       string
		FilesPath          string
		UploadMaxSize      int
		UploadExtensions   []string
		UploadContentTypes []string
	}

	DbCfg struct {
		Type    string
		Host    string
		Name    string
		User    string
		Passwd  string
		SSLMode string
		Path    string
	}

	// Global setting objects
	INSTALLER_URL = path.Join(Host.Path, "/goant/install")
	CustomConf    string
)

func NewContext() {
	Cfg, err := ini.Load(CONFIG_FILE)
	if err != nil {
		panic(fmt.Errorf("Err: Failed to load `%s`: %v", CONFIG_FILE, err.Error()))
	}

	if len(CustomConf) == 0 {
		CustomConf = CUSTOM_PATH + CONFIG_FILE
	}
	if utils.IsFile(CustomConf) {
		Cfg.Append(CustomConf)
	}

	EnableGzip = Cfg.Section("").Key("enable_gzip").MustBool(false)
	Host.Path = Cfg.Section("").Key("host.path").MustString("/")
	Host.Addr = Cfg.Section("").Key("host.addr").MustString("/")
	Host.Port = Cfg.Section("").Key("host.port").MustString("2015")

	if !strings.HasPrefix(Host.Path, "/") {
		Host.Path = "/" + Host.Path
	}
	Host.Path = strings.TrimSuffix(Host.Path, "/")

	// API settings
	sec := Cfg.Section("api")
	API.UploadFolder = sec.Key("api.upload_folder").MustString("./content/upload")
	API.FilesPath = sec.Key("api.files_path").MustString("/upload")
	API.UploadMaxSize = sec.Key("api.upload_max_size").MustInt(10 << 20)
	API.UploadExtensions = sec.Key("api.upload_extensions").Strings("|")
	API.UploadContentTypes = sec.Key("api.upload_content_types").Strings("|")
	if len(API.UploadExtensions) == 0 {
		API.UploadExtensions = []string{".jpg", ".jpeg", ".gif", ".png", ".svg"}
	}
	if len(API.UploadContentTypes) == 0 {
		API.UploadContentTypes = []string{"image/jpeg", "image/png", "image/gif", "image/svg+xml"}
	}

	API.UploadFolder = strings.TrimSuffix(API.UploadFolder, "/")
	if !strings.HasPrefix(API.UploadFolder, "./") {
		if strings.HasPrefix(API.UploadFolder, "/") {
			API.UploadFolder = "." + API.UploadFolder
		} else {
			API.UploadFolder = "./" + API.UploadFolder
		}
	}

	API.FilesPath = strings.TrimSuffix(API.FilesPath, "/")
	if !strings.HasPrefix(API.FilesPath, "/") {
		API.UploadFolder = "/" + API.UploadFolder
	}

	// database settings
	sec = Cfg.Section("database")
	DbCfg.Type = sec.Key("db.type").In("mysql", []string{"mysql"})
	DbCfg.Host = sec.Key("db.host").Value()
	DbCfg.Name = sec.Key("db.name").Value()
	DbCfg.User = sec.Key("db.user").Value()
	DbCfg.Passwd = sec.Key("db.passwd").Value()
	DbCfg.SSLMode = sec.Key("db.ssl").MustString("disable")
	DbCfg.Path = sec.Key("db.path").MustString("data/goant.db")

	// session settings
	sec = Cfg.Section("session")
	Session.Name = sec.Key("name").MustString("goant")
	Session.Key = sec.Key("key").MustString("goant")
	Session.Store = sec.Key("store").In("cookie", []string{"cookie"})

	sec = Cfg.Section("security")
	InstallLock = sec.Key("install_lock").MustBool(false)
}
