package models

import (
	"time"

	"github.com/QN-zhangzhuo/go-sdk/qiniu"
)

// ReqProductOrderNew 请求下单接口要用的信息
type ReqProductOrderNew struct {
	ProductID    int64    `json:"product_id"`
	Duration     uint     `json:"duration"`
	TimeDuration uint64   `json:"time_duration,omitempty"`
	Quantity     uint     `json:"quantity"`
	Property     *string  `json:"property,omitempty"`
	Fee          *float64 `json:"fee,omitempty"` // 不是必须传递的
}

// ReqProductOrderAccomplish 是发货完成需要回调给BO系统的信息
type ReqProductOrderAccomplish struct {
	ID        int64              `json:"id"`
	Property  string             `json:"property"`
	StartTime *qiniu.RFC3339Time `json:"start_time"`
	Force     bool               `json:"force"`
}

// ReqOrderNew 下单信息
type ReqOrderNew struct {
	BuyerID uint32 `json:"uid"`

	// 备注信息
	Memo string `json:"memo,omitempty"`

	Orders []ReqProductOrderNew `json:"orders"`
}

// OrderNewResponse 是BO GAEA API下单接口返回的数据结构
type OrderNewResponse struct {
	Data OrderHash `json:"data,omitempty"`

	Code int `json:"code"`

	Message string `json:"message,omitempty"`
}

// OrderHash 订单号
type OrderHash struct {
	Order string `json:"order_hash"`
}

// Order 时BO订单系统生成的信息
type Order struct {
	ID            int64           `json:"id"`
	OrderHash     string          `json:"order"`
	SellerID      int64           `json:"seller_id"`
	BuyerID       uint32          `json:"buyer_id"`
	Fee           float64         `json:"fee"`
	CFee          float64         `json:"c_fee"`
	Memo          string          `json:"memo"`
	TxnType       TxnType         `json:"txn_type"`
	PayTime       time.Time       `json:"pay_time"`
	PayActionTime time.Time       `json:"pay_action_time"`
	UpdateTime    time.Time       `json:"update_time"`
	CreateTime    time.Time       `json:"create_time"`
	ExpiredTime   time.Time       `json:"expired_time"`
	Status        OrderStatus     `json:"status"`
	Version       int             `json:"version"`
	Products      *[]Product      `json:"products,omitempty"`
	ProductOrders *[]ProductOrder `json:"product_orders,omitempty"`
}

// TxnType 订单支付流水版本
type TxnType int

const (
	// TxnOrder order trans
	TxnOrder TxnType = iota
	// TxnProductOrder po trans
	TxnProductOrder
)

// Product 是商品信息
type Product struct {
	ID             int64          `json:"id"`        // 商品 id
	Name           string         `json:"name"`      // 商品名称
	SellerID       int64          `json:"seller_id"` // 商家 id
	Model          string         `json:"model"`
	SPU            string         `json:"spu"`
	Unit           ProductUnit    `json:"unit"`            // 商品基础单位
	OriginalPrice  float64        `json:"original_price"`  // 商品原价
	Price          float64        `json:"price"`           // 商品现价
	ExpiresIn      uint64         `json:"expires_in"`      // 过期时长
	Property       string         `json:"property"`        // 商品属性
	Description    string         `json:"description"`     // 商品描述
	UpdateTime     time.Time      `json:"update_time"`     // 最后更新时间
	CreateTime     time.Time      `json:"create_time"`     // 商品创建时间
	StartTime      time.Time      `json:"start_time"`      // 上线时间
	EndTime        time.Time      `json:"end_time"`        // 下线时间
	Status         ProductStatus  `json:"status"`          // 商品状态
	SettlementMode SettlementMode `json:"settlement_mode"` // 结算方式
	CategoryID     int64          `json:"category_id"`     // 商品分类
	Version        int            `json:"version"`
}

// SettlementMode 租用方式
type SettlementMode int

const (
	// SettlementSaleOnce  一次性售出
	SettlementSaleOnce SettlementMode = iota + 1
	// SettlementRentOnTime 按时租用
	SettlementRentOnTime
	// SettlementRentOnVolume 按量租用
	SettlementRentOnVolume
	// SettlementProject 项目制
	SettlementProject
)

// ProductOrder 是一个订单中的子订单信息， 一个子订单对应一个商品购买记录
type ProductOrder struct {
	ID              int64              `json:"id"`
	ProductID       int64              `json:"product_id"`
	SellerID        int64              `json:"seller_id"`
	BuyerID         uint32             `json:"buyer_id"`
	OrderID         int64              `json:"order_id"`
	OrderHash       string             `json:"order_hash"`
	OrderType       OrderType          `json:"order_type"`
	ProductOrderID  int64              `json:"product_order_id"`
	ProductVersion  int                `json:"product_version"`
	ProductName     string             `json:"product_name"`
	ProductProperty string             `json:"product_property"`
	Property        string             `json:"property"`
	Duration        uint               `json:"duration"`
	TimeDuration    time.Duration      `json:"time_duration"`
	Quantity        uint               `json:"quantity"`
	Fee             float64            `json:"fee"`
	CFee            float64            `json:"c_fee"`
	UpdateTime      time.Time          `json:"update_time"`
	CreateTime      time.Time          `json:"create_time"`
	StartTime       time.Time          `json:"start_time"`
	EndTime         time.Time          `json:"end_time"`
	AbortTime       time.Time          `json:"abort_time"`
	RefundAck       time.Time          `json:"refund_ack"`
	ExpiredTime     time.Time          `json:"expired_time"`
	Status          ProductOrderStatus `json:"status"`
	AccomplishTime  time.Time          `json:"accomplish_time"`
	Product         *Product           `json:"product,omitempty"`
	Version         int                `json:"version"`
}

// OrderStatus 表示订单的状态
type OrderStatus int

const (
	// OrderStatusUnPay 未支付
	OrderStatusUnPay OrderStatus = iota + 1

	// OrderStatusPayed 已支付
	OrderStatusPayed

	// OrderStatusCancelled 订单已经取消
	OrderStatusCancelled
	_

	// OrderStatusPostPay 订单后付费
	OrderStatusPostPay

	// OrderStatusPaying 订单支付中
	OrderStatusPaying
)

// Valid 检测订单状态是否合法
func (o OrderStatus) Valid() bool {
	switch o {
	case OrderStatusUnPay, OrderStatusPayed, OrderStatusCancelled, OrderStatusPostPay, OrderStatusPaying:
		return true
	default:
		return false
	}
}

// IsPaying 如果订单在支付中返回true, 否则false
func (o OrderStatus) IsPaying() bool {
	return o == OrderStatusPaying
}

// IsUnPay 如果订单为支付，返回true, 否则false
func (o OrderStatus) IsUnPay() bool {
	return o == OrderStatusUnPay
}

// IsPayed 如果订单已经支付， 返回true, 否则返回false
func (o OrderStatus) IsPayed() bool {
	return o == OrderStatusPayed
}

// IsCancelled 如果订单已经取消返回true, 否则返回false
func (o OrderStatus) IsCancelled() bool {
	return o == OrderStatusCancelled
}

// Humanize 返回方便人类阅读的字符串
func (o OrderStatus) Humanize() string {
	switch o {
	case OrderStatusUnPay:
		return "未支付"
	case OrderStatusPayed:
		return "已支付"
	case OrderStatusCancelled:
		return "作废"
	case OrderStatusPostPay:
		return "后付费"
	case OrderStatusPaying:
		return "支付中"
	default:
		return "未知订单状态"
	}
}

// OrderType 订单类型
type OrderType int

const (
	// OrderTypeBuy 新的购买订单
	OrderTypeBuy OrderType = iota + 1
	_
	// OrderTypeUpgrade 升级订单
	OrderTypeUpgrade
	_

	// OrderTypeRefund 退款订单
	OrderTypeRefund
)

// Valid 检测订单类型是否有效
func (s OrderType) Valid() bool {
	switch s {
	case OrderTypeBuy, OrderTypeUpgrade, OrderTypeRefund:
		return true
	default:
		return false
	}
}

// Humanize 返回容易阅读的字符串
func (s OrderType) Humanize() string {
	switch s {
	case OrderTypeBuy:
		return "新购"
	case OrderTypeUpgrade:
		return "升级"
	case OrderTypeRefund:
		return "退款"
	default:
		return "未知订单类型"
	}
}

// ProductUnit 购买的商品的时间单位
type ProductUnit int

const (
	// Yearly  年
	Yearly ProductUnit = iota + 1

	// Monthly 月
	Monthly

	// Weekly 周
	Weekly

	// Daily 日
	Daily

	// UnLimited 无限制
	UnLimited ProductUnit = 99
)

// Valid 检测时间单位是否有效
func (pu ProductUnit) Valid() bool {
	switch pu {
	case Yearly, Monthly, Weekly, Daily, UnLimited:
		return true
	default:
		return false
	}
}

// String 返回时间单位的字符串表示
func (pu ProductUnit) String() string {
	switch pu {
	case Yearly:
		return "year"
	case Monthly:
		return "month"
	case Weekly:
		return "week"
	case Daily:
		return "day"
	case UnLimited:
		return "unlimited"
	default:
		return "unknown ProductUnit"
	}
}

// Humanize 返回方便人阅读的字符串
func (pu ProductUnit) Humanize() string {
	switch pu {
	case Yearly:
		return "按年"
	case Monthly:
		return "按月"
	case Weekly:
		return "按周"
	case Daily:
		return "按天"
	case UnLimited:
		return "一次性购买"
	default:
		return "未知计费单位"
	}
}

// AddDuration 根据时间单位加上相应的时间
func (pu ProductUnit) AddDuration(baseTime time.Time, duration int) time.Time {
	switch pu {
	case Yearly:
		return baseTime.AddDate(duration, 0, 0)
	case Monthly:
		return baseTime.AddDate(0, duration, 0)
	case Weekly:
		return baseTime.AddDate(0, 0, duration*7)
	case Daily:
		return baseTime.AddDate(0, 0, duration)
	case UnLimited:
		//mysql only support 2038-01-01 00:00:00 for timestamp datetype
		return time.Date(2038, time.January, 1, 0, 0, 0, 0, time.Local)
	default:
		return baseTime
	}
}

// ProductStatus 商品状态
type ProductStatus int

const (
	// ProductStatusNew 商品新建
	ProductStatusNew ProductStatus = iota + 1

	// ProductStatusOnline 商品在线
	ProductStatusOnline

	// ProductStatusDeprecated 商品失效
	ProductStatusDeprecated

	// ProductStatusDeleted 商品被删除
	ProductStatusDeleted
)

// Humanize 返回方便人阅读的字符串
func (p ProductStatus) Humanize() string {
	switch p {
	case ProductStatusNew:
		return "新建"
	case ProductStatusOnline:
		return "在线"
	case ProductStatusDeprecated:
		return "已失效"
	case ProductStatusDeleted:
		return "已删除"
	default:
		return "未知产品状态"
	}
}

// Valid 返回商品状态是否有效
func (p ProductStatus) Valid() bool {
	switch p {
	case ProductStatusNew, ProductStatusOnline, ProductStatusDeprecated, ProductStatusDeleted:
		return true
	default:
		return false
	}
}

// IsNew 如果商品是新建，返回true, 否则false
func (p ProductStatus) IsNew() bool {
	return p == ProductStatusNew
}

// IsOnline 如果商品是在线，返回true, 否则false
func (p ProductStatus) IsOnline() bool {
	return p == ProductStatusOnline
}

// IsDeprecated 如果商品是失效，返回true, 否则false
func (p ProductStatus) IsDeprecated() bool {
	return p == ProductStatusDeprecated
}

// IsDeleted 如果商品被删除，返回true, 否则false
func (p ProductStatus) IsDeleted() bool {
	return p == ProductStatusDeleted
}

// ProductOrderStatus 商品订单状态
type ProductOrderStatus int

const (
	// ProductOrderStatusNew 新建状态
	ProductOrderStatusNew ProductOrderStatus = iota + 1

	// ProductOrderStatusComplete 完成状态
	ProductOrderStatusComplete
)

/*
const (
	ProductCategoryRespack = "respack"
)

// 分类是否是 资源包分类
func (pc *ProductCategory) IsRespack() bool {
	return pc.Code == ProductCategoryRespack
}

// Package status
type PackageStatus int

const (
	// PackageStatusUndefined undefined
	PackageStatusUndefined PackageStatus = iota
	// PackageStatusInvalid invalid
	PackageStatusInvalid
	// PackageStatusValid valid
	PackageStatusValid
)

*/
