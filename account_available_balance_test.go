package hst_test

import (
	"context"
	"testing"

	"github.com/kainonly/hst"
)

func TestAvailableBalance(t *testing.T) {
	ctx := context.Background()

	dto := hst.NewAvailableBalanceDto(cfg.MerchantNo)
	result, _, err := client.AvailableBalance(ctx, dto)
	if err != nil {
		logResult(t, "account_available_balance", errorLogData{false, err.Error()})
		t.Fatalf("AvailableBalance 失败: %v", err)
	}

	logResult(t, "account_available_balance", result)
}

func TestAvailableBalanceSub(t *testing.T) {
	ctx := context.Background()

	sub := cfg.SubMerchantNo[0]
	dto := hst.NewAvailableBalanceDto(sub)
	result, _, err := client.AvailableBalance(ctx, dto)
	if err != nil {
		logResult(t, "account_available_balance", errorLogData{false, err.Error()})
		t.Fatalf("AvailableBalance 失败: %v", err)
	}

	logResult(t, "account_available_balance", result)
}
