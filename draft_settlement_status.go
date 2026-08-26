package hst

import (
	"context"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
)

type SettlementStatusDto struct {
	ReqTimestamp string `json:"reqTimestamp"` // 业务请求时间戳，须与外层信封 timestamp 一致（防重放）
	ChannelNo    string `json:"channelNo"`    // 渠道号
	DraftId      string `json:"draftId"`      // 草稿 ID
}

func (x *SettlementStatusDto) GetTs() string {
	return x.ReqTimestamp
}

func NewSettlementStatusDto(draftId string) *SettlementStatusDto {
	return &SettlementStatusDto{
		DraftId:      draftId,
		ReqTimestamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
	}
}

type SettlementStatusBizData struct {
	DraftId              string `json:"draftId"`              // 草稿 ID
	MerchantId           string `json:"merchantId"`           // 商户号（若草稿已生成商户则有值，否则为空）
	MerchantCreated      bool   `json:"merchantCreated"`      // 是否已生成商户号
	SettlementStatus     string `json:"settlementStatus"`     // 入驻状态码：0-审核中，1-成功，2-失败，3-待激活，4-激活中
	SettlementStatusDesc string `json:"settlementStatusDesc"` // 入驻状态描述
	MybankOrderNo        string `json:"mybankOrderNo"`        // 网商银行申请单号
	MybankMerchantId     string `json:"mybankMerchantId"`     // 网商银行商户编号
	FailReason           string `json:"failReason"`           // 开户失败原因
	ActivateUrl          string `json:"activateUrl"`          // 激活链接（settlementStatus = 3 待激活时使用）
}

func (x *Hst) SettlementStatus(ctx context.Context, dto *SettlementStatusDto) (result *SignObjectRespResult[*SettlementStatusBizData], signObjectResp *SignObjectResp, err error) {
	dto.ChannelNo = x.Option.ChannelId

	var signObjectReq *SignObjectReq
	if signObjectReq, err = x.NewSignObjectReq(dto); err != nil {
		return
	}

	if signObjectResp, err = x.Request(ctx,
		"/channel/merchant_info_draft/settlement_status", signObjectReq); err != nil {
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
