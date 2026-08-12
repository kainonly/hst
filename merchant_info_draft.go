package hst

import (
	"context"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
)

type CreatePrepareBody struct {
	ReqTimestamp    string `json:"reqTimestamp"`    // 业务请求时间戳，须与外层信封 timestamp 一致（防重放）
	ChannelNo       string `json:"channelNo"`       // 渠道号
	EnvelopeNo      string `json:"envelopeNo"`      // 商品订单号（用于查询预下单结果或发起支付）
	ProductCode     string `json:"productCode"`     // 产品编码（见附录《下单产品码字典》），例如 WECHAT_MINI_PROGRAM、CASHIER
	BuyerRealName   string `json:"buyerRealName"`   // 买家姓名（仅转账时必填）
	BuyerPhone      string `json:"buyerPhone"`      // 买家手机号
	ReceiverName    string `json:"receiverName"`    // 收货人姓名
	ReceiverPhone   string `json:"receiverPhone"`   // 收货人手机号
	NotifyUrl       string `json:"notifyUrl"`       // 商户回调地址
	VerifyCode      string `json:"verifyCode"`      // 商户内部的付款码值
	RedirectUrl     string `json:"redirectUrl"`     // 支付成功回跳地址
	Period          string `json:"period"`          // 收款周期
	PersonOrgType   string `json:"personOrgType"`   // 主体类型（PERSON 个人/ORG 机构）
	PersonCertNo    string `json:"personCertNo"`    // 法人身份证号
	PersonCertName  string `json:"personCertName"`  // 法人姓名
	PersonCardNo    string `json:"personCardNo"`    // 个人结算卡号
	PersonCertFront string `json:"personCertFront"` // 人卡正面照 url
	PersonCertBack  string `json:"personCertBack"`  // 人卡反面照 url
	PersonAgrNo     string `json:"personAgrNo"`     // 个人银行卡正面协议号
	ResveFieldOne   string `json:"resveFieldOne"`   // 预留字段一
	AmountTotal     string `json:"amountTotal"`     // 金额（单位：分）
	Discount        string `json:"discount"`        // 优惠金额
	Score           string `json:"score"`           // 积分
	Freight         string `json:"freight"`         // 运费
	Charge          string `json:"charge"`          // 附加费
	OriginalAmount  string `json:"originalAmount"`  // 订单原价
	BuyerInfo       string `json:"buyerInfo"`       // 买家扩展信息
	InvoiceTitle    string `json:"invoiceTitle"`    // 发票抬头
	InvoiceType     string `json:"invoiceType"`     // 发票类型（COMPANY/INDIVIDUAL）
	InvoiceCertNo   string `json:"invoiceCertNo"`   // 发票纳税人识别号
	InvoiceBank     string `json:"invoiceBank"`     // 发票开户行
	BankAcct        string `json:"bankAcct"`        // 发票账号
	InvoiceAddr     string `json:"invoiceAddr"`     // 发票地址
	InvoicePhone    string `json:"invoicePhone"`    // 发票电话
	InvoiceEmail    string `json:"invoiceEmail"`    // 发票邮箱
	InvoiceRemark   string `json:"invoiceRemark"`   // 发票备注
	BankAcctType    string `json:"bankAcctType"`    // 银行开户类型
	ExpirationTime  string `json:"expirationTime"`  // 有效时间
	StartTime       string `json:"startTime"`       // 开始时间
	PeriodStart     string `json:"periodStart"`     // 周期开始时间（格式：YYYY-MM-DD）
	PeriodEnd       string `json:"periodEnd"`       // 周期结束时间
	LoginToken      string `json:"loginToken"`      // 登录令牌
	CallerNumber    string `json:"callerNumber"`    // 主叫号码
	CalledNumber    string `json:"calledNumber"`    // 被叫号码
	SubAppid        string `json:"subAppid"`        // 微信子应用 ID
	Openid          string `json:"openid"`          // 用户 openid
	AuthCode        string `json:"authCode"`        // 支付授权码
	Appid           string `json:"appid"`           // 应用 ID
	Userid          string `json:"userid"`          // 用户 ID
	SubOpenid       string `json:"subOpenid"`       // 用户子 openid
	LimitPay        string `json:"limitPay"`        // 支付限制
	IdCardNo        string `json:"idCardNo"`        // 身份证号
	Accno           string `json:"accno"`           // 银行卡号
	Address         string `json:"address"`         // 收货地址
	VipId           string `json:"vipId"`           // 会员 ID
	ReceiverType    string `json:"receiverType"`    // 收货人类型
	IcData          string `json:"icData"`          // 芯片数据（IC 卡）
	Extend          string `json:"extend"`          // 拓展字段
}

//func (x *Hst) CreatePrepare(ctx context.Context, body *CreatePrepareBody) (resp *CreatePrepareResp, err error) {
//	resp = &CreatePrepareResp{}
//	_, err = x.Client.NewRequest().SetContext(ctx).SetBody(body).SetResult(resp).Post("/createPrepare")
//	return
//}

type SettlementStatusBody struct {
	ReqTimestamp string `json:"reqTimestamp"` // 业务请求时间戳，须与外层信封 timestamp 一致（防重放）
	ChannelNo    string `json:"channelNo"`    // 渠道号
	DraftId      string `json:"draftId"`      // 草稿 ID
}

func (x *SettlementStatusBody) GetTs() string {
	return x.ReqTimestamp
}

func NewSettlementStatusBody(draftId string) *SettlementStatusBody {
	return &SettlementStatusBody{DraftId: draftId}
}

type SettlementStatusBizData struct {
	DraftId              string `json:"draftId"`              // 草稿 ID
	MerchantId           string `json:"merchantId"`           // 商户号（若草稿已生成商户则有值，否则为空）
	MerchantCreated      bool   `json:"merchantCreated"`      // 是否已生成商户号
	SettlementStatus     string `json:"settlementStatus"`     // 入驻状态码：0-审核中，1-成功，2-失败，3-待激活，4-激活中
	SettlementStatusDesc string `json:"settlementStatusDesc"` // 入驻状态描述
	MybankOrderNo        string `json:"mybankOrderNo"`        // 网商银行申请单号
	MybankMerchantId     string `json:"mybankMerchantId"`     // 网商银行商户编号
	FailReason           string `json:"failReason"`           // 开户失败原因
	ActivateUrl          string `json:"activateUrl"`          // 激活链接（settlementStatus = 3 待激活时使用）
}

func (x *Hst) SettlementStatus(ctx context.Context, body *SettlementStatusBody) (bizData *SignObjectRespBody[SettlementStatusBizData], err error) {
	body.ReqTimestamp = strconv.FormatInt(time.Now().UnixMilli(), 10)
	body.ChannelNo = x.Option.ChannelId

	var signObjectReq *SignObjectReq
	if signObjectReq, err = x.NewSignObjectReq(body); err != nil {
		return
	}

	var b string
	if b, err = x.Request(ctx,
		"/merchant_info_draft/settlement_status", signObjectReq); err != nil {
		return
	}

	if err = sonic.UnmarshalString(b, &bizData); err != nil {
		return
	}
	return
}
