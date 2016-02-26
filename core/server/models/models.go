package models

import (
	"fmt"
	"github.com/go-ant/ant/core/server/modules/setting"
	"github.com/go-ant/ant/core/server/modules/startup"
	"github.com/go-ant/ant/core/server/modules/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

const APP_VERSION = "0.1.0"

var db gorm.DB

func init() {
	startup.Register(func() {
		if err := newEngine(); err == nil {
			db.LogMode(true)
		}
	})
}

func newEngine() (err error) {
	err = getEngine()
	if err != nil {
		return fmt.Errorf("connect to database: %v", err)
	}

	return nil
}

func getEngine() (err error) {
	connectionString := ""
	switch setting.DbCfg.Type {
	case "mysql":
		connectionString = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", setting.DbCfg.User, setting.DbCfg.Passwd, setting.DbCfg.Host, setting.DbCfg.Name)
	case "postgres":
		var host, port = "127.0.0.1", "5432"
		fields := strings.Split(setting.DbCfg.Host, ":")
		if len(fields) > 0 && len(strings.TrimSpace(fields[0])) > 0 {
			host = fields[0]
		}
		if len(fields) > 1 && len(strings.TrimSpace(fields[1])) > 0 {
			port = fields[1]
		}
		connectionString = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s", setting.DbCfg.User, setting.DbCfg.Passwd, host, port, setting.DbCfg.Name, setting.DbCfg.SSLMode)
		//	case "sqlite3":
		//		if !EnableSQLite3 {
		//			return nil, fmt.Errorf("Unknown database type: %s", setting.DbCfg.Type)
		//		}
		//		os.MkdirAll(path.Dir(setting.DbCfg.Path), os.ModePerm)
		//		connectionString = "file:" + setting.DbCfg.Path + "?cache=shared&mode=rwc"
	default:
		return fmt.Errorf("Unknown database type: %s", setting.DbCfg.Type)
	}
	db, err = gorm.Open(setting.DbCfg.Type, connectionString)
	return
}

// Model Assistant
type Options struct {
	User        *User
	Role        *Role
	Permission  *Permission
	Permissions []*Permission

	Limit   uint32
	Page    uint32
	Offset  uint32
	Include string
	GormAdp *GormAdapter
}

func (c *Options) IsInclude(name string) bool {
	return utils.StringInSlice(name, strings.Split(c.Include, ","))
}

// GormAdapter
/*
 *	models.Options{
 *		GormAdp: &models.GormAdapter{Query: "column=value1"}},
 *	}
 *	models.Options{
 *		GormAdp: &models.GormAdapter{Query: "column=?", Args: []interface{}{"column-value"}},
 *	}
 *	models.Options{
 *		GormAdp: &models.GormAdapter{Map: map[string]interface{}{"column1": "value1", "column2": "value2"}},
 *	}
 */
type GormAdapter struct {
	Columns []string
	Query   interface{}
	Args    []interface{}
	Map     map[string]interface{}
	OrderBy string
	Joins   string
}

// initDb init db object
func initDb(opts *Options) *gorm.DB {
	dbInit := db.New()
	if opts != nil {
		if opts.Page == 0 {
			opts.Page = 1
		}
		opts.Offset = (opts.Page - 1) * opts.Limit
		if opts.GormAdp != nil {
			if opts.GormAdp.Query != nil {
				dbInit = db.Where(opts.GormAdp.Query, opts.GormAdp.Args...)
			} else if opts.GormAdp.Map != nil {
				dbInit = db.Where(opts.GormAdp.Map)
			}
			if len(opts.GormAdp.OrderBy) > 0 {
				dbInit = dbInit.Order(opts.GormAdp.OrderBy)
			} else {
				dbInit = dbInit.Order("id desc")
			}
			if !utils.IsEmpty(opts.GormAdp.Joins) {
				if opts.GormAdp.Joins == "role" {
					dbInit = dbInit.Joins("left join roles_users on roles_users.user_id = users.id left join roles on roles_users.role_id = roles.id")
				}
			}
			if len(opts.GormAdp.Columns) > 0 {
				dbInit = dbInit.Select(opts.GormAdp.Columns)
			}
		}
	}
	return dbInit
}
