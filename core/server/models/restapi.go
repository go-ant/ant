package models

type RestApi struct {
	Data       interface{} `json:"data"`
	Error      ApiErr      `json:"error"`
	Pagination *Pagination `json:"pagination"`
}

type ApiErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (c *ApiErr) IsSuccess() bool {
	return c.Code == 0
}

type apiMsg struct {
	Success ApiErr
	Created ApiErr
	Saved   ApiErr
	Deleted ApiErr

	DatabaseFailed ApiErr
	NeedToSignIn   ApiErr
	NoPermission   ApiErr
	SaveFail       ApiErr
	LoadFail       ApiErr

	ErrRoleNameCanNotBeEmpty        ApiErr
	ErrRoleNameTooLong              ApiErr
	ErrRoleDescriptionTooLong       ApiErr
	ErrNotSupportNumericRole        ApiErr
	ErrRoleCanNotBeEmpty            ApiErr
	ErrRoleAlreadyExist             ApiErr
	ErrRoleNotFound                 ApiErr
	ErrRoleHaveUsersCanNotBeDeleted ApiErr
	ErrPermissionCanNotBeEmpty      ApiErr
	ErrPermissionAlreadyExist       ApiErr

	ErrUserNameCanNotBeEmpty     ApiErr
	ErrUserNameTooLong           ApiErr
	ErrNotSupportNumericUser     ApiErr
	ErrUserPasswordCanNotBeEmpty ApiErr
	ErrUserPasswordTooLong       ApiErr
	ErrUserAlreadyExist          ApiErr
	ErrUserNotFound              ApiErr
	ErrPasswordIncorrect         ApiErr
	ErrPasswordNotMatch          ApiErr
	ErrPasswordRequired          ApiErr
	ErrPasswordTooShort          ApiErr

	ErrPostTitleTooLong           ApiErr
	ErrPostSlugTooLong            ApiErr
	ErrPostMetaTitleTooLong       ApiErr
	ErrPostMetaDescriptionTooLong ApiErr
	ErrPostSlugAlreadyExist       ApiErr
	ErrPostNotFound               ApiErr

	ErrFileNotSupported ApiErr
	ErrFileTooLarge     ApiErr
}

var ApiMsg = apiMsg{
	Success: ApiErr{Code: 0, Message: ""},
	Created: ApiErr{Code: 0, Message: "msg.created"},
	Saved:   ApiErr{Code: 0, Message: "msg.saved"},
	Deleted: ApiErr{Code: 0, Message: "msg.deleted"},

	DatabaseFailed: ApiErr{Code: 1, Message: "msg.db_initialization_failed"},
	NeedToSignIn:   ApiErr{Code: 2, Message: "msg.need_to_sign_in"},
	NoPermission:   ApiErr{Code: 9, Message: "msg.no_permission"},
	SaveFail:       ApiErr{Code: 10, Message: "msg.save_fail"},
	LoadFail:       ApiErr{Code: 10, Message: "msg.load_fail"},

	ErrRoleNameCanNotBeEmpty:        ApiErr{Code: 101, Message: "msg.role_name_not_be_empty"},
	ErrRoleNameTooLong:              ApiErr{Code: 102, Message: "msg.role_name_too_long"},
	ErrNotSupportNumericRole:        ApiErr{Code: 103, Message: "msg.not_support_numeric_role"},
	ErrRoleDescriptionTooLong:       ApiErr{Code: 104, Message: "msg.role_description_too_long"},
	ErrRoleCanNotBeEmpty:            ApiErr{Code: 105, Message: "msg.role_not_be_empty"},
	ErrRoleAlreadyExist:             ApiErr{Code: 106, Message: "msg.role_exist"},
	ErrRoleNotFound:                 ApiErr{Code: 107, Message: "msg.role_not_found"},
	ErrRoleHaveUsersCanNotBeDeleted: ApiErr{Code: 108, Message: "msg.role_have_users_not_be_deleted"},
	ErrPermissionCanNotBeEmpty:      ApiErr{Code: 109, Message: "msg.perm_not_be_empty"},
	ErrPermissionAlreadyExist:       ApiErr{Code: 110, Message: "msg.perm_exist"},

	ErrUserNameCanNotBeEmpty:     ApiErr{Code: 201, Message: "msg.user_name_not_be_empty"},
	ErrUserNameTooLong:           ApiErr{Code: 202, Message: "msg.user_name_too_long"},
	ErrNotSupportNumericUser:     ApiErr{Code: 203, Message: "msg.not_support_numeric_user"},
	ErrUserPasswordCanNotBeEmpty: ApiErr{Code: 204, Message: "msg.user_password_not_be_empty"},
	ErrUserPasswordTooLong:       ApiErr{Code: 205, Message: "msg.user_password_too_long"},
	ErrUserAlreadyExist:          ApiErr{Code: 206, Message: "msg.user_exist"},
	ErrUserNotFound:              ApiErr{Code: 207, Message: "msg.user_not_found"},
	ErrPasswordIncorrect:         ApiErr{Code: 208, Message: "msg.password_incorrect"},
	ErrPasswordNotMatch:          ApiErr{Code: 209, Message: "msg.passwords_not_match"},
	ErrPasswordRequired:          ApiErr{Code: 210, Message: "msg.password_required"},
	ErrPasswordTooShort:          ApiErr{Code: 211, Message: "msg.password_at_least_8_char"},

	ErrPostTitleTooLong:           ApiErr{Code: 301, Message: "msg.post_title_too_long"},
	ErrPostSlugTooLong:            ApiErr{Code: 302, Message: "msg.post_slug_too_long"},
	ErrPostMetaTitleTooLong:       ApiErr{Code: 303, Message: "msg.post_meta_title_too_long"},
	ErrPostMetaDescriptionTooLong: ApiErr{Code: 304, Message: "msg.post_meta_description_too_long"},
	ErrPostSlugAlreadyExist:       ApiErr{Code: 305, Message: "msg.post_slug_exist"},
	ErrPostNotFound:               ApiErr{Code: 306, Message: "msg.post_not_found"},

	ErrFileNotSupported: ApiErr{Code: 901, Message: "msg.file_not_supported"},
	ErrFileTooLarge:     ApiErr{Code: 902, Message: "msg.file_too_large"},
}

func UnknowError(msg string) ApiErr {
	return ApiErr{Code: 10, Message: msg}
}
