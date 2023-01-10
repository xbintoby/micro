package model

import (
	"jam3.com/common/errs"
)

var (
	UsernameOrPwd = errs.NewError(1000, "username or password is null")
	NoLegal       = errs.NewError(1001, "invalid") // invalid param
	NoLegalUid    = errs.NewError(1002, "invalid uid")
	TokenIsNull   = errs.NewError(1003, "Token is null")
	JwtAuthFail   = errs.NewError(1004, "Jwt auth fail")
)
