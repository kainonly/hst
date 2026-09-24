package hst

import (
	"context"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
)

type TradeSplitStatusDto struct {
	ReqTimestamp string `json:"reqTimestamp"` // 业务请求时间戳，须与外层信封 timestamp 一致
	PartnerId    string `json:"partnerId"`    // 渠道商 ID
	MerchantId   string `json:"merchantId"`   // 商户 ID
	BusId        string `json:"busId"`        // 文档导入批次号 / 主记录 ID（由上传接口返回）
}

func (x *TradeSplitStatusDto) GetTs() string {
	return x.ReqTimestamp
}

// NewTradeSplitStatusDto 创建查询分账状态请求体。
func NewTradeSplitStatusDto(busId string) *TradeSplitStatusDto {
	return &TradeSplitStatusDto{
		BusId:        busId,
		ReqTimestamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
	}
}

// TradeSplitStatusBizData 查询分账状态响应 bizData。
type TradeSplitStatusBizData struct {
	BusId             string `json:"busId"`             // 分账记录唯一标识
	DocTradeFileBusId string `json:"docTradeFileBusId"` // 文档导入批次号
	PartnerId         string `json:"partnerId"`         // 渠道商 ID
	MerchantId        string `json:"merchantId"`        // 商户 ID
	DetailCount       int64  `json:"detailCount"`       // 本次选中的补单成功明细数
	TotalAmount       string `json:"totalAmount"`       // 本次分账汇总金额，单位：元，保留两位小数
	MybankBatchNo     string `json:"mybankBatchNo"`     // 网商银行来款打批批次号
	SplitStatus       string `json:"splitStatus"`       // 分账状态，取值见方法注释
	FailReason        string `json:"failReason"`        // 失败/余额不足原因，仅 FAILED/BALANCE_INSUFFICIENT 时有值
	CreateTime        string `json:"createTime"`        // 创建时间
	UpdateTime        string `json:"updateTime"`        // 更新时间
}

// TradeSplitStatus 查询分账状态。
// 按批次号查询最新一条分账记录的状态与进度。
// 补单完成（docStatus = COMPLETED / PARTIAL_SUCCESS）后平台自动触发一次分账，
// 可用本接口跟踪分账进度。
//
// 批次从未触发过分账时 bizData 为 null（BizData 为 nil），
// 此时应继续以 TradeStatus 轮询 docStatus。
//
// splitStatus 取值：
//   - PENDING             已受理，待异步执行（非终态）
//   - PROCESSING          分账处理中（非终态）
//   - SUCCESS             分账成功（终态），只有 SUCCESS 才代表分账已完成
//   - FAILED              分账失败（终态），需平台运营人工介入，禁止自行重试
//   - BALANCE_INSUFFICIENT 订单管理专户余额不足（终态），需平台运营人工介入，禁止自行重试
//   - NO_ELIGIBLE_RECORDS 批次下无可分账的成功补单记录（终态）
func (x *Hst) TradeSplitStatus(ctx context.Context, dto *TradeSplitStatusDto) (result *SignObjectRespResult[*TradeSplitStatusBizData], signObjectResp *SignObjectResp, err error) {
	dto.PartnerId = x.Option.ChannelId
	dto.MerchantId = x.Option.MerchantNo

	var signObjectReq *SignObjectReq
	if signObjectReq, err = x.NewSignObjectReq(dto); err != nil {
		return
	}

	if signObjectResp, err = x.Request(ctx,
		"/channel/doc-trade-file/split/status", signObjectReq); err != nil {
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
