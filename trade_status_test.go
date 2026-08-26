package hst_test

import (
	"context"
	"testing"

	"github.com/kainonly/hst"
)

// TestTradeStatus 流程 A：查询主记录状态（读 trade_import.log 的 busId）
func TestTradeStatus(t *testing.T) {
	ctx := context.Background()

	var busId string
	readLastLogBizData(t, "trade_import", &busId)
	t.Logf("从 trade_import 日志读取到 busId: %s", busId)

	dto := hst.NewTradeStatusDto(busId)
	result, _, err := client.TradeStatus(ctx, dto)
	if err != nil {
		logResult(t, "trade_status", errorLogData{false, err.Error()})
		t.Fatalf("TradeStatus 失败: %v", err)
	}
	logResult(t, "trade_status", result)
}

// TestTradeStatus2 流程 B：查询主记录状态（读 trade_import_2.log 的 busId）
func TestTradeStatus2(t *testing.T) {
	ctx := context.Background()

	var busId string
	readLastLogBizData(t, "trade_import_2", &busId)
	t.Logf("从 trade_import_2 日志读取到 busId: %s", busId)

	dto := hst.NewTradeStatusDto(busId)
	result, _, err := client.TradeStatus(ctx, dto)
	if err != nil {
		logResult(t, "trade_status_2", errorLogData{false, err.Error()})
		t.Fatalf("TradeStatus2 失败: %v", err)
	}
	logResult(t, "trade_status_2", result)
}
