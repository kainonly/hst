package hst

import (
	"context"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
)

type TradeStatusDto struct {
	ReqTimestamp string `json:"reqTimestamp"` // 业务请求时间戳，须与外层信封 timestamp 一致
	PartnerId    string `json:"partnerId"`    // 渠道商 ID
	MerchantId   string `json:"merchantId"`   // 商户 ID
	BusId        string `json:"busId"`        // 批次号 / 主记录 ID
}

func (x *TradeStatusDto) GetTs() string {
	return x.ReqTimestamp
}

// NewTradeStatusDto 创建查询主记录状态请求体。
// partnerId / merchantId / busId 为必填字段
// （reqTimestamp 由 TradeStatus 方法自动填充）。
func NewTradeStatusDto(
	partnerId string,
	merchantId string,
	busId string,
) *TradeStatusDto {
	return &TradeStatusDto{
		PartnerId:  partnerId,
		MerchantId: merchantId,
		BusId:      busId,
	}
}

// TradeStatusBizData 查询主记录状态响应 bizData。
type TradeStatusBizData struct {
	BusId            string `json:"busId"`            // 主记录唯一标识
	PartnerId        string `json:"partnerId"`        // 渠道商 ID
	MerchantId       string `json:"merchantId"`       // 商户 ID
	FileName         string `json:"fileName"`         // 文件名
	FileSize         int64  `json:"fileSize"`         // 文件大小
	TotalDetailCount int64  `json:"totalDetailCount"` // 导入明细总条数
	SuccessCount     int64  `json:"successCount"`     // 成功记录数
	FailCount        int64  `json:"failCount"`        // 失败记录数（补单失败）
	SuccessAmount    string `json:"successAmount"`    // 成功总金额，单位：元，保留两位小数
	FailAmount       string `json:"failAmount"`       // 失败总金额，单位：元，保留两位小数
	ProcessingAmount string `json:"processingAmount"` // 处理中金额，单位：元，保留两位小数
	DocStatus        string `json:"docStatus"`        // 处理状态（如 IMPORTING 解析中）
	CreateTime       string `json:"createTime"`       // 创建时间
	UpdateTime       string `json:"updateTime"`       // 更新时间
}

// TradeStatus 查询主记录状态。
// 根据批次号查询文档交易导入主记录的状态与进度。
//
// 金额口径说明：
//   - 三项金额按明细状态实时汇总：successAmount 对应处理成功的明细，
//     failAmount 对应补单失败的明细，processingAmount 对应待确认与处理中的明细。
//   - 已取消（CANCELLED）的明细不计入以上任何一项，因此三项之和不一定等于导入文件的总金额。
//   - 确认补单（/confirm）之前，明细均处于待确认状态，此时 processingAmount 即为本批次的待确认总金额。
//   - 批次下无明细时返回 "0.00"，不会返回 null。
func (x *Hst) TradeStatus(ctx context.Context, dto *TradeStatusDto) (result *SignObjectRespResult[*TradeStatusBizData], err error) {
	dto.ReqTimestamp = strconv.FormatInt(time.Now().UnixMilli(), 10)

	var signObjectReq *SignObjectReq
	if signObjectReq, err = x.NewSignObjectReq(dto); err != nil {
		return
	}

	var b string
	if b, err = x.Request(ctx,
		"/channel/doc-trade-file/status", signObjectReq); err != nil {
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
