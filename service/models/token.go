package models

import (
	"net/http"
	"strings"
	"time"
)

// Token 是Oauth 2.0 服务需要的请求资源的认证token
//
// 大部分用户不应该直接使用该结构体中的字段， 该字段可导出主要是给其他实现Oauth 2.0流程的包使用
type Token struct {
	// AccessToken 用来认证和授权请求访问资源
	AccessToken string `json:"access_token"`

	// TokenType 是token的类型，默认是"Bearer"
	TokenType string `json:"token_type,omitempty"`

	// RefreshToken 用来刷新access token
	RefreshToken string `json:"refresh_token,omitempty"`

	ExpiresIN int `json:"expires_in"`

	// Expiry 是可选的表明access token过期的时间
	//
	// 如果是0， TokenSource 实现会一直使用access token和Refresh token
	Expiry time.Time `json:"expiry,omitempty"`
}

// SetExpiry 通过整型设置过期时间， expiresIn 表示多少秒内过期
// 七牛oauth2/token接口返回的值是expiresIn, 那到该值后需要理解调用该函数
func (t *Token) SetExpiry(expiresIn int) {
	t.Expiry = time.Now().Add(time.Second * time.Duration(expiresIn))
}

// Type 如果t.TokenType非空，返回该字段， 否则返回"Bearer"
func (t *Token) Type() string {
	if strings.EqualFold(t.TokenType, "bearer") {
		return "Bearer"
	}
	if strings.EqualFold(t.TokenType, "mac") {
		return "MAC"
	}
	if strings.EqualFold(t.TokenType, "basic") {
		return "Basic"
	}
	if t.TokenType != "" {
		return t.TokenType
	}
	return "Bearer"
}

// SetAuthHeader 设置http Authorization Header
func (t *Token) SetAuthHeader(r *http.Request) {
	r.Header.Set("Authorization", t.Type()+" "+t.AccessToken)
}

// Expired 检测AccessToken是否过期
func (t *Token) Expired() bool {
	if t.Expiry.IsZero() {
		return false
	}
	return time.Now().After(t.Expiry)
}

// Valid 验证accessToken是否有效
func (t *Token) Valid() bool {
	return t != nil && t.AccessToken != "" && !t.Expired()
}
