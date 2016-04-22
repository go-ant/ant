package models

import (
	"crypto/sha256"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/go-ant/ant/core/server/modules/capabilities"
	"github.com/go-ant/ant/core/server/modules/capabilities/rbac"
	"github.com/go-ant/ant/core/server/modules/setting"
	"github.com/go-ant/ant/core/server/modules/utils"
	"github.com/go-ant/ant/core/server/modules/utils/slug"
)

type User struct {
	Id        uint32    `json:"id" gorm:"primary_key"`
	Name      string    `json:"name" sql:"type:varchar(50);"`
	Password  string    `json:"-" sql:"type:char(40);"`
	Slug      string    `json:"slug" sql:"type:varchar(100);unique_index;"`
	Email     string    `json:"email" sql:"type:varchar(50);"`
	Bio       string    `json:"bio" sql:"type:varchar(200);"`
	Avatar    string    `json:"avatar" sql:"type:varchar(200);"`
	Cover     string    `json:"cover" sql:"type:varchar(200);"`
	Website   string    `json:"website" sql:"type:varchar(200);"`
	Location  string    `json:"location" sql:"type:varchar(200);"`
	Language  string    `json:"language" sql:"type:varchar(6);"`
	Salt      string    `json:"-" sql:"type:char(10);"`
	LastLogin time.Time `json:"last_login" sql:"type:datetime;"`
	Status    string    `json:"status" sql:"type:varchar(20);"`
	CreatedAt time.Time `json:"-" sql:"type:datetime;"`
	CreatedBy uint32    `json:"-" sql:"type:bigint unsigned;"`
	UpdatedAt time.Time `json:"-" sql:"type:datetime;"`
	UpdatedBy uint32    `json:"-" sql:"type:bigint unsigned;"`

	Roles []Role `json:"roles" gorm:"many2many:roles_users;"`
}

// EncodePassword encodes password to safe format.
func (c *User) EncodePassword() {
	newPasswd := utils.PBKDF2([]byte(c.Password), []byte(c.Salt), 5000, 20, sha256.New)
	c.Password = fmt.Sprintf("%x", newPasswd)
}

func (c *User) GetRole() Role {
	if c.Roles == nil {
		db.Model(&c).Related(&c.Roles, "Roles")
	}
	return c.Roles[0]
}

const (
	SiteOwner string = "owner"
)

// IsGranted check if user has permission
func (c *User) IsGranted(permission string, assert rbac.AssertionFunc) bool {
	return capabilities.IsGranted(c.GetRole().Slug, permission, assert)
}

func (c *User) IsOwner() bool {
	return c.GetRole().Slug == SiteOwner
}

func CreateUser(u *User, opts *Options) ApiErr {

	if utils.IsEmpty(u.Name) {
		return ApiMsg.ErrUserNameCanNotBeEmpty
	}

	if utf8.RuneCountInString(u.Name) > 50 {
		return ApiMsg.ErrUserNameTooLong
	}

	if utils.IsNumeric(u.Name) {
		return ApiMsg.ErrNotSupportNumericUser
	}

	if utils.IsEmpty(u.Password) {
		return ApiMsg.ErrUserPasswordCanNotBeEmpty
	}

	if utf8.RuneCountInString(u.Password) < 8 {
		return ApiMsg.ErrPasswordTooShort
	}

	if utf8.RuneCountInString(u.Password) > 20 {
		return ApiMsg.ErrUserPasswordTooLong
	}

	if opts.Role.Id == 0 {
		return ApiMsg.ErrRoleCanNotBeEmpty
	}

	targetRole, _ := GetRoleById(opts.Role.Id, nil)
	if targetRole.Id == 0 {
		return ApiMsg.ErrRoleNotFound
	}

	// can not add user to owner-role if the owner exist
	if targetRole.Slug == SiteOwner {
		if _, recordCount, _ := GetUsers(nil); recordCount > 0 {
			return ApiMsg.NoPermission
		}
	}

	if existed, _, _ := isUserExist(u.Name); existed {
		return ApiMsg.ErrUserAlreadyExist
	}

	// default avatar
	if utils.IsEmpty(u.Avatar) {
		u.Avatar = setting.Host.Path + "/goant/assets/css/images/avatar.png"
	}

	u.Cover = setting.Host.Path + "/goant/assets/css/images/cover.jpg"
	u.Slug = slug.Make(u.Name)
	u.Salt = utils.GetRandomString(10)
	u.Status = "active"
	u.LastLogin = utils.ToTime(time.RFC1123)
	u.EncodePassword()

	tx := db.Begin()
	if err := tx.Create(u).Error; err != nil {
		tx.Rollback()
		return UnknowError(err.Error())
	}

	tx.Exec("INSERT INTO `roles_users` (`role_id`, `user_id`) VALUES (?, ?)", targetRole.Id, u.Id)

	tx.Commit()

	return ApiMsg.Created
}

func EditUser(u *User, opts *Options) ApiErr {

	if utils.IsEmpty(u.Name) {
		return ApiMsg.ErrUserNameCanNotBeEmpty
	}

	if utf8.RuneCountInString(u.Name) > 50 {
		return ApiMsg.ErrUserNameTooLong
	}

	if utils.IsNumeric(u.Name) {
		return ApiMsg.ErrNotSupportNumericUser
	}

	if opts.Role.Id == 0 {
		return ApiMsg.ErrRoleCanNotBeEmpty
	}

	// just owner can update itself
	if u.GetRole().Slug == SiteOwner && u.Id != opts.User.Id {
		return ApiMsg.NoPermission
	}

	targetRole, _ := GetRoleById(opts.Role.Id, nil)
	if targetRole.Id == 0 {
		return ApiMsg.ErrRoleNotFound
	}

	// can't update user to "Owner"
	if u.GetRole().Slug != SiteOwner && targetRole.Slug == SiteOwner {
		return ApiMsg.NoPermission
	}

	if existed, existUser, _ := isUserExist(u.Name); existed && existUser.Id != u.Id {
		return ApiMsg.ErrUserAlreadyExist
	}

	if utils.IsEmpty(u.Avatar) {
		u.Avatar = setting.Host.Path + "/goant/assets/css/images/avatar.png"
	}
	if utils.IsEmpty(u.Cover) {
		u.Cover = setting.Host.Path + "/goant/assets/css/images/cover.jpg"
	}

	u.Slug = slug.Make(u.Name)

	tx := db.Begin()
	if err := tx.Select([]string{"name", "slug", "email", "location", "website", "bio", "avatar", "cover", "updated_at", "updated_by"}).Save(u).Error; err != nil {
		tx.Rollback()
		return UnknowError(err.Error())
	}
	tx.Exec("UPDATE `roles_users` SET `role_id` = ? WHERE (`user_id` = ?)", targetRole.Id, u.Id)
	tx.Commit()

	return ApiMsg.Saved
}

func GetUsers(opts *Options) ([]*User, uint32, ApiErr) {
	var recordCount uint32
	if opts == nil {
		opts = &Options{Limit: 15, Page: 1}
	}
	dbInit := initDb(opts)
	users := make([]*User, 0)

	if err := dbInit.Model(User{}).Count(&recordCount).Limit(opts.Limit).Offset(opts.Offset).Find(&users).Error; err != nil {
		return nil, 0, UnknowError(err.Error())
	}
	if len(users) > 0 && opts.IsInclude("role") {
		for _, user := range users {
			db.Model(&user).Association("Roles").Find(&user.Roles)
		}
	}

	return users, recordCount, ApiMsg.Success
}

func GetUser(opts *Options) (*User, ApiErr) {
	dbInit := initDb(opts)
	user := new(User)
	if u := dbInit.First(user); u.Error != nil {
		return nil, ApiMsg.ErrUserNotFound
	}

	if opts != nil && opts.IsInclude("role") && user.Id > 0 {
		user.GetRole()
	}
	return user, ApiMsg.Success
}

func GetUserById(id uint32, opts *Options) (*User, ApiErr) {
	if opts == nil {
		opts = &Options{}
	}
	opts.GormAdp = &GormAdapter{
		Query: "id = ?",
		Args:  []interface{}{id},
	}
	user, errApi := GetUser(opts)

	return user, errApi
}

func GetUserByName(name string, opts *Options) (*User, ApiErr) {
	if opts == nil {
		opts = &Options{}
	}
	opts.GormAdp = &GormAdapter{
		Query: "name = ?",
		Args:  []interface{}{name},
	}
	user, errApi := GetUser(opts)

	return user, errApi
}

func ChangePassword(oldPassword, newPassword, verifyPassword string, userId uint32, opts *Options) (bool, ApiErr) {

	targetUser, _ := GetUserById(userId, nil)
	if targetUser.Id == 0 {
		return false, ApiMsg.ErrUserNotFound
	}

	if targetUser.Id != opts.User.Id && (!opts.User.IsGranted("edit-users", nil) || targetUser.GetRole().Slug == SiteOwner) {
		return false, ApiMsg.NoPermission
	}

	if newPassword != verifyPassword {
		return false, ApiMsg.ErrPasswordNotMatch
	}

	requirePassword := opts != nil && opts.User != nil && userId == opts.User.Id

	if requirePassword && utils.IsEmpty(oldPassword) {
		return false, ApiMsg.ErrPasswordRequired
	}

	if utf8.RuneCountInString(newPassword) < 8 {
		return false, ApiMsg.ErrPasswordTooShort
	}

	if utf8.RuneCountInString(newPassword) > 20 {
		return false, ApiMsg.ErrUserPasswordTooLong
	}

	targetUser, err := GetUserById(userId, nil)

	if err.Code != 0 {
		return false, ApiMsg.ErrUserNotFound
	}

	if requirePassword && !passwordCompare(oldPassword, targetUser.Password, targetUser.Salt) {
		return false, ApiMsg.ErrPasswordIncorrect
	}

	targetUser.Password = newPassword
	targetUser.EncodePassword()
	tx := db.Begin()
	if err := tx.Model(&targetUser).UpdateColumns(User{Password: targetUser.Password}).Error; err != nil {
		tx.Rollback()
		return false, UnknowError(err.Error())
	}
	tx.Commit()
	return true, ApiMsg.Saved
}

func UserSignin(u *User, opts *Options) (*User, ApiErr) {
	user, err := GetUserByName(u.Name, opts)
	if err.Code != 0 {
		return nil, ApiMsg.ErrUserNotFound
	}
	if !passwordCompare(u.Password, user.Password, user.Salt) {
		return nil, ApiMsg.ErrPasswordIncorrect
	}

	db.Model(user).UpdateColumn(User{LastLogin: time.Now()})

	return user, ApiMsg.Success
}

func isUserExist(name string) (bool, *User, error) {
	user := new(User)
	if len(name) == 0 {
		return false, user, nil
	}
	if err := db.Where(&User{Name: name}).First(user).Error; err != nil {
		return false, nil, err
	} else {
		return user.Id > 0, user, nil
	}
}

func passwordCompare(password, source, salt string) bool {
	newPasswd := utils.PBKDF2([]byte(password), []byte(salt), 5000, 20, sha256.New)
	return source == fmt.Sprintf("%x", newPasswd)
}
