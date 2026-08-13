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

	dto := hst.NewAvailableBalanceDto(outTradeNo)
	result, err := client.AvailableBalance(ctx, dto)
	if err != nil {
		logResult(t, "account_available_balance", errorLogData{false, err.Error()})
		t.Fatalf("AvailableBalance 失败: %v", err)
	}

	logResult(t, "account_available_balance", result)
}
