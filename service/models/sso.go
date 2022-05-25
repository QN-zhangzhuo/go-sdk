package models

// SSOUserInfo sso 用户信息
type SSOUserInfo struct {
	UID        uint32 `json:"uid"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	LoginToken string `json:"login_token"`
}
