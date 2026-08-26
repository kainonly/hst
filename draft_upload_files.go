package hst

import (
	"bytes"
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/kainonly/go/help"
	"resty.dev/v3"
)

// UploadFile 上传文件（文件名 + 内容字节）。
// 典型场景：业务方在自己的接口中已收到用户上传的文件（[]byte），
// Step 1 用同一份字节计算 SM3 哈希提交 FileManifest，
// Step 2 直接透传该字节构造本类型，全程无需落盘。
type UploadFile struct {
	Name string // multipart 文件名（含扩展名，如 sfz-a.jpg）
	Data []byte // 文件内容，须与 Step 1 计算 SM3 哈希的字节完全一致
}

// NewUploadFile 创建上传文件。
// name 为 multipart 文件名（含扩展名），data 为文件内容字节。
func NewUploadFile(name string, data []byte) *UploadFile {
	return &UploadFile{Name: name, Data: data}
}

// field 生成 resty multipart 字段。
// Content-Type 留空，由 resty 读取前 512 字节自动嗅探。
func (x *UploadFile) field(name string) *resty.MultipartField {
	return &resty.MultipartField{
		Name:     name,
		FileName: x.Name,
		Reader:   bytes.NewReader(x.Data),
		FileSize: int64(len(x.Data)),
	}
}

// uploadFileFields 合法的 multipart 文件字段名，与 FileManifest 字段一一对应，
// 按文档 fileManifest 顺序排列（服务端按字段名定位、字段内按顺序比对 SM3）。
var uploadFileFields = []string{
	"certPhotoAFiles",             // 身份证人像面
	"certPhotoBFiles",             // 身份证国徽面
	"licensePhotoFiles",           // 营业执照
	"prgPhotoFiles",               // 组织机构代码证
	"industryLicensePhotoFiles",   // 开户许可证
	"shopPhotoFiles",              // 门头照
	"otherPhotoFiles",             // 其他资料
	"certPhotoCFiles",             // 手持身份证
	"registerProtocolPhotoFiles",  // 商户入驻协议
	"contractPhotoFiles",          // 租赁协议
	"shopEntrancePhotoFiles",      // 门店内景
	"checkstandPhotoFiles",        // 收银台
	"merchantAgreementPhotoFiles", // 商户协议
}

// UploadFilesDto 上传资质文件请求体（multipart/form-data）。
// 文件按 fileManifest 字段名组织，仅 fileManifest 中出现且哈希数量 > 0 的字段
// 必须按顺序上传对应数量的文件；清单中为空或未出现的字段不可上传
// （否则视为夹带未签名文件而拒绝）。
type UploadFilesDto struct {
	ChannelId   string
	UploadToken string
	files       map[string][]*UploadFile
}

// NewUploadFilesDto 创建上传资质文件请求体。
func NewUploadFilesDto(uploadToken string) *UploadFilesDto {
	return &UploadFilesDto{
		UploadToken: uploadToken,
		files:       make(map[string][]*UploadFile),
	}
}

// SetFiles 按 fileManifest 字段名设置该资质的文件列表（同字段重复调用为覆盖）。
// field 须为 fileManifest 字段名之一（certPhotoAFiles、licensePhotoFiles 等），
// files 按上传顺序排列；UploadFiles 会对非法字段名返回错误。
func (x *UploadFilesDto) SetFiles(field string, files ...*UploadFile) *UploadFilesDto {
	x.files[field] = files
	return x
}

// validateFields 校验已设置的字段名全部合法，防止拼写错误导致文件被静默丢弃。
func (x *UploadFilesDto) validateFields() error {
	valid := make(map[string]struct{}, len(uploadFileFields))
	for _, name := range uploadFileFields {
		valid[name] = struct{}{}
	}
	for name := range x.files {
		if _, ok := valid[name]; !ok {
			return fmt.Errorf(`未知上传字段名 %q，须为 fileManifest 字段之一`, name)
		}
	}
	return nil
}

// multipartFields 按 fileManifest 字段顺序生成 multipart 文件字段列表。
// 同一字段多个文件按 slice 顺序生成多个 MultipartField（Name 相同）。
// channelId 与 uploadToken 由 UploadFiles 方法用 SetFormData 设置。
func (x *UploadFilesDto) multipartFields() []*resty.MultipartField {
	var fields []*resty.MultipartField
	for _, name := range uploadFileFields {
		for _, f := range x.files[name] {
			if f == nil {
				continue
			}
			fields = append(fields, f.field(name))
		}
	}
	return fields
}

// UploadFilesBizData 上传资质文件响应 bizData（MerchantInfoDraftVO 草稿视图对象）。
type UploadFilesBizData struct {
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

// UploadFilesResp 上传资质文件响应（普通 JSON，非加密信封）。
type UploadFilesResp struct {
	Code string              `json:"code"` // 网关响应码，SUCCESS 表示网关受理成功
	Msg  string              `json:"msg"`  // 网关响应消息
	Data UploadFilesRespData `json:"data"`
}

// UploadFilesRespData 上传资质文件响应 data 字段（含业务结果）。
type UploadFilesRespData struct {
	BizSuccess bool                `json:"bizSuccess"`
	BizCode    string              `json:"bizCode"`
	BizMsg     string              `json:"bizMsg"`
	BizData    *UploadFilesBizData `json:"bizData"`
}

// UploadFiles 上传资质文件（两步上传 Step 2）。
// 以 multipart/form-data 携带 uploadToken 与实际资质文件，
// 服务端按字段名 + 数组顺序逐一重算 SM3 比对，全部吻合后写入 / 更新草稿库。
//
// 与标准 JSON 加密流程相互独立，响应为普通 JSON（非加密信封）。
// uploadToken 为一次性凭证，无论本次上传成功与否都会被消费；上传失败需从 Step 1 重新开始。
func (x *Hst) UploadFiles(ctx context.Context, dto *UploadFilesDto) (bizData *UploadFilesBizData, err error) {
	if err = dto.validateFields(); err != nil {
		return
	}
	dto.ChannelId = x.Option.ChannelId

	// channel-multi-file 前缀，拼接完整 URL
	uploadURL := x.Option.BaseURL + "/channel-multi-file/merchant_info_draft/upload_files"
	var resp *resty.Response
	if resp, err = x.Client.R().SetContext(ctx).
		SetFormData(map[string]string{
			"channelId":   dto.ChannelId,
			"uploadToken": dto.UploadToken,
		}).
		SetMultipartFields(dto.multipartFields()...).
		Post(uploadURL); err != nil {
		return
	}

	if resp.StatusCode() != 200 {
		err = help.E(0, fmt.Sprintf(`第三方接口响应失败! status=%d body=%s`,
			resp.StatusCode(), resp.String()))
		return
	}

	var uploadResp *UploadFilesResp
	if err = sonic.Unmarshal(resp.Bytes(), &uploadResp); err != nil {
		return
	}

	if uploadResp.Code != "SUCCESS" {
		err = help.E(0, fmt.Sprintf(`第三方接口响应失败! code=%s msg=%s`,
			uploadResp.Code, uploadResp.Msg))
		return
	}

	if !uploadResp.Data.BizSuccess {
		err = bizError(uploadResp.Data.BizCode, uploadResp.Data.BizMsg)
		return
	}

	bizData = uploadResp.Data.BizData
	return
}
