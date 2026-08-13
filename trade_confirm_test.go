package hst_test

import (
	"context"
	"testing"

	"github.com/kainonly/hst"
)

func TestTradeConfirm(t *testing.T) {
	ctx := context.Background()

	// 从 logs/trade_import.log 读取上一次的 busId
	// trade_import 的 bizData 是裸字符串，直接用 string 读取
	var busId string
	readLastLogBizData(t, "trade_import", &busId)
	t.Logf("从 trade_import 日志读取到 busId: %s", busId)

	dto := hst.NewTradeConfirmDto(
		cfg.ChannelId,  // partnerId
		cfg.MerchantNo, // merchantId
		busId,          // busId
	)
	result, err := client.TradeConfirm(ctx, dto)
	if err != nil {
		logResult(t, "trade_confirm", map[string]any{
			"bizSuccess": false,
			"error":      err.Error(),
		})
		t.Fatalf("TradeConfirm 失败: %v", err)
	}

	logResult(t, "trade_confirm", result)
}
