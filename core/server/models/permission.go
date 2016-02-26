package models

type Permission struct {
	Id         uint32 `json:"id" gorm:"primary_key"`
	Name       string `json:"name" sql:"not null;type:varchar(20)"`
	Slug       string `json:"slug" sql:"not null;type:varchar(20)"`
	ObjectType string `json:"-" sql:"not null;type:varchar(20)"`
	IsCore     bool   `json:"is_core" sql:"not null;type:tinyint(1)"`
}

// CreatePermission creates a permission
func CreatePermission(perm *Permission, opts *Options) ApiErr {
	count := 0
	if db.Model(&Permission{}).Where("slug = ?", perm.Slug).Count(&count); count > 0 {
		return ApiMsg.ErrPermissionAlreadyExist
	}

	tx := db.Begin()
	if err := tx.Create(perm).Error; err != nil {
		tx.Rollback()
		return ApiMsg.SaveFail
	}
	tx.Commit()
	return ApiMsg.Created
}

// GetPermissions returns all permissions
func GetPermissions(opts *Options) ([]*Permission, ApiErr) {
	permissions := make([]*Permission, 0)
	dbInner := initDb(opts)
	if err := dbInner.Find(&permissions).Error; err != nil {
		return nil, ApiMsg.ErrRoleNotFound
	}
	return permissions, ApiMsg.Success
}

// GetPermissionsByIds returns all permissions by given ids
func GetPermissionsByIds(ids []string) ([]*Permission, ApiErr) {
	permissions := make([]*Permission, 0)
	if err := db.Where("id in (?)", ids).Order("Id ASC").Find(&permissions).Error; err != nil {
		return nil, ApiMsg.ErrRoleNotFound
	}
	return permissions, ApiMsg.Success
}
