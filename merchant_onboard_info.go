package hst

import (
	"context"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
)

type MerchantOnboardInfoDto struct {
	ReqTimestamp string `json:"reqTimestamp"` // 业务请求时间戳，须与外层信封 timestamp 一致（防重放）
	PartnerId    string `json:"partnerId"`    // 渠道商 ID（与信封 channelId 不是同一个字段，必须显式传入）
	MerchantId   string `json:"merchantId"`   // 商户号，须归属于 partnerId，否则返回「商户不存在或无权访问」
}

func (x *MerchantOnboardInfoDto) GetTs() string {
	return x.ReqTimestamp
}

// NewMerchantOnboardInfoDto 创建查询商户入驻信息请求体。
func NewMerchantOnboardInfoDto(merchantId string) *MerchantOnboardInfoDto {
	return &MerchantOnboardInfoDto{
		MerchantId:   merchantId,
		ReqTimestamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
	}
}

// MerchantOnboardInfoBizData 查询商户入驻信息响应 bizData。
// 响应只包含有值的字段：merchantBaseInfo（主体资料不存在）、mybankInfo（未发起过网商入驻）
// 不返回时为 nil，请按指针是否为 nil 判断，不要依赖空对象。
type MerchantOnboardInfoBizData struct {
	DraftId            string              `json:"draftId"`            // 入驻草稿 ID，即入驻结果通知中的 applyNo；非草稿流程创建的商户不返回
	MerchantInfo       *MerchantInfo       `json:"merchantInfo"`       // 商户信息
	MerchantBaseInfo   *MerchantBaseInfo   `json:"merchantBaseInfo"`   // 主体基础信息，主体资料不存在时为 nil
	SettlementAccounts []SettlementAccount `json:"settlementAccounts"` // 结算账户列表，默认账户在前；无结算账户时为空数组
	MybankInfo         *MybankInfo         `json:"mybankInfo"`         // 网商银行入驻信息（最近一笔申请），未发起过网商入驻时为 nil
}

// MerchantInfo 商户信息。
type MerchantInfo struct {
	MerchantId      string `json:"merchantId"`      // 商户号
	PartnerId       string `json:"partnerId"`       // 渠道商 ID
	MerchantType    string `json:"merchantType"`    // 商户业务类型，如 SCHOOL、GROUP_MEAL、HR
	LegalName       string `json:"legalName"`       // 法定名称
	ShortName       string `json:"shortName"`       // 商户简称/常用名
	AlipayPid       string `json:"alipayPid"`       // 支付宝商户 ID
	AlipayAccount   string `json:"alipayAccount"`   // 支付宝收款账号
	WxSubMchId      string `json:"wxSubMchId"`      // 微信商户号
	WxSubMchAccount string `json:"wxSubMchAccount"` // 微信收款账号
	LogicGroupId    string `json:"logicGroupId"`    // 支付宝学校、机构用户库 ID
	CreateTime      string `json:"createTime"`      // 创建时间，格式 yyyy-MM-dd HH:mm:ss
	UpdateTime      string `json:"updateTime"`      // 更新时间，格式 yyyy-MM-dd HH:mm:ss
}

// MerchantBaseInfo 主体基础信息。
// 字段与 CreatePrepareDto 的主体资料字段一致；
// 同一主体（同一营业执照或负责人证件）下的商户共用一份主体资料。
// certPhotoA ~ otherPhoto 为可直接访问的预签名地址，有效期 7 天，
// 同一字段多张图片用英文逗号分隔。
type MerchantBaseInfo struct {
	LegalName              string `json:"legalName"`              // 法定名称（营业执照上的公司名称）
	ShortName              string `json:"shortName"`              // 商户简称
	MerchantBaseType       string `json:"merchantBaseType"`       // 商户类型：01(自然人)/02(个体工商户)/03(企业商户)
	SubRoleType            string `json:"subRoleType"`            // 商户角色：shop_franchisee加盟商门店/shop_direct直营门店/shop_joint联营门店/supplier_brandSupplyChain品牌供应链公司/supplier_external外部供应商/brand_manage品牌管理公司/brand_area品牌区域公司/brand_other其他服务公司
	DealType               string `json:"dealType"`               // 商户经营类型：01(实体特约)/02(网络特约)/03(实体兼网络特约)
	Mcc                    string `json:"mcc"`                    // 经营类目（MCC 码）
	BusinessScope          string `json:"businessScope"`          // 经营范围代码
	BussAuthType           string `json:"bussAuthType"`           // 证照类型：12营业执照/统一社会信用代码/41政府机关证件/42部队证件/43社会团体证件/44事业单位证件/45民办非企业组织证件/99其它企业类型证件
	BussAuthNo             string `json:"bussAuthNo"`             // 证照号码（营业执照号或统一社会信用代码）
	BussAuthStart          string `json:"bussAuthStart"`          // 证照有效期起
	BussAuthVld            string `json:"bussAuthVld"`            // 证照有效期止
	TaxNum                 string `json:"taxNum"`                 // 税务登记证号码
	PrincipalPerson        string `json:"principalPerson"`        // 负责人名称或企业法人代表姓名
	PrincipalMobile        string `json:"principalMobile"`        // 负责人手机号
	PrincipalCertType      string `json:"principalCertType"`      // 负责人证件类型：100身份证/102护照/108外国人永久居留身份证/114台湾居民来往大陆通行证/115港澳居民来往内地通行证/116港澳居民居住证/117台湾居民居住证
	PrincipalCertNo        string `json:"principalCertNo"`        // 负责人证件号码
	PrincipalCertStart     string `json:"principalCertStart"`     // 负责人证件有效期起
	PrincipalCertVld       string `json:"principalCertVld"`       // 负责人证件有效期止
	ContactName            string `json:"contactName"`            // 联系人姓名
	ContactMobile          string `json:"contactMobile"`          // 联系人手机号
	Email                  string `json:"email"`                  // 邮箱
	MangerLogonId          string `json:"mangerLogonId"`          // 管理员支付宝登录号
	ServicePhoneNo         string `json:"servicePhoneNo"`         // 商户客服电话
	ShareholderName        string `json:"shareholderName"`        // 控股股东或实际控制人姓名
	ShareholderCertType    string `json:"shareholderCertType"`    // 控股股东或实际控制人证件类型
	ShareholderCertNo      string `json:"shareholderCertNo"`      // 控股股东或实际控制人证件号码
	ShareholderCertVld     string `json:"shareholderCertVld"`     // 控股股东或实际控制人证件有效期
	PersonSex              string `json:"personSex"`              // 性别（自然人商户）：M(男性)/F(女性)
	PersonProfession       string `json:"personProfession"`       // 职业（自然人商户）
	PersonCertVld          string `json:"personCertVld"`          // 证件有效期（自然人商户）
	Province               string `json:"province"`               // 省（国标省市区号）
	City                   string `json:"city"`                   // 市（国标省市区号）
	District               string `json:"district"`               // 区（国标省市区号）
	Address                string `json:"address"`                // 详细地址
	CertPhotoA             string `json:"certPhotoA"`             // 负责人证件人像面（预签名地址，7 天内有效）
	CertPhotoB             string `json:"certPhotoB"`             // 负责人证件国徽面（预签名地址，7 天内有效）
	CertPhotoC             string `json:"certPhotoC"`             // 手持身份证照片（预签名地址，7 天内有效）
	LicensePhoto           string `json:"licensePhoto"`           // 营业执照（预签名地址，7 天内有效）
	PrgPhoto               string `json:"prgPhoto"`               // 组织机构代码证（预签名地址，7 天内有效）
	IndustryLicensePhoto   string `json:"industryLicensePhoto"`   // 开户许可证（预签名地址，7 天内有效）
	ShopPhoto              string `json:"shopPhoto"`              // 门头照（预签名地址，7 天内有效）
	ShopEntrancePhoto      string `json:"shopEntrancePhoto"`      // 门店内景照（预签名地址，7 天内有效）
	CheckstandPhoto        string `json:"checkstandPhoto"`        // 收银台照片（预签名地址，7 天内有效）
	ContractPhoto          string `json:"contractPhoto"`          // 租赁协议/产权证明（预签名地址，7 天内有效）
	RegisterProtocolPhoto  string `json:"registerProtocolPhoto"`  // 餐饮平台照片/入驻协议照片（预签名地址，7 天内有效）
	MerchantAgreementPhoto string `json:"merchantAgreementPhoto"` // 商户协议照片（预签名地址，7 天内有效）
	OtherPhoto             string `json:"otherPhoto"`             // 其他照片（预签名地址，7 天内有效）
	CreateTime             string `json:"createTime"`             // 创建时间，格式 yyyy-MM-dd HH:mm:ss
	UpdateTime             string `json:"updateTime"`             // 更新时间，格式 yyyy-MM-dd HH:mm:ss
}

// SettlementAccount 结算账户。
type SettlementAccount struct {
	AccountId             string `json:"accountId"`             // 结算账户 ID
	MerchantId            string `json:"merchantId"`            // 所属商户号
	SettlementAccountType string `json:"settlementAccountType"` // 结算类型：01(银行卡)/02(支付宝)/03(支付宝虚拟账户)
	BankCardNo            string `json:"bankCardNo"`            // 银行卡号
	BankCertName          string `json:"bankCertName"`          // 银行账户户名
	AccountType           string `json:"accountType"`           // 账户类型：01(对私账户)/02(对公账户)
	ContactLine           string `json:"contactLine"`           // 联行号
	BranchName            string `json:"branchName"`            // 开户支行名称
	BranchProvince        string `json:"branchProvince"`        // 开户支行所在省
	BranchCity            string `json:"branchCity"`            // 开户支行所在市
	CertType              string `json:"certType"`              // 持卡人证件类型：01(身份证)
	CertNo                string `json:"certNo"`                // 持卡人证件号码
	CardHolderAddress     string `json:"cardHolderAddress"`     // 持卡人地址
	LogonId               string `json:"logonId"`               // 支付宝登录账号
	UserId                string `json:"userId"`                // 支付宝用户 ID
	VerifyStatus          string `json:"verifyStatus"`          // 验证状态：PENDING(待验证)/VERIFIED(已验证)
	IsDefault             string `json:"isDefault"`             // 是否默认提现账户：T(是)/F(否)
	CreateTime            string `json:"createTime"`            // 创建时间，格式 yyyy-MM-dd HH:mm:ss
	UpdateTime            string `json:"updateTime"`            // 更新时间，格式 yyyy-MM-dd HH:mm:ss
}

// MybankInfo 网商银行入驻信息（最近一笔申请）。
type MybankInfo struct {
	MybankOrderNo    string `json:"mybankOrderNo"`    // 网商银行申请单号；入驻前图片被网商拒绝时不返回
	MybankMerchantId string `json:"mybankMerchantId"` // 网商银行商户号（入驻成功/待激活后返回）
	Status           string `json:"status"`           // 入驻状态码：0-审核中，1-成功，2-失败，3-待激活，4-激活中
	StatusDesc       string `json:"statusDesc"`       // 入驻状态描述
	FailReason       string `json:"failReason"`       // 失败原因（status = 2 时返回）
	CreateTime       string `json:"createTime"`       // 申请时间，格式 yyyy-MM-dd HH:mm:ss
	UpdateTime       string `json:"updateTime"`       // 最近更新时间，格式 yyyy-MM-dd HH:mm:ss
}

// MerchantOnboardInfo 查询商户入驻信息。
// 按商户号查询商户的全部入驻资料：商户信息、主体基础信息（证照、负责人、联系人、地址、证件图片地址）、
// 全部结算账户；商户发起过网商银行入驻时，附带最近一笔网商入驻信息（含入驻状态与失败原因）。
// 可用于核对进件资料，或在收到网商入驻失败通知后查看当前资料再修改重提。
//
// 响应只包含有值的字段，请按字段存在性判断，不要依赖空字符串：
//   - MerchantBaseInfo / MybankInfo 为 nil 表示主体资料不存在 / 未发起过网商入驻
//   - SettlementAccounts 无结算账户时为空数组，默认账户在前
//   - 图片地址（CertPhotoA ~ OtherPhoto）为预签名地址，有效期 7 天，
//     同一字段多张图片用英文逗号分隔；请勿长期保存地址，过期后重新调用本接口获取
//
// mybankInfo.status 取值：
//   - 0 审核中、1 成功、2 失败、3 待激活、4 激活中
//   - 仅 status = 2（失败）时允许按 failReason 修改资料后用原 DraftId 经
//     UpdatePrepare → UploadFiles → Confirm 重提（商户号、结算账户 ID 保持不变），
//     其余状态不允许修改；重提结果通过 MYBANK_SUB_MERCHANT_ONBOARD_SUCCESS / FAIL 事件通知
func (x *Hst) MerchantOnboardInfo(ctx context.Context, dto *MerchantOnboardInfoDto) (result *SignObjectRespResult[*MerchantOnboardInfoBizData], signObjectResp *SignObjectResp, err error) {
	dto.PartnerId = x.Option.ChannelId

	var signObjectReq *SignObjectReq
	if signObjectReq, err = x.NewSignObjectReq(dto); err != nil {
		return
	}

	if signObjectResp, err = x.Request(ctx,
		"/channel/merchant/onboard-info", signObjectReq); err != nil {
		return
	}
	if err = sonic.UnmarshalString(signObjectResp.Body, &result); err != nil {
		return
	}
	if !result.BizSuccess {
		err = bizError(result.BizCode, result.BizMsg)
	}
	return
}
