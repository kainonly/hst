package hst_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kainonly/hst"
)

func TestTradeQuery(t *testing.T) {
	ctx := context.Background()

	// 生成外部提现单号（与 apply 测试一致，查不到也无所谓）
	outWithdrawNo := fmt.Sprintf("W%s", time.Now().Format("20060102150405"))

	dto := hst.NewTradeQueryDto(
		cfg.MerchantNo, // merchantNo
		outWithdrawNo,  // outWithdrawNo
	)

	result, err := client.TradeQuery(ctx, dto)
	if err != nil {
		logResult(t, "withdrawal_query", errorLogData{false, err.Error()})
		t.Logf("TradeQuery 返回错误（预期可能失败）: %v", err)
		return
	}

	logResult(t, "withdrawal_query", result)
}
