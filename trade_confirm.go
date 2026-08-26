package hst

import (
	"context"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
)

type TradeConfirmDto struct {
	ReqTimestamp string `json:"reqTimestamp"` // 业务请求时间戳，须与外层信封 timestamp 一致
	PartnerId    string `json:"partnerId"`    // 渠道商 ID
	MerchantId   string `json:"merchantId"`   // 商户 ID
	BusId        string `json:"busId"`        // 批次号 / 主记录 ID（由上传接口返回）
}

func (x *TradeConfirmDto) GetTs() string {
	return x.ReqTimestamp
}

// NewTradeConfirmDto 创建确认导入请求体。
func NewTradeConfirmDto(busId string) *TradeConfirmDto {
	return &TradeConfirmDto{
		ReqTimestamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
		BusId:        busId,
	}
}

// TradeConfirm 确认导入。
// 确认已上传的文档交易批次并触发补单分账，需在文件解析完成后调用。
// 文件仍在解析中（docStatus = IMPORTING 或明细数未就绪）时不允许确认。
// 响应 bizData 为布尔值，true 表示确认成功。
func (x *Hst) TradeConfirm(ctx context.Context, dto *TradeConfirmDto) (result *SignObjectRespResult[bool], signObjectResp *SignObjectResp, err error) {
	dto.PartnerId = x.Option.ChannelId
	dto.MerchantId = x.Option.MerchantNo

	var signObjectReq *SignObjectReq
	if signObjectReq, err = x.NewSignObjectReq(dto); err != nil {
		return
	}

	if signObjectResp, err = x.Request(ctx,
		"/channel/doc-trade-file/confirm", signObjectReq); err != nil {
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
