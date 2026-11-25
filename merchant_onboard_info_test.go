package hst_test

import (
	"context"
	"testing"

	"github.com/kainonly/hst"
)

func TestMerchantOnboardInfo(t *testing.T) {
	ctx := context.Background()

	dto := hst.NewMerchantOnboardInfoDto(cfg.SubMerchantNo[0])

	result, _, err := client.MerchantOnboardInfo(ctx, dto)
	if err != nil {
		logResult(t, "merchant_onboard_info", errorLogData{false, err.Error()})
		return
	}

	logResult(t, "merchant_onboard_info", result)
}
