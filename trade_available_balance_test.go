package hst_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kainonly/hst"
)

func TestAvailableBalance(t *testing.T) {
	ctx := context.Background()

	// outTradeNo 仅作请求标识记入日志，不参与查询、不做幂等，但不可为空
	outTradeNo := fmt.Sprintf("Q%s", time.Now().Format("20060102150405"))

	dto := hst.NewAvailableBalanceDto(
		cfg.ChannelId,  // partnerId（渠道商 ID）
		cfg.MerchantNo, // merchantNo（平台商户号）
		outTradeNo,     // outTradeNo（外部请求流水号）
	)
	result, err := client.AvailableBalance(ctx, dto)
	if err != nil {
		logResult(t, "available_balance", map[string]any{
			"bizSuccess": false,
			"error":      err.Error(),
		})
		t.Fatalf("AvailableBalance 失败: %v", err)
	}

	logResult(t, "available_balance", result)
}
