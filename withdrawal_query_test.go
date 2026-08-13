package hst_test

import (
	"context"
	"testing"

	"github.com/kainonly/hst"
)

func TestTradeQuery(t *testing.T) {
	ctx := context.Background()

	// 生成外部提现单号（与 apply 测试一致，查不到也无所谓）

	dto := hst.NewTradeQueryDto(
		cfg.SubMerchantNo[0],
		`W20260813161712`, // outWithdrawNo
	)

	result, err := client.TradeQuery(ctx, dto)
	if err != nil {
		logResult(t, "withdrawal_query", errorLogData{false, err.Error()})
		return
	}

	logResult(t, "withdrawal_query", result)
}
