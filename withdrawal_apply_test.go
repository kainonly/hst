package hst_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kainonly/hst"
)

func TestWithdrawalApply(t *testing.T) {
	ctx := context.Background()

	// 生成唯一外部提现单号（幂等键）
	outWithdrawNo := fmt.Sprintf("W%s", time.Now().Format("20060102150405"))

	dto := hst.NewWithdrawalApplyDto(
		cfg.SubMerchantNo[0], // merchantNo
		outWithdrawNo,        // outWithdrawNo（幂等键）
		"1.00",               // totalAmount（单位元，精确到分）
	).SetRemark("测试提现")

	result, _, err := client.WithdrawalApply(ctx, dto)
	if err != nil {
		logResult(t, "withdrawal_apply", errorLogData{false, err.Error()})
		return
	}

	logResult(t, "withdrawal_apply", result)
}
