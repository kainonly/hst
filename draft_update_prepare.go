package hst

import (
	"context"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
)

type UpdatePrepareDto struct {
	ReqTimestamp          string        `json:"reqTimestamp"`                    // 业务请求时间戳，须与外层信封 timestamp 一致（防重放）
	ChannelNo             string        `json:"channelNo"`                       // 渠道号
	DraftId               string        `json:"draftId"`                         // 草稿 ID（UploadFiles 响应返回，仅 EDITING/FAILED 状态可更新）
	ProductCode           []string      `json:"productCode,omitempty"`           // 已选产品编码列表，WiCoinProductCode 枚举：WICOIN_PAY/WICOIN_ZFT/WICOIN_MYBANK_SPLIT_ACCOUNT/WICOIN_SECURE_PAY
	LegalName             string        `json:"legalName,omitempty"`             // 法定名称（营业执照上的公司名称）
	ShortName             string        `json:"shortName,omitempty"`             // 商户简称 / 常用名
	MerchantBaseType      string        `json:"merchantBaseType,omitempty"`      // 商户类型：01(自然人)/02(个体工商户)/03(企业商户)
	SubRoleType           string        `json:"subRoleType,omitempty"`           // 商户角色：shop_franchisee加盟商门店/shop_direct直营门店/shop_joint联营门店/supplier_brandSupplyChain品牌供应链公司/supplier_external外部供应商/brand_manage品牌管理公司/brand_area品牌区域公司/brand_other其他服务公司
	DealType              string        `json:"dealType,omitempty"`              // 商户经营类型：01(实体特约)/02(网络特约)/03(实体兼网络特约)
	Mcc                   string        `json:"mcc,omitempty"`                   // 经营类目（MCC 码）
	ContactMobile         string        `json:"contactMobile,omitempty"`         // 联系人手机号
	ContactName           string        `json:"contactName,omitempty"`           // 联系人姓名
	Email                 string        `json:"email,omitempty"`                 // 联系人邮箱
	PrincipalMobile       string        `json:"principalMobile,omitempty"`       // 负责人手机号（即法人手机号）
	PrincipalCertType     string        `json:"principalCertType,omitempty"`     // 负责人证件类型：100身份证/102护照/108外国人永久居留身份证/114台湾居民来往大陆通行证/115港澳居民来往内地通行证/116港澳居民居住证/117台湾居民居住证
	PrincipalCertNo       string        `json:"principalCertNo,omitempty"`       // 负责人证件号码
	PrincipalPerson       string        `json:"principalPerson,omitempty"`       // 负责人名称或企业法人代表姓名
	PrincipalCertVld      string        `json:"principalCertVld,omitempty"`      // 负责人证件有效期，格式 yyyy-MM-dd HH:mm:ss
	MangerLogonId         string        `json:"mangerLogonId,omitempty"`         // 管理员支付宝登录号
	Province              string        `json:"province,omitempty"`              // 省份（国标省市区号）
	City                  string        `json:"city,omitempty"`                  // 城市（国标省市区号）
	District              string        `json:"district,omitempty"`              // 区（国标省市区号）
	Address               string        `json:"address,omitempty"`               // 详细地址
	ServicePhoneNo        string        `json:"servicePhoneNo,omitempty"`        // 商户客服电话
	TaxNum                string        `json:"taxNum,omitempty"`                // 税务登记证号码
	BussAuthVld           string        `json:"bussAuthVld,omitempty"`           // 营业执照有效期，格式 yyyy-MM-dd HH:mm:ss
	ShareholderName       string        `json:"shareholderName,omitempty"`       // 控股股东或实际控制人姓名
	ShareholderCertType   string        `json:"shareholderCertType,omitempty"`   // 控股股东或实际控制人证件类型
	ShareholderCertNo     string        `json:"shareholderCertNo,omitempty"`     // 控股股东或实际控制人证件号码
	ShareholderCertVld    string        `json:"shareholderCertVld,omitempty"`    // 控股股东或实际控制人证件有效期，格式 yyyy-MM-dd HH:mm:ss
	PersonSex             string        `json:"personSex,omitempty"`             // 性别（自然人商户）：M(男性)/F(女性)
	PersonProfession      string        `json:"personProfession,omitempty"`      // 职业（自然人商户）
	PersonCertVld         string        `json:"personCertVld,omitempty"`         // 身份证件有效期限（自然人商户），格式 yyyy-MM-dd HH:mm:ss
	BussAuthType          string        `json:"bussAuthType,omitempty"`          // 营业执照证件类型：12营业执照/统一社会信用代码/41政府机关证件/42部队证件/43社会团体证件/44事业单位证件/45民办非企业组织证件/99其它企业类型证件
	BussAuthNo            string        `json:"bussAuthNo,omitempty"`            // 证件号码（营业执照号或统一社会信用代码）
	Remark                string        `json:"remark,omitempty"`                // 备注
	PartnerId             string        `json:"partnerId,omitempty"`             // 渠道商 ID
	MerchantType          string        `json:"merchantType,omitempty"`          // 商户业务类型，如 SCHOOL、GROUP_MEAL、HR 等
	AgreementNo           string        `json:"agreementNo,omitempty"`           // 支付宝签约记录编号（安全发）
	AlipayPid             string        `json:"alipayPid,omitempty"`             // 支付宝商户 ID
	AlipayAccount         string        `json:"alipayAccount,omitempty"`         // 支付宝收款账号
	LogicGroupId          string        `json:"logicGroupId,omitempty"`          // 支付宝学校、机构用户库 ID
	WxSubMchId            string        `json:"wxSubMchId,omitempty"`            // 微信商户号
	WxSubMchAccount       string        `json:"wxSubMchAccount,omitempty"`       // 微信收款账号
	SettlementAccountType string        `json:"settlementAccountType,omitempty"` // 结算类型：01(银行卡)/02(支付宝)/03(支付宝虚拟账户)
	BankCardNo            string        `json:"bankCardNo,omitempty"`            // 银行卡号
	BankCertName          string        `json:"bankCertName,omitempty"`          // 银行账户户名（仅对公账户 02 需要）
	AccountType           string        `json:"accountType,omitempty"`           // 账户类型：01(对私账户)/02(对公账户)
	ContactLine           string        `json:"contactLine,omitempty"`           // 联行号
	BranchName            string        `json:"branchName,omitempty"`            // 开户支行名称
	BranchProvince        string        `json:"branchProvince,omitempty"`        // 开户支行所在省
	BranchCity            string        `json:"branchCity,omitempty"`            // 开户支行所在市
	CertType              string        `json:"certType,omitempty"`              // 持卡人证件类型：01(身份证)
	CertNo                string        `json:"certNo,omitempty"`                // 持卡人证件号码
	CardHolderAddress     string        `json:"cardHolderAddress,omitempty"`     // 持卡人地址
	LogonId               string        `json:"logonId,omitempty"`               // 支付宝登陆账号
	UserId                string        `json:"userId,omitempty"`                // 支付宝用户 ID
	FileManifest          *FileManifest `json:"fileManifest,omitempty"`          // 文件哈希清单（非必填）：仅包含有新文件的字段，未提供的字段草稿保留原有文件；无需更新文件时整体省略
}

func (x *UpdatePrepareDto) GetTs() string {
	return x.ReqTimestamp
}

// NewUpdatePrepareDto 创建更新草稿并申请上传凭证请求体。
// 除 draftId 外的必填参数与 CreatePrepareDto 完全一致；
// 业务字段须整体重新提交（非增量更新），与创建请求的必填规则相同。
func NewUpdatePrepareDto(
	draftId string,
	legalName string,
	shortName string,
	productCode []string,
	merchantBaseType string,
	subRoleType string,
	dealType string,
	mcc string,
	contactMobile string,
	contactName string,
	email string,
) *UpdatePrepareDto {
	return &UpdatePrepareDto{
		ReqTimestamp:     strconv.FormatInt(time.Now().UnixMilli(), 10),
		DraftId:          draftId,
		LegalName:        legalName,
		ShortName:        shortName,
		ProductCode:      productCode,
		MerchantBaseType: merchantBaseType,
		SubRoleType:      subRoleType,
		DealType:         dealType,
		Mcc:              mcc,
		ContactMobile:    contactMobile,
		ContactName:      contactName,
		Email:            email,
	}
}

// SetPrincipal 一次性设置负责人相关信息（手机号、证件类型、证件号码、姓名、证件有效期）。
func (x *UpdatePrepareDto) SetPrincipal(mobile, certType, certNo, person, certVld string) *UpdatePrepareDto {
	x.PrincipalMobile = mobile
	x.PrincipalCertType = certType
	x.PrincipalCertNo = certNo
	x.PrincipalPerson = person
	x.PrincipalCertVld = certVld
	return x
}

// SetLocation 一次性设置地址信息（省/市/区/详细地址/客服电话）。
func (x *UpdatePrepareDto) SetLocation(province, city, district, address, servicePhoneNo string) *UpdatePrepareDto {
	x.Province = province
	x.City = city
	x.District = district
	x.Address = address
	x.ServicePhoneNo = servicePhoneNo
	return x
}

// SetPerson 一次性设置自然人商户相关信息（性别、职业）。
func (x *UpdatePrepareDto) SetPerson(sex, profession string) *UpdatePrepareDto {
	x.PersonSex = sex
	x.PersonProfession = profession
	return x
}

// SetSettlementAccountType 设置结算类型。
func (x *UpdatePrepareDto) SetSettlementAccountType(i string) *UpdatePrepareDto {
	x.SettlementAccountType = i
	return x
}

// SetBank 一次性设置结算账户信息
// （bankCardNo 银行卡号、accountType 账户类型 01对私/02对公、
// branchName 开户支行名称、branchProvince 开户支行所在省、branchCity 开户支行所在市）。
// 银行账户户名 bankCertName 仅对公账户（02）需要，经 SetBankCertName 单独设置；
// 联行号 contactLine 为非必填，经 SetContactLine 单独设置。
func (x *UpdatePrepareDto) SetBank(bankCardNo, accountType, branchName, branchProvince, branchCity string) *UpdatePrepareDto {
	x.BankCardNo = bankCardNo
	x.AccountType = accountType
	x.BranchName = branchName
	x.BranchProvince = branchProvince
	x.BranchCity = branchCity
	return x
}

// SetBankCertName 设置银行账户户名（仅对公账户 02 需要）。
func (x *UpdatePrepareDto) SetBankCertName(i string) *UpdatePrepareDto {
	x.BankCertName = i
	return x
}

// SetContactLine 设置联行号（非必填）。
func (x *UpdatePrepareDto) SetContactLine(i string) *UpdatePrepareDto {
	x.ContactLine = i
	return x
}

// SetCardHolder 一次性设置持卡人信息（certType 证件类型、certNo 证件号码、cardHolderAddress 持卡人地址）。
func (x *UpdatePrepareDto) SetCardHolder(certType, certNo, cardHolderAddress string) *UpdatePrepareDto {
	x.CertType = certType
	x.CertNo = certNo
	x.CardHolderAddress = cardHolderAddress
	return x
}

// SetLogonId 设置支付宝登录账号。
func (x *UpdatePrepareDto) SetLogonId(i string) *UpdatePrepareDto {
	x.LogonId = i
	return x
}

// SetUserId 设置支付宝用户 ID。
func (x *UpdatePrepareDto) SetUserId(i string) *UpdatePrepareDto {
	x.UserId = i
	return x
}

// SetFileManifest 设置文件哈希清单（非必填）。
// 仅需包含有新文件的字段，未提供的字段草稿保留原有文件；无需更新文件时可不调用。
func (x *UpdatePrepareDto) SetFileManifest(i *FileManifest) *UpdatePrepareDto {
	x.FileManifest = i
	return x
}

// SetMangerLogonId 设置管理员支付宝登录号。
func (x *UpdatePrepareDto) SetMangerLogonId(i string) *UpdatePrepareDto {
	x.MangerLogonId = i
	return x
}

// SetBussAuthVld 设置营业执照有效期，格式 yyyy-MM-dd HH:mm:ss。
func (x *UpdatePrepareDto) SetBussAuthVld(i string) *UpdatePrepareDto {
	x.BussAuthVld = i
	return x
}

// SetTaxNum 设置税务登记证号码。
func (x *UpdatePrepareDto) SetTaxNum(i string) *UpdatePrepareDto {
	x.TaxNum = i
	return x
}

// SetShareholder 一次性设置控股股东信息
// （shareholderName 姓名、shareholderCertType 证件类型、
// shareholderCertNo 证件号码、shareholderCertVld 证件有效期）。
func (x *UpdatePrepareDto) SetShareholder(shareholderName, shareholderCertType, shareholderCertNo, shareholderCertVld string) *UpdatePrepareDto {
	x.ShareholderName = shareholderName
	x.ShareholderCertType = shareholderCertType
	x.ShareholderCertNo = shareholderCertNo
	x.ShareholderCertVld = shareholderCertVld
	return x
}

// SetPersonCertVld 设置身份证件有效期限（自然人商户），格式 yyyy-MM-dd HH:mm:ss。
func (x *UpdatePrepareDto) SetPersonCertVld(i string) *UpdatePrepareDto {
	x.PersonCertVld = i
	return x
}

// SetBussAuthType 设置营业执照证件类型。
func (x *UpdatePrepareDto) SetBussAuthType(i string) *UpdatePrepareDto {
	x.BussAuthType = i
	return x
}

// SetBussAuthNo 设置证件号码（营业执照号或统一社会信用代码）。
func (x *UpdatePrepareDto) SetBussAuthNo(i string) *UpdatePrepareDto {
	x.BussAuthNo = i
	return x
}

// SetRemark 设置备注。
func (x *UpdatePrepareDto) SetRemark(i string) *UpdatePrepareDto {
	x.Remark = i
	return x
}

// SetPartnerId 设置渠道商 ID。
func (x *UpdatePrepareDto) SetPartnerId(i string) *UpdatePrepareDto {
	x.PartnerId = i
	return x
}

// SetMerchantType 设置商户业务类型，如 SCHOOL、GROUP_MEAL、HR 等。
func (x *UpdatePrepareDto) SetMerchantType(i string) *UpdatePrepareDto {
	x.MerchantType = i
	return x
}

// SetAgreementNo 设置支付宝签约记录编号（安全发）。
func (x *UpdatePrepareDto) SetAgreementNo(i string) *UpdatePrepareDto {
	x.AgreementNo = i
	return x
}

// SetAlipayPid 设置支付宝商户 ID。
func (x *UpdatePrepareDto) SetAlipayPid(i string) *UpdatePrepareDto {
	x.AlipayPid = i
	return x
}

// SetAlipayAccount 设置支付宝收款账号。
func (x *UpdatePrepareDto) SetAlipayAccount(i string) *UpdatePrepareDto {
	x.AlipayAccount = i
	return x
}

// SetLogicGroupId 设置支付宝学校、机构用户库 ID。
func (x *UpdatePrepareDto) SetLogicGroupId(i string) *UpdatePrepareDto {
	x.LogicGroupId = i
	return x
}

// SetWxSubMchId 设置微信商户号。
func (x *UpdatePrepareDto) SetWxSubMchId(i string) *UpdatePrepareDto {
	x.WxSubMchId = i
	return x
}

// SetWxSubMchAccount 设置微信收款账号。
func (x *UpdatePrepareDto) SetWxSubMchAccount(i string) *UpdatePrepareDto {
	x.WxSubMchAccount = i
	return x
}

type UpdatePrepareBizData struct {
	UploadToken   string `json:"uploadToken"`
	ExpireSeconds int64  `json:"expireSeconds"`
}

// UpdatePrepare 更新草稿并申请上传凭证（两步上传 Step 1）。
// 对已存在草稿提交修改后的业务字段与（若有新文件）文件哈希清单，换取一次性上传凭证；
// 凭证有效期内调用 UploadFiles 完成实际文件上传后，草稿数据才会被更新。
//
// 约束：
//   - 仅 EDITING（编辑中）或 FAILED（提交失败）状态的草稿可更新
//   - 业务字段须整体重新提交（非增量更新），必填规则与 CreatePrepare 相同
//   - fileManifest 非必填，仅需包含有新文件的字段；未提供的字段草稿保留原有文件
func (x *Hst) UpdatePrepare(ctx context.Context, dto *UpdatePrepareDto) (result *SignObjectRespResult[*UpdatePrepareBizData], err error) {
	dto.ChannelNo = x.Option.ChannelId

	var signObjectReq *SignObjectReq
	if signObjectReq, err = x.NewSignObjectReq(dto); err != nil {
		return
	}

	var b string
	if b, err = x.Request(ctx,
		"/channel/merchant_info_draft/update/prepare", signObjectReq); err != nil {
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
