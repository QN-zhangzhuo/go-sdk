package account

import "time"

const (
	// USER_TYPE_ADMIN 管理员账户
	USER_TYPE_ADMIN        = 0x0001
	USER_TYPE_VIP          = 0x0002
	USER_TYPE_STDUSER      = 0x0004
	USER_TYPE_STDUSER2     = 0x0008
	USER_TYPE_EXPUSER      = 0x0010
	USER_TYPE_PARENTUSER   = 0x0020
	USER_TYPE_OP           = 0x0040
	USER_TYPE_SUPPORT      = 0x0080
	USER_TYPE_CC           = 0x0100
	USER_TYPE_QCOS         = 0x0200
	USER_TYPE_FUSION       = 0x0400
	USER_TYPE_PILI         = 0x0800
	USER_TYPE_PANDORA      = 0x1000
	USER_TYPE_DISTRIBUTION = 0x2000
	USER_TYPE_QVM          = 0x4000
	USER_TYPE_DISABLED     = 0x8000

	USER_TYPE_USERS   = USER_TYPE_STDUSER | USER_TYPE_STDUSER2 | USER_TYPE_EXPUSER
	USER_TYPE_SUDOERS = USER_TYPE_ADMIN | USER_TYPE_OP | USER_TYPE_SUPPORT
)

// UserInfo 用户账户的信息
type UserInfo struct {
	UID                   uint32    `json:"uid,omitempty"`
	UserID                string    `json:"userid,omitempty"`
	Email                 string    `json:"email,omitempty"`
	Username              string    `json:"username,omitempty"`
	ParentUID             uint32    `json:"parent_uid,omitempty"`
	IsActivated           bool      `json:"is_activated,omitempty"`
	UserType              UserType  `json:"user_type,omitempty"`
	DeviceNum             int       `json:"device_num,omitempty"`
	InvitationNum         int       `json:"invitation_num,omitempty"`
	LastParentOperationAt time.Time `json:"last_parent_operation_at,omitempty"`
}

// SudoerInfo 七牛内部用户信息， 通过ldap登陆
type SudoerInfo struct {
	UserInfo
	Sudoer  uint32 `json:"suid,omitempty"`
	UtypeSu uint32 `json:"sut,omitempty"`
}

// IsAdmin 判断用户是否是管理员
func (u *UserInfo) IsAdmin() bool {
	return u.UserType.IsAdmin()
}

// UserType 七牛账户的信息
type UserType uint32

// IsAdmin 判断账户是否是管理员账户
func (u UserType) IsAdmin() bool {
	return u&USER_TYPE_ADMIN != 0
}

// DeveloperInfo 开发者信息
type DeveloperInfo struct {
	UID          uint32 `json:"uid"`
	Email        string `json:"email"`
	FullName     string `json:"fullname"`
	IsEnterprise bool   `json:"is_enterprise"`
	IsCertified  bool   `json:"is_certified"`
	IsInternal   bool   `json:"is_internal"`
}
