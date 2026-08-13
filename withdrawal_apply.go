package hst

import (
	"context"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
)

type ApplyDto struct {
	ReqTimestamp  string `json:"reqTimestamp"`  // 业务请求时间戳，须与外层信封 timestamp 一致（防重放）
	PartnerId     string `json:"partnerId"`     // 渠道商 ID（与信封 channelId 不是同一个字段，必须显式传入）
	MerchantNo    string `json:"merchantNo"`    // 平台商户号，须归属于 partnerId
	OutWithdrawNo string `json:"outWithdrawNo"` // 外部提现单号，渠道商系统唯一（幂等键，重复调用返回首次结果）
	TotalAmount   string `json:"totalAmount"`   // 提现金额，单位元，精确到分（如 100.00），必须大于 0
	Remark        string `json:"remark"`        // 备注
}

func (x *ApplyDto) GetTs() string {
	return x.ReqTimestamp
}

// NewApplyDto 创建商户提现申请请求体。
func NewApplyDto(
	merchantNo string,
	outWithdrawNo string,
	totalAmount string,
) *ApplyDto {
	return &ApplyDto{
		MerchantNo:    merchantNo,
		OutWithdrawNo: outWithdrawNo,
		TotalAmount:   totalAmount,
		ReqTimestamp:  strconv.FormatInt(time.Now().UnixMilli(), 10),
	}
}

// SetRemark 设置备注。
func (x *ApplyDto) SetRemark(i string) *ApplyDto {
	x.Remark = i
	return x
}

// ApplyBizData 商户提现申请响应 bizData。
//
// 首次受理只回带 withdrawNo / status / errorDesc；
// 命中幂等（同一 outWithdrawNo 重复提交）时回带库中订单的完整快照
// （含 outWithdrawNo / totalAmount / withdrawApplyDate / withdrawFinishDate）。
//
// 不要依赖申请接口回带金额或时间做对账，一律以查询提现订单的返回为准。
type ApplyBizData struct {
	WithdrawNo         string `json:"withdrawNo"`         // 平台提现单号，唯一
	OutWithdrawNo      string `json:"outWithdrawNo"`      // 外部提现单号（原样回带，仅幂等命中时返回）
	TotalAmount        string `json:"totalAmount"`        // 提现金额，单位元（仅幂等命中时返回）
	Status             string `json:"status"`             // 提现状态
	WithdrawApplyDate  string `json:"withdrawApplyDate"`  // 提现申请时间（仅幂等命中时返回）
	WithdrawFinishDate string `json:"withdrawFinishDate"` // 提现完成时间（仅幂等命中时返回）
	ErrorDesc          string `json:"errorDesc"`          // 错误描述，status=FAIL 时有值
}

// Apply 商户提现申请。
// 渠道商代其名下商户发起提现，把网商银行账户中的可用余额提现到商户已绑定的结算银行卡。
// 平台内部合并了网商银行的「申请 + 短信确认」两步，渠道商无需处理短信验证码。
//
// 危险操作提示：
//   - 调用超时、连接中断时资金可能已在出账。此时须用原 outWithdrawNo 调用查询提现订单查明结果；
//     换一个单号重发等于再提现一笔。
//   - outWithdrawNo 是幂等键，重复调用返回首次结果。
func (x *Hst) Apply(ctx context.Context, dto *ApplyDto) (result *SignObjectRespResult[*ApplyBizData], err error) {
	dto.PartnerId = x.Option.ChannelId
	dto.MerchantNo = x.Option.MerchantNo

	var signObjectReq *SignObjectReq
	if signObjectReq, err = x.NewSignObjectReq(dto); err != nil {
		return
	}

	var b string
	if b, err = x.Request(ctx,
		"/channel/merchant/withdrawal/apply", signObjectReq); err != nil {
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
