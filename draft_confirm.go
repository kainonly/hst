package hst

import (
	"context"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
)

type ConfirmDto struct {
	ReqTimestamp string `json:"reqTimestamp"` // 业务请求时间戳，须与外层信封 timestamp 一致（防重放）
	ChannelNo    string `json:"channelNo"`    // 渠道号
	DraftId      string `json:"draftId"`      // 草稿 ID（由创建草稿或上传资质文件返回）
}

func (x *ConfirmDto) GetTs() string {
	return x.ReqTimestamp
}

func NewConfirmDto(draftId string) *ConfirmDto {
	return &ConfirmDto{
		DraftId:      draftId,
		ReqTimestamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
	}
}

// ConfirmBizData 确认提交草稿响应 bizData（MerchantInfoDraftVO 草稿视图对象）。
// 确认成功后 draftStatus = CONFIRMED，并回填 merchantId、orgId、accountId；
// 确认失败则 draftStatus = FAILED 并填充 failReason。
type ConfirmBizData struct {
	DraftId                string `json:"draftId"`                // 草稿唯一 ID
	DraftStatus            string `json:"draftStatus"`            // 草稿状态：EDITING/SUBMITTING/CONFIRMED/FAILED
	PartnerId              string `json:"partnerId"`              // 渠道商 ID
	MerchantId             string `json:"merchantId"`             // 商户 ID（确认后填充）
	OrgId                  string `json:"orgId"`                  // 企业唯一号（确认后填充）
	AccountId              string `json:"accountId"`              // 结算账户 ID（确认后填充）
	FailReason             string `json:"failReason"`             // 提交失败原因
	LegalName              string `json:"legalName"`              // 法定名称（营业执照上的公司名称）
	ShortName              string `json:"shortName"`              // 商户简称 / 常用名
	MerchantBaseType       string `json:"merchantBaseType"`       // 商户类型：01/02/03
	SubRoleType            string `json:"subRoleType"`            // 商户角色
	DealType               string `json:"dealType"`               // 商户经营类型：01/02/03
	Mcc                    string `json:"mcc"`                    // 经营类目（MCC 码）
	ContactMobile          string `json:"contactMobile"`          // 联系人手机号
	ContactName            string `json:"contactName"`            // 联系人姓名
	Email                  string `json:"email"`                  // 邮箱
	PrincipalMobile        string `json:"principalMobile"`        // 负责人手机号
	PrincipalCertType      string `json:"principalCertType"`      // 负责人证件类型
	PrincipalCertNo        string `json:"principalCertNo"`        // 负责人证件号码
	PrincipalPerson        string `json:"principalPerson"`        // 负责人名称或企业法人代表姓名
	PrincipalCertVld       string `json:"principalCertVld"`       // 负责人证件有效期
	Province               string `json:"province"`               // 省
	City                   string `json:"city"`                   // 市
	District               string `json:"district"`               // 区
	Address                string `json:"address"`                // 详细地址
	ServicePhoneNo         string `json:"servicePhoneNo"`         // 商户客服电话
	TaxNum                 string `json:"taxNum"`                 // 税务登记证号码
	BussAuthVld            string `json:"bussAuthVld"`            // 营业执照有效期
	ShareholderName        string `json:"shareholderName"`        // 控股股东或实际控制人姓名
	ShareholderCertType    string `json:"shareholderCertType"`    // 控股股东或实际控制人证件类型
	ShareholderCertNo      string `json:"shareholderCertNo"`      // 控股股东或实际控制人证件号码
	ShareholderCertVld     string `json:"shareholderCertVld"`     // 控股股东或实际控制人证件有效期
	PersonSex              string `json:"personSex"`              // 性别（自然人商户）：M/F
	PersonProfession       string `json:"personProfession"`       // 职业（自然人商户）
	PersonCertVld          string `json:"personCertVld"`          // 身份证件有效期限（自然人商户）
	BussAuthType           string `json:"bussAuthType"`           // 证件类型
	BussAuthNo             string `json:"bussAuthNo"`             // 证件号码
	Remark                 string `json:"remark"`                 // 备注
	CertPhotoA             string `json:"certPhotoA"`             // 身份证人像面 Minio URL
	CertPhotoB             string `json:"certPhotoB"`             // 身份证国徽面 Minio URL
	LicensePhoto           string `json:"licensePhoto"`           // 营业执照 Minio URL
	PrgPhoto               string `json:"prgPhoto"`               // 组织机构代码证 Minio URL
	IndustryLicensePhoto   string `json:"industryLicensePhoto"`   // 开户许可证 Minio URL
	ShopPhoto              string `json:"shopPhoto"`              // 门头照 Minio URL
	OtherPhoto             string `json:"otherPhoto"`             // 其他资料 Minio URL
	CertPhotoC             string `json:"certPhotoC"`             // 手持身份证 Minio URL
	RegisterProtocolPhoto  string `json:"registerProtocolPhoto"`  // 商户入驻协议 Minio URL
	ContractPhoto          string `json:"contractPhoto"`          // 租赁协议 Minio URL
	ShopEntrancePhoto      string `json:"shopEntrancePhoto"`      // 门店内景 Minio URL
	CheckstandPhoto        string `json:"checkstandPhoto"`        // 收银台 Minio URL
	MerchantAgreementPhoto string `json:"merchantAgreementPhoto"` // 商户协议 Minio URL
	MerchantType           string `json:"merchantType"`           // 商户业务类型
	AgreementNo            string `json:"agreementNo"`            // 支付宝签约记录编号（安全发）
	AlipayPid              string `json:"alipayPid"`              // 支付宝商户 ID
	AlipayAccount          string `json:"alipayAccount"`          // 支付宝收款账号
	LogicGroupId           string `json:"logicGroupId"`           // 支付宝学校、机构用户库 ID
	WxSubMchId             string `json:"wxSubMchId"`             // 微信商户号
	WxSubMchAccount        string `json:"wxSubMchAccount"`        // 微信收款账号
	SettlementAccountType  string `json:"settlementAccountType"`  // 结算类型：01/02/03
	BankCardNo             string `json:"bankCardNo"`             // 银行卡号
	BankCertName           string `json:"bankCertName"`           // 银行账户户名
	AccountType            string `json:"accountType"`            // 账户类型：01/02
	ContactLine            string `json:"contactLine"`            // 联行号
	BranchName             string `json:"branchName"`             // 开户支行名称
	BranchProvince         string `json:"branchProvince"`         // 开户支行所在省
	BranchCity             string `json:"branchCity"`             // 开户支行所在市
	CertType               string `json:"certType"`               // 持卡人证件类型
	CertNo                 string `json:"certNo"`                 // 持卡人证件号码
	CardHolderAddress      string `json:"cardHolderAddress"`      // 持卡人地址
	LogonId                string `json:"logonId"`                // 支付宝登陆账号
	UserId                 string `json:"userId"`                 // 支付宝用户 ID
	CreateTime             string `json:"createTime"`             // 创建时间，格式 yyyy-MM-dd HH:mm:ss
	UpdateTime             string `json:"updateTime"`             // 更新时间，格式 yyyy-MM-dd HH:mm:ss
}

// Confirm 确认提交草稿。
// 将两步上传完成后的草稿数据分发写入商户基础信息、商户信息、结算账户三个正式业务表，
// 触发商户入驻申请。确认后草稿状态变更为 CONFIRMED 或 FAILED；
// 若为 FAILED，可通过 UpdatePrepare 修改后重新走两步上传流程再次确认。
func (x *Hst) Confirm(ctx context.Context, dto *ConfirmDto) (result *SignObjectRespResult[*ConfirmBizData], err error) {
	dto.ChannelNo = x.Option.ChannelId

	var signObjectReq *SignObjectReq
	if signObjectReq, err = x.NewSignObjectReq(dto); err != nil {
		return
	}

	var b string
	if b, err = x.Request(ctx,
		"/channel/merchant_info_draft/confirm", signObjectReq); err != nil {
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
