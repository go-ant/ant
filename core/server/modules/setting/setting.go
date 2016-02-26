package setting

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	CONFIG_PATH = "config/app.ini"
)

var (
	EnableGzip bool
	Installed  bool

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
		Type, Host, Name, User, Passwd, SSLMode string
	}
)

func init() {
	Config, err := loadConfig(CONFIG_PATH)
	if err != nil || Config == nil {
		panic(fmt.Errorf("Err: Failed to load `config/app.ini`"))
	}

	EnableGzip, _ = strconv.ParseBool(Config.StringDefault("enable_gzip", "true"))
	Host.Path = Config.StringDefault("host.path", "/")
	Host.Addr = Config.StringDefault("host.addr", "")
	Host.Port = Config.StringDefault("host.port", "2015")

	if !strings.HasPrefix(Host.Path, "/") {
		Host.Path = "/" + Host.Path
	}
	Host.Path = strings.TrimSuffix(Host.Path, "/")

	// API settings
	Config.Section("api")
	API.UploadFolder = Config.StringDefault("api.upload_folder", "./content/upload")
	API.FilesPath = Config.StringDefault("api.files_path", "/upload")
	API.UploadMaxSize = Config.IntDefault("api.upload_max_size", 10<<20)
	API.UploadExtensions = Config.ArrayDefault("api.upload_extensions", []string{".jpg", ".jpeg", ".gif", ".png", ".svg"})
	API.UploadContentTypes = Config.ArrayDefault("api.upload_content_types", []string{"image/jpeg", "image/png", "image/gif", "image/svg+xml"})

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
	Config.Section("database")
	DbCfg.Type = "mysql"
	DbCfg.Host = Config.StringDefault("db.host", "")
	DbCfg.Name = Config.StringDefault("db.name", "")
	DbCfg.User = Config.StringDefault("db.user", "")
	DbCfg.Passwd = Config.StringDefault("db.passwd", "")
	DbCfg.SSLMode = Config.StringDefault("db.ssl", "disable")

}
