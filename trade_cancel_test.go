package hst_test

import (
	"context"
	"testing"

	"github.com/kainonly/hst"
)

// TestTradeCancel 流程 B：取消导入（读 trade_import_2.log 的 busId）
func TestTradeCancel(t *testing.T) {
	ctx := context.Background()

	var busId string
	readLastLogBizData(t, "trade_import_2", &busId)
	t.Logf("从 trade_import_2 日志读取到 busId: %s", busId)

	dto := hst.NewTradeCancelDto(cfg.ChannelId, cfg.MerchantNo, busId)
	result, err := client.TradeCancel(ctx, dto)
	if err != nil {
		logResult(t, "trade_cancel", errorLogData{false, err.Error()})
		t.Fatalf("TradeCancel 失败: %v", err)
	}
	logResult(t, "trade_cancel", result)
}
