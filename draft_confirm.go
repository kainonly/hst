package hst

import (
	"context"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
)

type ConfirmDto struct {
	ReqTimestamp string `json:"reqTimestamp"` // 业务请求时间戳，须与外层信封 timestamp 一致（防重放）
	ChannelNo    string `json:"channelNo"`    // 渠道号
	DraftId      string `json:"draftId"`      // 草稿 ID（由创建草稿或上传资质文件返回）
}

func (x *ConfirmDto) GetTs() string {
	return x.ReqTimestamp
}

func NewConfirmDto(draftId string) *ConfirmDto {
	return &ConfirmDto{
		DraftId:      draftId,
		ReqTimestamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
	}
}

// ConfirmBizData 确认提交草稿响应 bizData。
// 确认成功后 draftStatus = CONFIRMED，并回填 merchantId、orgId、accountId；
// 确认失败则 draftStatus = FAILED 并填充 failReason。
// 若需草稿全量视图（含业务字段与文件 URL），使用 UploadFiles 响应的 UploadFilesBizData。
type ConfirmBizData struct {
	DraftId     string `json:"draftId"`     // 草稿唯一 ID
	DraftStatus string `json:"draftStatus"` // 草稿状态：EDITING/SUBMITTING/CONFIRMED/FAILED
	PartnerId   string `json:"partnerId"`   // 渠道商 ID
	MerchantId  string `json:"merchantId"`  // 商户 ID（确认成功后回填）
	OrgId       string `json:"orgId"`       // 企业唯一号（确认成功后回填）
	AccountId   string `json:"accountId"`   // 结算账户 ID（确认成功后回填）
	FailReason  string `json:"failReason"`  // 提交失败原因（draftStatus = FAILED 时填充）
}

// Confirm 确认提交草稿。
// 将两步上传完成后的草稿数据分发写入商户基础信息、商户信息、结算账户三个正式业务表，
// 触发商户入驻申请。确认后草稿状态变更为 CONFIRMED 或 FAILED；
// 若为 FAILED，可通过 UpdatePrepare 修改后重新走两步上传流程再次确认。
func (x *Hst) Confirm(ctx context.Context, dto *ConfirmDto) (result *SignObjectRespResult[*ConfirmBizData], err error) {
	dto.ChannelNo = x.Option.ChannelId

	var signObjectReq *SignObjectReq
	if signObjectReq, err = x.NewSignObjectReq(dto); err != nil {
		return
	}

	var b string
	if b, err = x.Request(ctx,
		"/channel/merchant_info_draft/confirm", signObjectReq); err != nil {
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
