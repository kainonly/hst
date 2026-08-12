package hst

import (
	"context"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
)

type GetUploadTokenDto struct {
	ReqTimestamp string `json:"reqTimestamp"` // 业务请求时间戳，须与外层信封 timestamp 一致
	PartnerId    string `json:"partnerId"`    // 渠道商 ID
	MerchantId   string `json:"merchantId"`   // 商户 ID
	FileName     string `json:"fileName"`     // 文件名（用于记录与日志）
	FileSM3Hash  string `json:"fileSM3Hash"`  // 文件内容 SM3 哈希（64 位十六进制）
}

func (x *GetUploadTokenDto) GetTs() string {
	return x.ReqTimestamp
}

// NewGetUploadTokenDto 创建申请文件上传凭证请求体。
// partnerId / merchantId / fileName / fileSM3Hash 为必填字段
// （reqTimestamp 由 GetUploadToken 方法自动填充）。
func NewGetUploadTokenDto(
	partnerId string,
	merchantId string,
	fileName string,
	fileSM3Hash string,
) *GetUploadTokenDto {
	return &GetUploadTokenDto{
		PartnerId:   partnerId,
		MerchantId:  merchantId,
		FileName:    fileName,
		FileSM3Hash: fileSM3Hash,
	}
}

type GetUploadTokenBizData struct {
	UploadToken   string `json:"uploadToken"`   // 上传凭证（UUID），用于上传文件鉴权
	ExpireSeconds int64  `json:"expireSeconds"` // 凭证有效期（秒），默认 300，过期需重新申请
}

// GetUploadToken 申请文件上传凭证（分账文件上传 Step 1）。
// 客户端预先计算文件 SM3 哈希并纳入签名体，服务端生成一次性 uploadToken（默认有效期 5 分钟）。
// uploadToken 一次性消费，无论上传成功与否都会失效；过期或失败需重新申请。
func (x *Hst) GetUploadToken(ctx context.Context, dto *GetUploadTokenDto) (result *SignObjectRespResult[*GetUploadTokenBizData], err error) {
	dto.ReqTimestamp = strconv.FormatInt(time.Now().UnixMilli(), 10)

	var signObjectReq *SignObjectReq
	if signObjectReq, err = x.NewSignObjectReq(dto); err != nil {
		return
	}

	var b string
	if b, err = x.Request(ctx,
		"/doc-trade-file/getUploadToken", signObjectReq); err != nil {
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
