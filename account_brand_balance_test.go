package hst_test

import (
	"context"
	"testing"

	"github.com/kainonly/hst"
)

func TestBrandBalance(t *testing.T) {
	ctx := context.Background()

	dto := hst.NewBrandBalanceDto(cfg.MerchantNo)
	result, err := client.BrandBalance(ctx, dto)
	if err != nil {
		logResult(t, "account_brand_balance", errorLogData{false, err.Error()})
		t.Fatalf("BrandBalance 失败: %v", err)
	}

	logResult(t, "account_brand_balance", result)
}

func TestBrandBalanceSub(t *testing.T) {
	ctx := context.Background()

	sub := cfg.SubMerchantNo[0]
	dto := hst.NewBrandBalanceDto(sub)
	result, err := client.BrandBalance(ctx, dto)
	if err != nil {
		logResult(t, "account_brand_balance", errorLogData{false, err.Error()})
		t.Fatalf("BrandBalance 失败: %v", err)
	}

	logResult(t, "account_brand_balance", result)
}
