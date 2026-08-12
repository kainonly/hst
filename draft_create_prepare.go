package hst

import (
	"context"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
)

type CreatePrepareDto struct {
	ReqTimestamp          string        `json:"reqTimestamp"`          // 业务请求时间戳，须与外层信封 timestamp 一致（防重放）
	ChannelNo             string        `json:"channelNo"`             // 渠道号
	ProductCode           []string      `json:"productCode"`           // 已选产品编码列表（WiCoinProductCode 枚举）
	LegalName             string        `json:"legalName"`             // 法定名称（营业执照上的公司名称）
	ShortName             string        `json:"shortName"`             // 商户简称 / 常用名
	MerchantBaseType      string        `json:"merchantBaseType"`      // 商户类型：01(自然人)/02(个体工商户)/03(企业商户)
	SubRoleType           string        `json:"subRoleType"`           // 商户角色
	DealType              string        `json:"dealType"`              // 商户经营类型：01(实体特约)/02(网络特约)/03(实体兼网络特约)
	Mcc                   string        `json:"mcc"`                   // 经营类目（MCC 码）
	ContactMobile         string        `json:"contactMobile"`         // 联系人手机号
	ContactName           string        `json:"contactName"`           // 联系人姓名
	Email                 string        `json:"email"`                 // 联系人邮箱
	PrincipalMobile       string        `json:"principalMobile"`       // 负责人手机号（即法人手机号）
	PrincipalCertType     string        `json:"principalCertType"`     // 负责人证件类型
	PrincipalCertNo       string        `json:"principalCertNo"`       // 负责人证件号码
	PrincipalPerson       string        `json:"principalPerson"`       // 负责人名称或企业法人代表姓名
	PrincipalCertVld      string        `json:"principalCertVld"`      // 负责人证件有效期，格式 yyyy-MM-dd HH:mm:ss
	MangerLogonId         string        `json:"mangerLogonId"`         // 管理员支付宝登录号
	Province              string        `json:"province"`              // 省份（国标省市区号）
	City                  string        `json:"city"`                  // 城市（国标省市区号）
	District              string        `json:"district"`              // 区（国标省市区号）
	Address               string        `json:"address"`               // 详细地址
	ServicePhoneNo        string        `json:"servicePhoneNo"`        // 商户客服电话
	TaxNum                string        `json:"taxNum"`                // 税务登记证号码
	BussAuthVld           string        `json:"bussAuthVld"`           // 营业执照有效期，格式 yyyy-MM-dd HH:mm:ss
	ShareholderName       string        `json:"shareholderName"`       // 控股股东或实际控制人姓名
	ShareholderCertType   string        `json:"shareholderCertType"`   // 控股股东或实际控制人证件类型
	ShareholderCertNo     string        `json:"shareholderCertNo"`     // 控股股东或实际控制人证件号码
	ShareholderCertVld    string        `json:"shareholderCertVld"`    // 控股股东或实际控制人证件有效期，格式 yyyy-MM-dd HH:mm:ss
	PersonSex             string        `json:"personSex"`             // 性别（自然人商户）：M(男性)/F(女性)
	PersonProfession      string        `json:"personProfession"`      // 职业（自然人商户）
	PersonCertVld         string        `json:"personCertVld"`         // 身份证件有效期限（自然人商户），格式 yyyy-MM-dd HH:mm:ss
	BussAuthType          string        `json:"bussAuthType"`          // 营业执照证件类型
	BussAuthNo            string        `json:"bussAuthNo"`            // 证件号码（营业执照号或统一社会信用代码）
	Remark                string        `json:"remark"`                // 备注
	PartnerId             string        `json:"partnerId"`             // 渠道商 ID
	MerchantType          string        `json:"merchantType"`          // 商户业务类型，如 SCHOOL、GROUP_MEAL、HR 等
	AgreementNo           string        `json:"agreementNo"`           // 支付宝签约记录编号（安全发）
	AlipayPid             string        `json:"alipayPid"`             // 支付宝商户 ID
	AlipayAccount         string        `json:"alipayAccount"`         // 支付宝收款账号
	LogicGroupId          string        `json:"logicGroupId"`          // 支付宝学校、机构用户库 ID
	WxSubMchId            string        `json:"wxSubMchId"`            // 微信商户号
	WxSubMchAccount       string        `json:"wxSubMchAccount"`       // 微信收款账号
	SettlementAccountType string        `json:"settlementAccountType"` // 结算类型：01(银行卡)/02(支付宝)/03(支付宝虚拟账户)
	BankCardNo            string        `json:"bankCardNo"`            // 银行卡号
	BankCertName          string        `json:"bankCertName"`          // 银行账户户名
	AccountType           string        `json:"accountType"`           // 账户类型：01(对私账户)/02(对公账户)
	ContactLine           string        `json:"contactLine"`           // 联行号
	BranchName            string        `json:"branchName"`            // 开户支行名称
	BranchProvince        string        `json:"branchProvince"`        // 开户支行所在省
	BranchCity            string        `json:"branchCity"`            // 开户支行所在市
	CertType              string        `json:"certType"`              // 持卡人证件类型：01(身份证)
	CertNo                string        `json:"certNo"`                // 持卡人证件号码
	CardHolderAddress     string        `json:"cardHolderAddress"`     // 持卡人地址
	LogonId               string        `json:"logonId"`               // 支付宝登陆账号
	UserId                string        `json:"userId"`                // 支付宝用户 ID
	FileManifest          *FileManifest `json:"fileManifest"`          // 文件哈希清单：字段名 -> 按上传顺序排列的 SM3 哈希列表（64 位 Hex）
}

func (x *CreatePrepareDto) GetTs() string {
	return x.ReqTimestamp
}

// NewCreatePrepareDto 创建创建草稿并申请上传凭证请求体。
// 必填字段作为参数传入（reqTimestamp / channelNo 由 CreatePrepare 方法自动填充）。
func NewCreatePrepareDto(
	productCode []string,
	legalName string,
	shortName string,
	merchantBaseType string,
	subRoleType string,
	dealType string,
	mcc string,
	contactMobile string,
	contactName string,
	email string,
	principalMobile string,
	principalCertType string,
	principalCertNo string,
	principalPerson string,
	principalCertVld string,
	province string,
	city string,
	district string,
	address string,
	servicePhoneNo string,
	personSex string,
	personProfession string,
	settlementAccountType string,
	bankCardNo string,
	bankCertName string,
	accountType string,
	contactLine string,
	branchName string,
	branchProvince string,
	branchCity string,
	certType string,
	certNo string,
	cardHolderAddress string,
	logonId string,
	userId string,
	fileManifest *FileManifest,
) *CreatePrepareDto {
	return &CreatePrepareDto{
		ProductCode:           productCode,
		LegalName:             legalName,
		ShortName:             shortName,
		MerchantBaseType:      merchantBaseType,
		SubRoleType:           subRoleType,
		DealType:              dealType,
		Mcc:                   mcc,
		ContactMobile:         contactMobile,
		ContactName:           contactName,
		Email:                 email,
		PrincipalMobile:       principalMobile,
		PrincipalCertType:     principalCertType,
		PrincipalCertNo:       principalCertNo,
		PrincipalPerson:       principalPerson,
		PrincipalCertVld:      principalCertVld,
		Province:              province,
		City:                  city,
		District:              district,
		Address:               address,
		ServicePhoneNo:        servicePhoneNo,
		PersonSex:             personSex,
		PersonProfession:      personProfession,
		SettlementAccountType: settlementAccountType,
		BankCardNo:            bankCardNo,
		BankCertName:          bankCertName,
		AccountType:           accountType,
		ContactLine:           contactLine,
		BranchName:            branchName,
		BranchProvince:        branchProvince,
		BranchCity:            branchCity,
		CertType:              certType,
		CertNo:                certNo,
		CardHolderAddress:     cardHolderAddress,
		LogonId:               logonId,
		UserId:                userId,
		FileManifest:          fileManifest,
	}
}

// SetMangerLogonId 设置管理员支付宝登录号。
func (x *CreatePrepareDto) SetMangerLogonId(i string) *CreatePrepareDto {
	x.MangerLogonId = i
	return x
}

// SetTaxNum 设置税务登记证号码。
func (x *CreatePrepareDto) SetTaxNum(i string) *CreatePrepareDto {
	x.TaxNum = i
	return x
}

// SetBussAuthVld 设置营业执照有效期，格式 yyyy-MM-dd HH:mm:ss。
func (x *CreatePrepareDto) SetBussAuthVld(i string) *CreatePrepareDto {
	x.BussAuthVld = i
	return x
}

// SetShareholderName 设置控股股东或实际控制人姓名。
func (x *CreatePrepareDto) SetShareholderName(i string) *CreatePrepareDto {
	x.ShareholderName = i
	return x
}

// SetShareholderCertType 设置控股股东或实际控制人证件类型。
func (x *CreatePrepareDto) SetShareholderCertType(i string) *CreatePrepareDto {
	x.ShareholderCertType = i
	return x
}

// SetShareholderCertNo 设置控股股东或实际控制人证件号码。
func (x *CreatePrepareDto) SetShareholderCertNo(i string) *CreatePrepareDto {
	x.ShareholderCertNo = i
	return x
}

// SetShareholderCertVld 设置控股股东或实际控制人证件有效期，格式 yyyy-MM-dd HH:mm:ss。
func (x *CreatePrepareDto) SetShareholderCertVld(i string) *CreatePrepareDto {
	x.ShareholderCertVld = i
	return x
}

// SetPersonCertVld 设置身份证件有效期限（自然人商户），格式 yyyy-MM-dd HH:mm:ss。
func (x *CreatePrepareDto) SetPersonCertVld(i string) *CreatePrepareDto {
	x.PersonCertVld = i
	return x
}

// SetBussAuthType 设置营业执照证件类型。
func (x *CreatePrepareDto) SetBussAuthType(i string) *CreatePrepareDto {
	x.BussAuthType = i
	return x
}

// SetBussAuthNo 设置证件号码（营业执照号或统一社会信用代码）。
func (x *CreatePrepareDto) SetBussAuthNo(i string) *CreatePrepareDto {
	x.BussAuthNo = i
	return x
}

// SetRemark 设置备注。
func (x *CreatePrepareDto) SetRemark(i string) *CreatePrepareDto {
	x.Remark = i
	return x
}

// SetPartnerId 设置渠道商 ID。
func (x *CreatePrepareDto) SetPartnerId(i string) *CreatePrepareDto {
	x.PartnerId = i
	return x
}

// SetMerchantType 设置商户业务类型，如 SCHOOL、GROUP_MEAL、HR 等。
func (x *CreatePrepareDto) SetMerchantType(i string) *CreatePrepareDto {
	x.MerchantType = i
	return x
}

// SetAgreementNo 设置支付宝签约记录编号（安全发）。
func (x *CreatePrepareDto) SetAgreementNo(i string) *CreatePrepareDto {
	x.AgreementNo = i
	return x
}

// SetAlipayPid 设置支付宝商户 ID。
func (x *CreatePrepareDto) SetAlipayPid(i string) *CreatePrepareDto {
	x.AlipayPid = i
	return x
}

// SetAlipayAccount 设置支付宝收款账号。
func (x *CreatePrepareDto) SetAlipayAccount(i string) *CreatePrepareDto {
	x.AlipayAccount = i
	return x
}

// SetLogicGroupId 设置支付宝学校、机构用户库 ID。
func (x *CreatePrepareDto) SetLogicGroupId(i string) *CreatePrepareDto {
	x.LogicGroupId = i
	return x
}

// SetWxSubMchId 设置微信商户号。
func (x *CreatePrepareDto) SetWxSubMchId(i string) *CreatePrepareDto {
	x.WxSubMchId = i
	return x
}

// SetWxSubMchAccount 设置微信收款账号。
func (x *CreatePrepareDto) SetWxSubMchAccount(i string) *CreatePrepareDto {
	x.WxSubMchAccount = i
	return x
}

type CreatePrepareBizData struct {
	UploadToken   string `json:"uploadToken"`
	ExpireSeconds int64  `json:"expireSeconds"`
}

func (x *Hst) CreatePrepare(ctx context.Context, dto *CreatePrepareDto) (result *SignObjectRespResult[*CreatePrepareBizData], err error) {
	dto.ReqTimestamp = strconv.FormatInt(time.Now().UnixMilli(), 10)
	dto.ChannelNo = x.Option.ChannelId

	var signObjectReq *SignObjectReq
	if signObjectReq, err = x.NewSignObjectReq(dto); err != nil {
		return
	}

	var b string
	if b, err = x.Request(ctx,
		"/merchant_info_draft/create/prepare", signObjectReq); err != nil {
		return
	}
	if err = sonic.UnmarshalString(b, &result); err != nil {
		return
	}
	if !result.BizSuccess {
		err = bizError(result.BizCode, result.BizMsg)
	}
	return
}
