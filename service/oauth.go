package service

// Config 包含了第三方应用的信息
type Config struct {
	// ClientID 应用的ID
	ClientID string

	// ClientSecret 应用的Secret
	ClientSecret string

	// Endpoint 包含了资源拥有者的token 服务器的URL
	Endpoint Endpoint

	// RedirectURL 重定向地址
	RedirectURL string

	// Scope 指定可选的而外的权限
	Scopes []string
}

// Endpoint 代表了OAuth 2.0 服务的提供者的认证和token服务器的URL
type Endpoint struct {
	AuthURL  string
	TokenURL string

	// AuthStyle 可选，指示client_id, client_secret如何传递给服务器
	AuthStyle AuthStyle
}

// AuthStyle 指定获取token的请求如果传递给认证服务器
type AuthStyle int

const (
	// AuthStyleAutoDetect 自动检测认证服务器的认证方式， 并且后续使用同样的方法传递client_id, client_secret
	AuthStyleAutoDetect AuthStyle = 0

	// AuthStyleInParams 表明在 POST 请求体中发送client_id, client_secret
	// 并且Content-Type 类型为application/x-www-form-urlencoded
	AuthStyleInParams AuthStyle = 1

	// AuthStyleInHeader 使用http Basic Authorization方式， 发送client_id, client_password
	// 该方式在OAuth2 RFC 6749 section 2.31中有说明
	AuthStyleInHeader AuthStyle = 2
)
