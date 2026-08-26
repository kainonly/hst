package hst

import (
	"context"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
)

type TradeQueryDto struct {
	ReqTimestamp  string `json:"reqTimestamp"`  // 业务请求时间戳，须与外层信封 timestamp 一致（防重放）
	PartnerId     string `json:"partnerId"`     // 渠道商 ID（与信封 channelId 不是同一个字段，必须显式传入）
	MerchantNo    string `json:"merchantNo"`    // 平台商户号
	OutWithdrawNo string `json:"outWithdrawNo"` // 外部提现单号，即申请时传入的单号
}

func (x *TradeQueryDto) GetTs() string {
	return x.ReqTimestamp
}

// NewTradeQueryDto 创建查询提现订单请求体。
func NewTradeQueryDto(
	merchantNo string,
	outWithdrawNo string,
) *TradeQueryDto {
	return &TradeQueryDto{
		MerchantNo:    merchantNo,
		OutWithdrawNo: outWithdrawNo,
		ReqTimestamp:  strconv.FormatInt(time.Now().UnixMilli(), 10),
	}
}

// TradeQueryBizData 查询提现订单响应 bizData。
// 查询返回的是订单完整快照（与申请接口首次受理时不同）。
type TradeQueryBizData struct {
	WithdrawNo         string `json:"withdrawNo"`         // 平台提现单号，唯一
	OutWithdrawNo      string `json:"outWithdrawNo"`      // 外部提现单号（原样回带）
	TotalAmount        string `json:"totalAmount"`        // 提现金额，单位元
	Status             string `json:"status"`             // 提现状态
	WithdrawApplyDate  string `json:"withdrawApplyDate"`  // 提现申请时间
	WithdrawFinishDate string `json:"withdrawFinishDate"` // 提现完成时间，未完成时为空
	ErrorDesc          string `json:"errorDesc"`          // 错误描述，status=FAIL 时有值
}

// TradeQuery 查询提现订单。
// 按外部提现单号查询提现订单的最新状态，用于商户提现申请之后的结果确认与对账。
//
// 订单须同时匹配 outWithdrawNo 与 partnerId 才能查到；查不到时返回「提现订单不存在」。
// 刚提交后立即查询到该错误，通常意味着申请请求根本没有落库，应重新提交（用原单号）。
//
// 状态说明：
//   - 只有 SUCCESS 才代表资金已到卡
//   - DEALING / WAIT_CONFIRM 都是中间态，须继续轮询
//   - UNKNOWN 须联系平台核实，不能自行判为失败
func (x *Hst) TradeQuery(ctx context.Context, dto *TradeQueryDto) (result *SignObjectRespResult[*TradeQueryBizData], signObjectResp *SignObjectResp, err error) {
	dto.PartnerId = x.Option.ChannelId

	var signObjectReq *SignObjectReq
	if signObjectReq, err = x.NewSignObjectReq(dto); err != nil {
		return
	}

	if signObjectResp, err = x.Request(ctx,
		"/channel/merchant/withdrawal/query", signObjectReq); err != nil {
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
