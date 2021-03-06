// Package defaults 包含了一些帮助函数来获取SDK的默认配置和默认的handlers
//
// 一般情况下，这个包不应该被直接使用，使用Session初始化的时候会调用该包配置默认项
package defaults

import (
	"net/http"

	"github.com/QN-zhangzhuo/go-sdk/qiniu"
	"github.com/QN-zhangzhuo/go-sdk/qiniu/corehandlers"
	"github.com/QN-zhangzhuo/go-sdk/qiniu/credentials"
	"github.com/QN-zhangzhuo/go-sdk/qiniu/defs"
	"github.com/QN-zhangzhuo/go-sdk/qiniu/request"
)

// Defaults 提供了默认的配置选项和handlers
type Defaults struct {
	Config   *qiniu.Config
	Handlers request.Handlers
}

// Get 获取默认的配置
func Get() Defaults {
	cfg := Config()
	handlers := Handlers()

	return Defaults{
		Config:   cfg,
		Handlers: handlers,
	}
}

// Config 返回默认的配置， 密钥信息默认是没有设置的
func Config() *qiniu.Config {
	return qiniu.NewConfig().
		WithHTTPClient(http.DefaultClient).
		WithLogger(qiniu.NewDefaultLogger()).
		WithLogLevel(qiniu.LogOff).
		WithRsHost(defs.DefaultRsHost).
		WithRsfHost(defs.DefaultRsfHost).
		WithAPIHost(defs.DefaultAPIHost).
		WithUCHost(defs.DefaultUcHost).
		WithMorseHost(defs.DefaultMorseHost).
		WithEmailClientID(defs.DefaultEmailClientID)
}

// Handlers 返回默认的请求handlers
func Handlers() request.Handlers {
	var handlers request.Handlers

	handlers.Build.PushBackNamed(corehandlers.SDKVersionUserAgentHandler)
	handlers.Build.PushBackNamed(corehandlers.AddHostExecEnvUserAgentHander)
	handlers.Build.PushBackNamed(corehandlers.BodyHandler)
	handlers.Build.AfterEachFn = request.HandlerListStopOnError
	handlers.Sign.PushBackNamed(corehandlers.BuildContentLengthHandler)
	handlers.Sign.AfterEachFn = request.HandlerListStopOnError
	handlers.Send.PushBackNamed(corehandlers.SendHandler)
	handlers.AfterRetry.PushBackNamed(corehandlers.AfterRetryHandler)
	handlers.ValidateResponse.PushBackNamed(corehandlers.ValidateResponseHandler)
	handlers.Unmarshal.PushBackNamed(corehandlers.UnmarshalHandler)
	handlers.Complete.PushBackNamed(corehandlers.CompleteHandler)

	return handlers
}

// CredChain 返回默认的密钥配置链
func CredChain(cfg *qiniu.Config, handlers request.Handlers) *credentials.Credentials {
	return credentials.NewCredentials(&credentials.ChainProvider{
		VerboseErrors: qiniu.BoolValue(cfg.CredentialsChainVerboseErrors),
		Providers:     CredProviders(cfg, handlers),
	})
}

// CredProviders 提供链式provider
func CredProviders(cfg *qiniu.Config, handlers request.Handlers) []credentials.Provider {
	return []credentials.Provider{
		&credentials.EnvProvider{},
	}
}
