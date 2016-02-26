package models

func InitialDatabase() ApiErr {
	tx := db.Begin()

	tx.DropTableIfExists(&User{})
	tx.DropTableIfExists(&Role{})
	tx.DropTableIfExists(&Permission{})
	tx.DropTableIfExists(&Post{})
	tx.DropTableIfExists(&Setting{})
	tx.DropTableIfExists("permissions_roles")
	tx.DropTableIfExists("roles_users")

	if err := tx.CreateTable(&User{}).Error; err != nil {
		tx.Rollback()
		return ApiMsg.SaveFail
	}
	if err := tx.CreateTable(&Role{}).Error; err != nil {
		tx.Rollback()
		return ApiMsg.SaveFail
	}
	if err := tx.CreateTable(&Permission{}).Error; err != nil {
		tx.Rollback()
		return ApiMsg.SaveFail
	}
	if err := tx.CreateTable(&Post{}).Error; err != nil {
		tx.Rollback()
		return ApiMsg.SaveFail
	}
	if err := tx.CreateTable(&Setting{}).Error; err != nil {
		tx.Rollback()
		return ApiMsg.SaveFail
	}


	tx.Commit()

	return ApiMsg.Success
}
