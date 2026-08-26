package hst

import (
	"context"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/kainonly/go/help"
)

type BrandBalanceDto struct {
	ReqTimestamp string `json:"reqTimestamp"` // 业务请求时间戳，须与外层信封 timestamp 一致（防重放）
	PartnerId    string `json:"partnerId"`    // 渠道商 ID（与信封 channelId 不是同一个字段，必须显式传入）
	MerchantNo   string `json:"merchantNo"`   // 平台商户号，须归属于 partnerId（仅用于定位其所属的平台配置，查的不是该商户自己的账户）
	OutTradeNo   string `json:"outTradeNo"`   // 外部请求流水号（仅作请求标识记入日志，不参与查询、不做幂等）
}

func (x *BrandBalanceDto) GetTs() string {
	return x.ReqTimestamp
}

// NewBrandBalanceDto 创建查询品牌商户订单管理专户余额请求体。
func NewBrandBalanceDto(merchantNo string) *BrandBalanceDto {
	return &BrandBalanceDto{
		ReqTimestamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
		MerchantNo:   merchantNo,
		OutTradeNo:   help.SID(),
	}
}

// BrandBalance 查询品牌商户订单管理专户余额。
// 查询商户所属平台（ISV）配置下的品牌订单管理专户余额，用于平台侧备付金监控。
//
// 注意：
//   - 此余额是平台侧的备付金，同一平台配置下的不同商户号查到的是同一个余额，
//     不能作为该商户可用/可提现额度使用 —— 商户自己的余额请用 AvailableBalance。
//   - 查询失败即抛错（MYBANK_ACCOUNT_BALANCE_QUERY_FAILED），不返回 0 或空串；
//     收到成功响应就意味着这是一个真实余额，不必再判空。
//
// 响应 bizData 为裸字符串（品牌商户订单管理专户余额，单位元）。
func (x *Hst) BrandBalance(ctx context.Context, dto *BrandBalanceDto) (result *SignObjectRespResult[string], signObjectResp *SignObjectResp, err error) {
	dto.PartnerId = x.Option.ChannelId

	var signObjectReq *SignObjectReq
	if signObjectReq, err = x.NewSignObjectReq(dto); err != nil {
		return
	}

	if signObjectResp, err = x.Request(ctx,
		"/channel/merchant_account/brand-balance", signObjectReq); err != nil {
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
