package hst

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/bytedance/sonic"
	"github.com/kainonly/go/help"
	"resty.dev/v3"
)

// UploadFile 上传文件源。
// Path 与 Reader 二选一：设置 Reader 时以流方式上传（无需落盘），否则按 Path 打开本地文件。
// Name 为 multipart 文件名，Path 模式下留空自动取文件基名；
// Type 为 Content-Type，留空时由 resty 自动嗅探（读取前 512 字节）。
type UploadFile struct {
	Path   string
	Reader io.Reader
	Name   string
	Type   string
	Size   int64 // Reader 模式下的内容大小（可选，进度回调用）
}

// NewUploadFileFromPath 从本地文件路径构造上传文件。
func NewUploadFileFromPath(path string) *UploadFile {
	return &UploadFile{Path: path}
}

// NewUploadFileFromReader 从 io.Reader 构造上传文件（内存流、网络流等，无需落盘）。
// name 为 multipart 文件名（含扩展名，如 sfz-a.jpg）。
func NewUploadFileFromReader(name string, reader io.Reader) *UploadFile {
	return &UploadFile{Name: name, Reader: reader}
}

// NewUploadFileFromFileHeader 从 *multipart.FileHeader 构造上传文件，
// 可直接传入 Hertz / Gin 等框架 FormFile() 的返回值实现请求文件流式转发，无需落盘。
// 注意：FileHeader 打开的流是一次性的，构造后不可重复使用。
func NewUploadFileFromFileHeader(fh *multipart.FileHeader) (file *UploadFile, err error) {
	var f multipart.File
	if f, err = fh.Open(); err != nil {
		return
	}
	file = &UploadFile{
		Name:   fh.Filename,
		Type:   fh.Header.Get("Content-Type"),
		Size:   fh.Size,
		Reader: f,
	}
	return
}

// field 生成 resty multipart 字段。
func (x *UploadFile) field(name string) *resty.MultipartField {
	return &resty.MultipartField{
		Name:        name,
		FileName:    x.Name,
		ContentType: x.Type,
		Reader:      x.Reader,
		FilePath:    x.Path,
		FileSize:    x.Size,
	}
}

// UploadFilesDto 上传资质文件请求体（multipart/form-data）。
// 每个字段对应一种资质，值为按上传顺序排列的上传文件列表，
// 支持 NewUploadFileFromPath / NewUploadFileFromReader / NewUploadFileFromFileHeader 构造。
// 仅 fileManifest 中出现且哈希数量 > 0 的字段必须按顺序上传对应数量的文件；
// 清单中为空或未出现的字段不可上传（否则视为夹带未签名文件而拒绝）。
type UploadFilesDto struct {
	ChannelId                   string
	UploadToken                 string
	CertPhotoAFiles             []*UploadFile // 身份证人像面
	CertPhotoBFiles             []*UploadFile // 身份证国徽面
	LicensePhotoFiles           []*UploadFile // 营业执照
	PrgPhotoFiles               []*UploadFile // 组织机构代码证
	IndustryLicensePhotoFiles   []*UploadFile // 开户许可证
	ShopPhotoFiles              []*UploadFile // 门头照
	OtherPhotoFiles             []*UploadFile // 其他资料
	CertPhotoCFiles             []*UploadFile // 手持身份证
	RegisterProtocolPhotoFiles  []*UploadFile // 商户入驻协议
	ContractPhotoFiles          []*UploadFile // 租赁协议
	ShopEntrancePhotoFiles      []*UploadFile // 门店内景
	CheckstandPhotoFiles        []*UploadFile // 收银台
	MerchantAgreementPhotoFiles []*UploadFile // 商户协议
}

// NewUploadFilesDto 创建上传资质文件请求体。
func NewUploadFilesDto(uploadToken string) *UploadFilesDto {
	return &UploadFilesDto{UploadToken: uploadToken}
}

// SetCertPhotoAFiles 设置身份证人像面文件列表。
func (x *UploadFilesDto) SetCertPhotoAFiles(i ...*UploadFile) *UploadFilesDto {
	x.CertPhotoAFiles = i
	return x
}

// SetCertPhotoBFiles 设置身份证国徽面文件列表。
func (x *UploadFilesDto) SetCertPhotoBFiles(i ...*UploadFile) *UploadFilesDto {
	x.CertPhotoBFiles = i
	return x
}

// SetLicensePhotoFiles 设置营业执照文件列表。
func (x *UploadFilesDto) SetLicensePhotoFiles(i ...*UploadFile) *UploadFilesDto {
	x.LicensePhotoFiles = i
	return x
}

// SetPrgPhotoFiles 设置组织机构代码证文件列表。
func (x *UploadFilesDto) SetPrgPhotoFiles(i ...*UploadFile) *UploadFilesDto {
	x.PrgPhotoFiles = i
	return x
}

// SetIndustryLicensePhotoFiles 设置开户许可证文件列表。
func (x *UploadFilesDto) SetIndustryLicensePhotoFiles(i ...*UploadFile) *UploadFilesDto {
	x.IndustryLicensePhotoFiles = i
	return x
}

// SetShopPhotoFiles 设置门头照文件列表。
func (x *UploadFilesDto) SetShopPhotoFiles(i ...*UploadFile) *UploadFilesDto {
	x.ShopPhotoFiles = i
	return x
}

// SetOtherPhotoFiles 设置其他资料文件列表。
func (x *UploadFilesDto) SetOtherPhotoFiles(i ...*UploadFile) *UploadFilesDto {
	x.OtherPhotoFiles = i
	return x
}

// SetCertPhotoCFiles 设置手持身份证文件列表。
func (x *UploadFilesDto) SetCertPhotoCFiles(i ...*UploadFile) *UploadFilesDto {
	x.CertPhotoCFiles = i
	return x
}

// SetRegisterProtocolPhotoFiles 设置商户入驻协议文件列表。
func (x *UploadFilesDto) SetRegisterProtocolPhotoFiles(i ...*UploadFile) *UploadFilesDto {
	x.RegisterProtocolPhotoFiles = i
	return x
}

// SetContractPhotoFiles 设置租赁协议文件列表。
func (x *UploadFilesDto) SetContractPhotoFiles(i ...*UploadFile) *UploadFilesDto {
	x.ContractPhotoFiles = i
	return x
}

// SetShopEntrancePhotoFiles 设置门店内景文件列表。
func (x *UploadFilesDto) SetShopEntrancePhotoFiles(i ...*UploadFile) *UploadFilesDto {
	x.ShopEntrancePhotoFiles = i
	return x
}

// SetCheckstandPhotoFiles 设置收银台文件列表。
func (x *UploadFilesDto) SetCheckstandPhotoFiles(i ...*UploadFile) *UploadFilesDto {
	x.CheckstandPhotoFiles = i
	return x
}

// SetMerchantAgreementPhotoFiles 设置商户协议文件列表。
func (x *UploadFilesDto) SetMerchantAgreementPhotoFiles(i ...*UploadFile) *UploadFilesDto {
	x.MerchantAgreementPhotoFiles = i
	return x
}

// multipartFields 按 fileManifest 顺序生成 multipart 文件字段列表。
// 同一字段多个文件按 slice 顺序生成多个 MultipartField（Name 相同）。
// channelId 与 uploadToken 由 UploadFiles 方法用 SetFormData 设置。
func (x *UploadFilesDto) multipartFields() []*resty.MultipartField {
	var fields []*resty.MultipartField
	fileGroups := []struct {
		name  string
		files []*UploadFile
	}{
		{"certPhotoAFiles", x.CertPhotoAFiles},
		{"certPhotoBFiles", x.CertPhotoBFiles},
		{"licensePhotoFiles", x.LicensePhotoFiles},
		{"prgPhotoFiles", x.PrgPhotoFiles},
		{"industryLicensePhotoFiles", x.IndustryLicensePhotoFiles},
		{"shopPhotoFiles", x.ShopPhotoFiles},
		{"otherPhotoFiles", x.OtherPhotoFiles},
		{"certPhotoCFiles", x.CertPhotoCFiles},
		{"registerProtocolPhotoFiles", x.RegisterProtocolPhotoFiles},
		{"contractPhotoFiles", x.ContractPhotoFiles},
		{"shopEntrancePhotoFiles", x.ShopEntrancePhotoFiles},
		{"checkstandPhotoFiles", x.CheckstandPhotoFiles},
		{"merchantAgreementPhotoFiles", x.MerchantAgreementPhotoFiles},
	}
	for _, g := range fileGroups {
		for _, f := range g.files {
			if f == nil {
				continue
			}
			fields = append(fields, f.field(g.name))
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
