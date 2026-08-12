package hst

import (
	"context"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
)

type TradeCancelDto struct {
	ReqTimestamp string `json:"reqTimestamp"` // 业务请求时间戳，须与外层信封 timestamp 一致
	PartnerId    string `json:"partnerId"`    // 渠道商 ID
	MerchantId   string `json:"merchantId"`   // 商户 ID
	BusId        string `json:"busId"`        // 批次号 / 主记录 ID
}

func (x *TradeCancelDto) GetTs() string {
	return x.ReqTimestamp
}

// NewTradeCancelDto 创建取消导入请求体。
// partnerId / merchantId / busId 为必填字段
// （reqTimestamp 由 TradeCancel 方法自动填充）。
func NewTradeCancelDto(
	partnerId string,
	merchantId string,
	busId string,
) *TradeCancelDto {
	return &TradeCancelDto{
		PartnerId:  partnerId,
		MerchantId: merchantId,
		BusId:      busId,
	}
}

// TradeCancel 取消导入。
// 取消尚未确认的文档交易导入批次。
// 响应 bizData 为布尔值，true 表示取消成功。
func (x *Hst) TradeCancel(ctx context.Context, dto *TradeCancelDto) (result *SignObjectRespResult[bool], err error) {
	dto.ReqTimestamp = strconv.FormatInt(time.Now().UnixMilli(), 10)

	var signObjectReq *SignObjectReq
	if signObjectReq, err = x.NewSignObjectReq(dto); err != nil {
		return
	}

	var b string
	if b, err = x.Request(ctx,
		"/channel/doc-trade-file/cancel", signObjectReq); err != nil {
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
