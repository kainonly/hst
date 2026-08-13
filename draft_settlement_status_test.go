package hst_test

import (
	"context"
	"testing"

	"github.com/kainonly/hst"
)

func TestSettlementStatus(t *testing.T) {
	ctx := context.Background()

	// 从 logs/confirm.log 读取上一次的 draftId
	var confirmResult struct {
		DraftId     string `json:"draftId"`
		DraftStatus string `json:"draftStatus"`
		MerchantId  string `json:"merchantId"`
	}
	readLastLogBizData(t, "confirm", &confirmResult)
	t.Logf("从 confirm 日志读取到 draftId: %s", confirmResult.DraftId)

	dto := hst.NewSettlementStatusDto(confirmResult.DraftId)
	result, err := client.SettlementStatus(ctx, dto)
	if err != nil {
		logResult(t, "settlement_status", errorLogData{false, err.Error()})
		t.Fatalf("SettlementStatus 失败: %v", err)
	}

	logResult(t, "settlement_status", result)
}
