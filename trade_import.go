package hst

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/bytedance/sonic"
	"github.com/kainonly/go/help"
	"resty.dev/v3"
)

// TradeImportDto 上传交易订单文件请求体（multipart/form-data）。
type TradeImportDto struct {
	ChannelId   string // 渠道商 ID（用于凭证归属校验）
	UploadToken string // 上传凭证（由 GetUploadToken 颁发）
	FilePath    string // XLSX 文件本地路径
}

// NewTradeImportDto 创建上传交易订单文件请求体。
// channelId / uploadToken / filePath 为必填字段。
func NewTradeImportDto(channelId string, uploadToken string, filePath string) *TradeImportDto {
	return &TradeImportDto{
		ChannelId:   channelId,
		UploadToken: uploadToken,
		FilePath:    filePath,
	}
}

// multipartFields 生成 multipart 文件字段列表。
// channelId / uploadToken 由 TradeImport 方法用 SetFormData 设置。
func (x *TradeImportDto) multipartFields() []*resty.MultipartField {
	return []*resty.MultipartField{
		{
			Name:     "file",
			FileName: filepath.Base(x.FilePath),
			FilePath: x.FilePath,
		},
	}
}

// TradeImportResp 上传交易订单文件响应（普通 JSON，非加密信封）。
type TradeImportResp struct {
	Code string              `json:"code"` // 网关响应码，SUCCESS 表示网关受理成功
	Msg  string              `json:"msg"`  // 网关响应消息
	Data TradeImportRespData `json:"data"`
}

// TradeImportRespData 上传交易订单文件响应 data 字段（含业务结果）。
// bizData 为业务主记录唯一 ID busId，用于后续确认 / 查询 / 取消。
type TradeImportRespData struct {
	BizSuccess bool   `json:"bizSuccess"`
	BizCode    string `json:"bizCode"`
	BizMsg     string `json:"bizMsg"`
	BizData    string `json:"bizData"` // busId，业务主记录唯一 ID
}

// TradeImport 上传交易订单文件（分账文件上传 Step 2）。
// 以 multipart/form-data 携带 uploadToken 与 XLSX 文件，
// 服务端重新计算文件 SM3 哈希并与凭证中的预期值比对，确保内容未被篡改。
//
// 与标准 JSON 加密流程相互独立，响应为普通 JSON（非加密信封）。
// uploadToken 为一次性凭证，无论上传成功与否都会被消费；上传失败需从 Step 1 重新开始。
// 返回 busId，用于后续确认 / 查询 / 取消。
func (x *Hst) TradeImport(ctx context.Context, dto *TradeImportDto) (busId string, err error) {
	// channel-file 前缀与 BaseURL 的 channel 前缀不同，用拼接 URL
	importURL := x.Option.BaseURL + "-file/doc-trade-file/import"
	var resp *resty.Response
	if resp, err = x.Client.R().SetContext(ctx).
		SetFormData(map[string]string{
			"channelId":   dto.ChannelId,
			"uploadToken": dto.UploadToken,
		}).
		SetMultipartFields(dto.multipartFields()...).
		Post(importURL); err != nil {
		return
	}

	if resp.StatusCode() != 200 {
		err = help.E(0, fmt.Sprintf(`第三方接口响应失败! status=%d body=%s`,
			resp.StatusCode(), resp.String()))
		return
	}

	var importResp *TradeImportResp
	if err = sonic.Unmarshal(resp.Bytes(), &importResp); err != nil {
		return
	}

	if importResp.Code != "SUCCESS" {
		err = help.E(0, fmt.Sprintf(`第三方接口响应失败! code=%s msg=%s`,
			importResp.Code, importResp.Msg))
		return
	}

	if !importResp.Data.BizSuccess {
		err = bizError(importResp.Data.BizCode, importResp.Data.BizMsg)
		return
	}

	busId = importResp.Data.BizData
	return
}
