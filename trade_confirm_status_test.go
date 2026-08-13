package hst_test

import (
	"context"
	"testing"

	"github.com/kainonly/hst"
)

// TestTradeConfirm 流程 A：确认导入（读 trade_import.log 的 busId）
func TestTradeConfirm(t *testing.T) {
	ctx := context.Background()

	var busId string
	readLastLogBizData(t, "trade_import", &busId)
	t.Logf("从 trade_import 日志读取到 busId: %s", busId)

	dto := hst.NewTradeConfirmDto(cfg.ChannelId, cfg.MerchantNo, busId)
	result, err := client.TradeConfirm(ctx, dto)
	if err != nil {
		logResult(t, "trade_confirm", errorLogData{false, err.Error()})
		t.Fatalf("TradeConfirm 失败: %v", err)
	}
	logResult(t, "trade_confirm", result)
}

// TestTradeStatus 流程 A：查询主记录状态（读 trade_import.log 的 busId）
func TestTradeStatus(t *testing.T) {
	ctx := context.Background()

	var busId string
	readLastLogBizData(t, "trade_import", &busId)
	t.Logf("从 trade_import 日志读取到 busId: %s", busId)

	dto := hst.NewTradeStatusDto(cfg.ChannelId, cfg.MerchantNo, busId)
	result, err := client.TradeStatus(ctx, dto)
	if err != nil {
		logResult(t, "trade_status", errorLogData{false, err.Error()})
		t.Fatalf("TradeStatus 失败: %v", err)
	}
	logResult(t, "trade_status", result)
}
