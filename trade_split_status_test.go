package hst_test

import (
	"context"
	"testing"

	"github.com/kainonly/hst"
)

// TestTradeSplitStatus 流程 A：查询分账状态（读 trade_import.log 的 busId）
func TestTradeSplitStatus(t *testing.T) {
	ctx := context.Background()

	var busId string
	readLastLogBizData(t, "trade_import", &busId)
	t.Logf("从 trade_import 日志读取到 busId: %s", busId)

	dto := hst.NewTradeSplitStatusDto(busId)
	result, _, err := client.TradeSplitStatus(ctx, dto)
	if err != nil {
		logResult(t, "trade_split_status", errorLogData{false, err.Error()})
		t.Fatalf("TradeSplitStatus 失败: %v", err)
	}
	logResult(t, "trade_split_status", result)
}

// TestTradeSplitStatus2 流程 B：查询分账状态（读 trade_import_2.log 的 busId）
func TestTradeSplitStatus2(t *testing.T) {
	ctx := context.Background()

	var busId string
	readLastLogBizData(t, "trade_import_2", &busId)
	t.Logf("从 trade_import_2 日志读取到 busId: %s", busId)

	dto := hst.NewTradeSplitStatusDto(busId)
	result, _, err := client.TradeSplitStatus(ctx, dto)
	if err != nil {
		logResult(t, "trade_split_status_2", errorLogData{false, err.Error()})
		t.Fatalf("TradeSplitStatus2 失败: %v", err)
	}
	logResult(t, "trade_split_status_2", result)
}
