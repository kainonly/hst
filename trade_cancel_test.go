package hst_test

import (
	"context"
	"testing"

	"github.com/kainonly/hst"
)

func TestTradeCancel(t *testing.T) {
	ctx := context.Background()

	// 从 logs/trade_import.log 读取上一次的 busId
	// trade_import 的 bizData 是裸字符串，直接用 string 读取
	var busId string
	readLastLogBizData(t, "trade_import", &busId)
	t.Logf("从 trade_import 日志读取到 busId: %s", busId)

	dto := hst.NewTradeCancelDto(
		cfg.ChannelId,  // partnerId
		cfg.MerchantNo, // merchantId
		busId,          // busId
	)
	result, err := client.TradeCancel(ctx, dto)
	if err != nil {
		logResult(t, "trade_cancel", map[string]any{
			"bizSuccess": false,
			"error":      err.Error(),
		})
		t.Fatalf("TradeCancel 失败: %v", err)
	}

	logResult(t, "trade_cancel", result)
}
