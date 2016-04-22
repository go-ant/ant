package models

import (
	"time"
	"unicode/utf8"

	"github.com/go-ant/ant/core/server/modules/capabilities"
	"github.com/go-ant/ant/core/server/modules/utils/slug"
)

type Role struct {
	Id          uint32    `json:"id" gorm:"primary_key"`
	Name        string    `json:"name" sql:"not null;type:varchar(50);"`
	Description string    `json:"description" sql:"type:varchar(150)"`
	Slug        string    `json:"slug" sql:"type:varchar(100);"`
	CreatedAt   time.Time `json:"created_at" sql:"type:datetime;"`
	CreatedBy   uint32    `json:"-" sql:"type:bigint unsigned;"`
	UpdatedAt   time.Time `json:"updated_at" sql:"type:datetime;"`
	UpdatedBy   uint32    `json:"-" sql:"type:bigint unsigned;"`

	Permissions []Permission `json:"permissions" gorm:"many2many:permissions_roles;"`
}

func (c *Role) GetPermissions() {
	db.Model(&c).Related(&c.Permissions, "Permissions")
}

func CreateRole(role *Role, opts *Options) ApiErr {
	if existed, existRole, _ := IsRoleExist(role.Name); existed && existRole.Id != role.Id {
		return ApiMsg.ErrRoleAlreadyExist
	}

	if len(role.Name) == 0 {
		return ApiMsg.ErrRoleNameCanNotBeEmpty
	}

	if utf8.RuneCountInString(role.Name) > 50 {
		return ApiMsg.ErrRoleNameTooLong
	}

	if utf8.RuneCountInString(role.Description) > 150 {
		return ApiMsg.ErrRoleDescriptionTooLong
	}

	if len(opts.Permissions) == 0 {
		return ApiMsg.ErrPermissionCanNotBeEmpty
	}

	role.Slug = slug.Make(role.Name)

	tx := db.Begin()
	if err := tx.Create(role).Error; err != nil {
		tx.Rollback()
		return UnknowError(err.Error())
	}

	for _, permission := range opts.Permissions {
		tx.Exec("INSERT INTO `permissions_roles` (`permission_id`, `role_id`) VALUES (?, ?)", permission.Id, role.Id)
	}

	tx.Commit()

	permissionsToAdd := make([]string, len(opts.Permissions))
	for _, perm := range opts.Permissions {
		permissionsToAdd = append(permissionsToAdd, perm.Slug)
	}
	capabilities.SetRole(role.Slug, permissionsToAdd)

	return ApiMsg.Created
}

func EditRole(role *Role, opts *Options) ApiErr {
	if existed, existRole, _ := IsRoleExist(role.Name); existed && existRole.Id != role.Id {
		return ApiMsg.ErrRoleAlreadyExist
	}

	if len(role.Name) == 0 {
		return ApiMsg.ErrRoleNameCanNotBeEmpty
	}

	if utf8.RuneCountInString(role.Name) > 50 {
		return ApiMsg.ErrRoleNameTooLong
	}

	if utf8.RuneCountInString(role.Description) > 150 {
		return ApiMsg.ErrRoleDescriptionTooLong
	}

	if len(opts.Permissions) == 0 {
		return ApiMsg.ErrPermissionCanNotBeEmpty
	}

	role.Slug = slug.Make(role.Name)
	role.UpdatedAt = time.Now()

	tx := db.Begin()
	if err := tx.Select([]string{"name", "description", "slug", "updated_at", "updated_by"}).Save(role).Error; err != nil {
		tx.Rollback()
		return UnknowError(err.Error())
	}

	tx.Exec("DELETE FROM `permissions_roles` WHERE `role_id` = ?", role.Id)
	for _, permission := range opts.Permissions {
		tx.Exec("INSERT INTO `permissions_roles` (`permission_id`, `role_id`) VALUES (?, ?)", permission.Id, role.Id)
	}

	tx.Commit()

	permissionsToAdd := make([]string, len(opts.Permissions))
	for _, perm := range opts.Permissions {
		permissionsToAdd = append(permissionsToAdd, perm.Slug)
	}
	capabilities.SetRole(role.Slug, permissionsToAdd)

	return ApiMsg.Saved
}

func DeleteRole(id uint32, opts *Options) ApiErr {
	// Note: A Role with users cannot perform delete operation.
	count := 0
	if db.Table("roles_users").Where("role_id = ?", id).Count(&count); count > 0 {
		return ApiMsg.ErrRoleHaveUsersCanNotBeDeleted
	}

	tx := db.Begin()
	tx.Where("id = ?", id).Delete(Role{})
	tx.Commit()

	return ApiMsg.Deleted
}

func GetRole(opts *Options) (*Role, ApiErr) {
	role := new(Role)
	dbInner := initDb(opts)
	dbInner.First(&role)

	if opts != nil && opts.IsInclude("permission") {
		role.GetPermissions()
	}

	if role.Id == 0 {
		return nil, ApiMsg.Success
	}
	return role, ApiMsg.Success
}

func GetRoleById(id uint32, opts *Options) (*Role, ApiErr) {
	role := new(Role)
	if u := db.First(role, id); u.Error != nil {
		return nil, ApiMsg.ErrRoleNotFound
	}

	if opts != nil && opts.IsInclude("permission") {
		role.GetPermissions()
	}
	return role, ApiMsg.Success
}

func GetRoles(opts *Options) ([]*Role, ApiErr) {
	dbInner := initDb(opts)
	roles := make([]*Role, 0)
	if err := dbInner.Find(&roles).Error; err != nil {
		return nil, UnknowError(err.Error())
	}
	return roles, ApiMsg.Success
}

func IsRoleExist(name string) (bool, *Role, ApiErr) {
	role := new(Role)
	if len(name) == 0 {
		return false, role, ApiMsg.Success
	}

	role.Name = name
	u := db.Where(&Role{Name: name}).First(role)
	if u.Error != nil {
		return false, role, ApiMsg.Success
	}
	return !u.RecordNotFound(), role, ApiMsg.Success
}
