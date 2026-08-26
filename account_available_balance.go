package hst

import (
	"context"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/kainonly/go/help"
)

type AvailableBalanceDto struct {
	ReqTimestamp string `json:"reqTimestamp"` // 业务请求时间戳，须与外层信封 timestamp 一致（防重放）
	PartnerId    string `json:"partnerId"`    // 渠道商 ID（与信封 channelId 不是同一个字段，必须显式传入）
	MerchantNo   string `json:"merchantNo"`   // 平台商户号，须归属于 partnerId
	OutTradeNo   string `json:"outTradeNo"`   // 外部请求流水号（仅作请求标识记入日志，不参与查询、不做幂等）
}

func (x *AvailableBalanceDto) GetTs() string {
	return x.ReqTimestamp
}

// NewAvailableBalanceDto 创建查询商户可用余额请求体。
func NewAvailableBalanceDto(merchantNo string) *AvailableBalanceDto {
	return &AvailableBalanceDto{
		ReqTimestamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
		OutTradeNo:   help.SID(),
		MerchantNo:   merchantNo,
	}
}

// BalanceInfo 余额明细。
type BalanceInfo struct {
	AccountType string `json:"accountType"` // 金额类型：AVAILABLE_BALANCE 可用余额 / PENDING_BALANCE 待结算金额
	TotalAmount string `json:"totalAmount"` // 金额，单位元
	Currency    string `json:"currency"`    // 币种，默认 CNY（人民币）
}

// AvailableBalanceBizData 查询商户可用余额响应 bizData。
//
// balanceInfos 按实际存在的账户返回，不是固定两条。
// 网商侧四类账户中，平台只透出：
//   - 提现账户 -> AVAILABLE_BALANCE（可提现）
//   - 余额冻结账户 -> PENDING_BALANCE（待结算，不可提现）
//
// 某类账户不存在时该条目直接缺席，请按 accountType 取值，不要按下标取。
// 查无账户时为空数组。
type AvailableBalanceBizData struct {
	BalanceInfos []BalanceInfo `json:"balanceInfos"` // 余额明细列表
}

// AvailableBalance 查询商户可用余额。
// 查询商户在网商银行资金账户下的可用余额与待结算金额，
// 用于结算对账与提现前的额度判断。
//
// 注意：PENDING_BALANCE 是尚未解冻的待结算金额，不可提现。
// 把 AVAILABLE_BALANCE 与 PENDING_BALANCE 相加当作可提现额度会导致提现申请以余额不足失败。
func (x *Hst) AvailableBalance(ctx context.Context, dto *AvailableBalanceDto) (result *SignObjectRespResult[*AvailableBalanceBizData], signObjectResp *SignObjectResp, err error) {
	dto.PartnerId = x.Option.ChannelId

	var signObjectReq *SignObjectReq
	if signObjectReq, err = x.NewSignObjectReq(dto); err != nil {
		return
	}

	if signObjectResp, err = x.Request(ctx,
		"/channel/merchant_account/available_balance", signObjectReq); err != nil {
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
