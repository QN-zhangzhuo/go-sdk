package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/QN-zhangzhuo/go-sdk/qiniu"
	"github.com/QN-zhangzhuo/go-sdk/qiniu/account"
	"github.com/QN-zhangzhuo/go-sdk/qiniu/client"
	"github.com/QN-zhangzhuo/go-sdk/qiniu/request"
	"github.com/QN-zhangzhuo/go-sdk/qiniu/session"
	"github.com/QN-zhangzhuo/go-sdk/service/models"
)

const (
	// ServiceName 总的服务入口
	ServiceName = "Service"
)

// Service 是总的对外入口服务，集成了kodo, cdn, oauth, account等服务的功能
type Service struct {
	*client.BaseClient

	token *models.Token
	mu    sync.Mutex
}

// New 新建一个服务实例
func New() *Service {
	sess := session.Must(session.New())
	c := sess.ClientConfig()
	return &Service{
		BaseClient: client.New(
			*c.Config,
			c.Handlers,
		),
	}
}

// NewService 使用ConfigProvider 新建一个Service实例
func NewService(p client.ConfigProvider, cfgs ...*qiniu.Config) *Service {
	c := p.ClientConfig(cfgs...)
	svc := &Service{
		BaseClient: client.New(
			*c.Config,
			c.Handlers,
		),
	}
	return svc
}

// QualifiedAccessToken 验证accessToken是否合法， 合法返回true, 否则返回false
func (s *Service) QualifiedAccessToken(accessToken string) (bool, *account.UserInfo, error) {
	userInfo, err := s.UserInfoFromAccessToken(accessToken)
	if err != nil {
		return false, nil, err
	}
	if userInfo.IsAdmin() && userInfo.IsActivated {
		return true, userInfo, nil
	}
	return false, nil, nil
}

// UserInfoFromAccessToken 使用bearer access token查询用户的基本信息
func (s *Service) UserInfoFromAccessToken(accessToken string) (*account.UserInfo, error) {
	api := &request.API{
		Method:      "POST",
		Path:        "/user/info",
		Host:        qiniu.StringValue(s.Config.AccHost),
		ContentType: "application/x-www-form-urlencoded",
		ServiceName: ServiceName,
		APIName:     "UserInfo",
	}
	var info account.UserInfo
	v := url.Values{}
	v.Set("access_token", accessToken)

	req := s.NewRequest(api, strings.NewReader(v.Encode()), &info)

	return &info, req.Send()
}

// ProductOrderAccomplish 发货完成后，通知BO系统发货完成
func (s *Service) ProductOrderAccomplish(input *models.ReqProductOrderAccomplish, out interface{}) error {
	v := url.Values{}
	v.Set("id", strconv.FormatInt(input.ID, 10))
	v.Set("property", input.Property)
	v.Set("start_time", input.StartTime.String())
	v.Set("force", strconv.FormatBool(input.Force))

	token, err := s.getToken()
	if err != nil {
		return err
	}
	api := &request.API{
		Host:        qiniu.StringValue(s.Config.TradeHost),
		ContentType: "application/x-www-form-urlencoded",
		Method:      "POST",
		Path:        "/product/order/accomplish",
		ServiceName: ServiceName,
		APIName:     "ProductOrderAccomplish",
	}

	req := s.NewRequest(api, strings.NewReader(v.Encode()), out)
	token.SetAuthHeader(req.HTTPRequest)
	return req.Send()
}

// Products 从BO系统获取商品列表
func (s *Service) Products(sellerID int, out interface{}) error {
	if s.Config.TradeHost == nil {
		return errors.New("trade host cannot be empty")
	}
	token, err := s.getToken()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Set("seller_id", strconv.FormatInt(int64(sellerID), 10))
	v.Set("status", "2")
	v.Set("page_size", "20")

	api := &request.API{
		Method:      "GET",
		Host:        qiniu.StringValue(s.Config.TradeHost),
		Path:        "/seller/product?" + v.Encode(),
		ServiceName: ServiceName,
		APIName:     "Products",
	}
	req := s.NewRequest(api, nil, out)
	token.SetAuthHeader(req.HTTPRequest)

	return req.Send()
}

// CreateOrder 向BO系统下订单
//
// 通过 admin oauth 获取token，调用 bo接口，创建订单
func (s *Service) CreateOrder(input *models.ReqOrderNew, out interface{}) error {
	token, err := s.getToken()
	if err != nil {
		return err
	}
	api := &request.API{
		Path:        "/api/order/new",
		Method:      "POST",
		ContentType: "application/json",
		Host:        qiniu.StringValue(s.Config.GaeaHost),
		ServiceName: ServiceName,
		APIName:     "CreateOrder",
	}
	req := s.NewRequest(api, input, out)
	token.SetAuthHeader(req.HTTPRequest)

	return req.Send()
}

func (s *Service) getToken() (*models.Token, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.token.Valid() {
		return s.token, nil
	}
	if s.token != nil {
		token, err := s.RefreshToken(s.token.RefreshToken)
		if err != nil {
			s.token = token
			return token, nil
		}
	}
	token, err := s.PasswordCredentialsToken(nil, qiniu.StringValue(s.Config.User), qiniu.StringValue(s.Config.Pass))
	if err == nil {
		s.token = token
	}
	return token, err
}

// PasswordCredentialsToken 使用账号和密码来获取access token
func (s *Service) PasswordCredentialsToken(ctx context.Context, username, password string) (*models.Token, error) {

	if s.Config.User != nil && s.Config.Pass != nil {
		username = qiniu.StringValue(s.Config.User)
		password = qiniu.StringValue(s.Config.Pass)
	}
	if username == "" || password == "" || s.Config.AccHost == nil {
		return nil, errors.New("username, password, acc_host cannot be empty")
	}
	v := url.Values{}
	v.Set("grant_type", "password")
	v.Set("username", username)
	v.Set("password", password)

	api := &request.API{
		Method:      "POST",
		Path:        "/oauth2/token?" + v.Encode(),
		Host:        qiniu.StringValue(s.Config.AccHost),
		ServiceName: ServiceName,
		APIName:     "PasswordCredentialsToken",
	}
	var token models.Token
	req := s.NewRequest(api, nil, &token)
	if ctx != nil {
		req.SetContext(ctx)
	}
	if err := req.Send(); err != nil {
		return nil, err
	}
	token.SetExpiry(token.ExpiresIN)
	return &token, nil
}

// RefreshToken 刷新当前的Token
// 调用该方法之前要确保当前的token指针不为空
func (s *Service) RefreshToken(refreshToken string) (*models.Token, error) {
	v := url.Values{}
	v.Set("grant_type", "refresh_token")
	v.Set("refresh_token", refreshToken)

	api := &request.API{
		Host:   qiniu.StringValue(s.Config.AccHost),
		Method: "POST",
		Path:   "/oauth2/token?" + v.Encode(),
		//ContentType: "application/x-www-form-urlencoded",
		ServiceName: ServiceName,
		APIName:     "RefreshToken",
	}
	var token models.Token
	req := s.NewRequest(api, nil, &token)

	if err := req.Send(); err != nil {
		return nil, err
	}
	token.SetExpiry(token.ExpiresIN)
	return &token, nil
}

// LoginRequired 检查token有没有过期， 如果没有过期返回用户信息
func (s *Service) LoginRequired(clientID, loginToken string) (*models.SSOUserInfo, error) {
	params := url.Values{}
	params.Set("client_id", clientID)

	if len(loginToken) != 0 {
		params.Set("login_token", loginToken)
	}
	api := &request.API{
		Path:        "/loginrequired?" + params.Encode(),
		Host:        qiniu.StringValue(s.Config.SSOHost),
		Method:      "GET",
		ServiceName: ServiceName,
		APIName:     "SSOLoginRequired",
	}
	var info models.SSOUserInfo

	req := s.NewRequest(api, nil, &info)
	if err := req.Send(); err != nil {
		return nil, err
	}
	return &info, nil
}

// DeveloperInfoUID 通过UID获取开发者信息
func (s *Service) DeveloperInfoUID(uid uint32) (info *account.DeveloperInfo, err error) {
	req, err := s.DeveloperRequestUID(uid, info)
	if err != nil {
		return nil, err
	}
	err = req.Send()
	return
}

// DeveloperInfoEmail 通过email获取开发者信息
func (s *Service) DeveloperInfoEmail(email string) (info *account.DeveloperInfo, err error) {
	req, err := s.DeveloperRequestEmail(email, info)
	if err != nil {
		return nil, err
	}
	err = req.Send()
	return
}

// DeveloperRequestUID 生成一个请求开发这信息的请求
func (s *Service) DeveloperRequestUID(uid uint32, out interface{}) (*request.Request, error) {
	api := &request.API{
		Host:        qiniu.StringValue(s.Config.APIHost),
		Path:        fmt.Sprintf("/api/developer/%d/overview", uid),
		Method:      "GET",
		ServiceName: ServiceName,
		APIName:     "DeveloperInfo",
	}
	return s.newRequestWithAdminToken(api, nil, out)
}

func (s *Service) newRequestWithAdminToken(op *request.API, params interface{}, data interface{}) (*request.Request, error) {
	req := s.NewRequest(op, params, data)
	token, err := s.getToken()
	if err != nil {
		return nil, err
	}
	token.SetAuthHeader(req.HTTPRequest)

	return req, nil
}

// DeveloperRequestEmail 生成一个请求开发这信息的请求
func (s *Service) DeveloperRequestEmail(email string, out interface{}) (*request.Request, error) {
	api := &request.API{
		Host:        qiniu.StringValue(s.Config.APIHost),
		Path:        fmt.Sprintf("/api/developer/%s/overview", email),
		Method:      "GET",
		ServiceName: ServiceName,
		APIName:     "DeveloperInfo",
	}
	return s.newRequestWithAdminToken(api, nil, out)
}

type LicenseVersion string

const (
	LicenseNoFound LicenseVersion = ""
	LicenseActived LicenseVersion = "0.1"
)

type Gender int

const (
	GenderMale   Gender = 0
	GenderFemale Gender = 1
)

type ImCategory int

const (
	QQ ImCategory = iota
	MSN
	GTalk
	Skype
	Other
)

type InternalCategory int

const (
	_InternalCategoryMin InternalCategory = -1
	NormalUser           InternalCategory = 0
	InternalUser         InternalCategory = 1
	TestUser             InternalCategory = 2
	_InternalCategoryMax InternalCategory = 3
)

type TotpStatus int

const (
	TotpStatusDisabled TotpStatus = iota
	TotpStatusEnabled
)

type TotpType int

type resp struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

const (
	TotpTypeAuthenticator TotpType = iota
	TotpTypeMobile
)

type Developer struct {
	Uid             uint32         `json:"uid"`
	Email           string         `json:"email"`
	FullName        string         `json:"fullName"`
	Gender          Gender         `json:"gender"`
	PhoneNumber     string         `json:"phoneNumber"`
	ImNumber        string         `json:"imNumber"`
	ImCategory      ImCategory     `json:"imCategory"`
	WebSite         string         `json:"webSite"`
	CompanyName     string         `json:"companyName"`
	ContractAddress string         `json:"contractAddress"`
	MobileBinded    bool           `json:"mobileBinded"`
	LicenseVersion  LicenseVersion `json:"licenseVersion"`

	Tags []string `json:"tags"`

	RegisterIp       string `json:"registerIp"`
	RegisterState    string `json:"registerState"`
	RegisterRegion   string `json:"registerRegion"`
	RegisterCity     string `json:"registerCity"`
	LocationProvince string `json:"locationProvince"`
	LocationCity     string `json:"locationCity"`

	Referrer           string           `json:"referrer"`
	InviterUid         uint32           `json:"inviterUid"`
	InviteBySales      bool             `json:"inviteBySales"`
	IsActivated        bool             `json:"isActivated"`
	InternalCategory   InternalCategory `json:"internalCategory"`
	InternalDepartment int              `json:"internalDepartment"`

	CreateAt               int64     `json:"createAt"`
	CreatedAtTime          time.Time `json:"createdAtTime"`
	UpdateAt               time.Time `json:"updateAt"`
	UpgradeStdAt           time.Time `json:"upgradeStdAt"`
	UpgradeVipAt           time.Time `json:"upgradeVipAt"`
	LastPasswordModifyTime time.Time `json:"lastPasswordModifyTime"`
	LastEmailModifyTime    time.Time `json:"lastEmailModifyTime"`

	TotpStatus TotpStatus `json:"totpStatus"`
	TotpType   TotpType   `json:"totpType"`

	EmailHistory []string `json:"emailHistory"`

	SfIsEnterprise  bool   `json:"sfIsEnterprise"`
	SfLeadsId       string `json:"sfLeadsId"`
	SfAccountId     string `json:"sfAccountId"`
	SfOpportunityId string `json:"sfOpportunityId"`
	SfSalesId       string `json:"sfSalesId"`
	SfInviteCode    string `json:"sfInviteCode"`
	SfUserId        string `json:"sfUserId"`
}

// GetDeveloper 获取开发者的信息
func (s *Service) GetDeveloper(uid uint32) (*Developer, error) {
	req, developer, err := s.GetDeveloperRequest(uid)
	if err != nil {
		return nil, err
	}
	if err := req.Send(); err != nil {
		return nil, err
	}
	return developer, nil
}

func (s *Service) GetDeveloperRequest(uid uint32) (req *request.Request, developer *Developer, err error) {
	op := &request.API{
		Method:      "GET",
		Path:        fmt.Sprintf("/api/developer?uid=%d", uid),
		Host:        qiniu.StringValue(s.Config.GaeaHost),
		ServiceName: ServiceName,
		APIName:     "GetDeveloper",
	}

	var resp = struct {
		resp
		Data Developer `json:"data"`
	}{}
	req, err = s.newRequestWithAdminToken(op, nil, &resp)
	developer = &resp.Data

	return
}

type User struct {
	Id          string
	Name        string
	CnName      string
	Email       string
	Mobile      string
	QQ          string
	GitHub      string
	WikiDot     string
	Slack       string
	Extension   string
	SfSalesId   string `json:"sf_sales_id" bson:"sf_sales_id"`
	Status      byte
	Delete      byte
	IsTotpOpen  bool
	Create_time time.Time
	Update_time time.Time
}

func (s *Service) GetUserRequest(salesID string) (req *request.Request, user *User, err error) {
	op := &request.API{
		Method:      "GET",
		Path:        "/api/user?salesId=" + salesID,
		Host:        qiniu.StringValue(s.Config.GaeaHost),
		ServiceName: ServiceName,
		APIName:     "GetUser",
	}
	user = &User{}
	var resp = struct {
		resp
		Data User `json:"data"`
	}{}
	req, err = s.newRequestWithAdminToken(op, nil, &resp)
	user = &resp.Data

	return
}

// GetUser 获取用户信息
func (s *Service) GetUser(salesID string) (*User, error) {
	req, user, err := s.GetUserRequest(salesID)
	if err != nil {
		return nil, err
	}
	if err := req.Send(); err != nil {
		return nil, err
	}
	return user, nil
}

var (
	// ErrEmailReceiverEmpty 邮件接收者没有设置
	ErrEmailReceiverEmpty = errors.New("email receiver empty")

	// ErrEmailSubjectEmpty 邮件主题没有设置
	ErrEmailSubjectEmpty = errors.New("email subject empty")
)

// Email 是要发送的一封邮件信息， 包含了主题， 接收方， 发送的信息
// 发送邮件会直接调用七牛的morse服务， 不会进行邮件格式验证， 这些交给该服务验证
type Email struct {
	// 主题
	Subject string `json:"subject"`

	// 接收方
	To []string `json:"to"`

	// 要抄送的人
	CC []string `json:"cc,omitempty"`

	// 密送方地址
	BCC []string `json:"bcc,omitempty"`

	// 要发送的数据
	Message string `json:"content"`

	// 用户UID
	UID uint32 `json:"uid"`
}

// Validate 验证邮件必填字段是否设置了
// 如果有必填字段没有设置，返回相应的错误信息
func (e *Email) Validate() error {
	e.Subject = strings.TrimSpace(e.Subject)

	if len(e.Subject) == 0 {
		return ErrEmailSubjectEmpty
	}
	if len(e.To) == 0 {
		return ErrEmailReceiverEmpty
	}
	return nil
}

// SendEmail 发送一封邮件， 如果发生了错误返回错误信息, 否则错误为nil
func SendEmail(sender Emailer, e *Email) error {
	return sender.SendEmail(e)
}

// Emailer 用来发送邮件
type Emailer interface {
	SendEmail(*Email) error
}

// SendEmail 使用七牛morse服务发送邮件
func (s *Service) SendEmail(e *Email) error {
	api := &request.API{
		APIName:     "SendEmail",
		ServiceName: ServiceName,
		Host:        qiniu.StringValue(s.Config.MorseHost),
		Path:        "/api/notification/send/mail",
		Method:      "POST",
		ContentType: "application/json",
	}
	req := s.NewRequest(api, e, nil)
	req.HTTPRequest.Header.Add("Client-Id", qiniu.StringValue(s.Config.EmailClientID))
	return req.Send()
}
