package hst_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kainonly/hst"
)

func TestBrandBalance(t *testing.T) {
	ctx := context.Background()

	outTradeNo := fmt.Sprintf("Q%s", time.Now().Format("20060102150405"))

	dto := hst.NewBrandBalanceDto(
		cfg.ChannelId,  // partnerId（渠道商 ID）
		cfg.MerchantNo, // merchantNo（仅用于定位其所属平台配置）
		outTradeNo,     // outTradeNo（外部请求流水号）
	)
	result, err := client.BrandBalance(ctx, dto)
	if err != nil {
		logResult(t, "account_brand_balance", errorLogData{false, err.Error()})
		t.Fatalf("BrandBalance 失败: %v", err)
	}

	logResult(t, "account_brand_balance", result)
}
