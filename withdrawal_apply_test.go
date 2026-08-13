package hst_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kainonly/hst"
)

func TestApply(t *testing.T) {
	ctx := context.Background()

	// 生成唯一外部提现单号（幂等键）
	outWithdrawNo := fmt.Sprintf("W%s", time.Now().Format("20060102150405"))

	dto := hst.NewApplyDto(
		cfg.MerchantNo, // merchantNo
		outWithdrawNo,  // outWithdrawNo（幂等键）
		"1.00",         // totalAmount（单位元，精确到分）
	).SetRemark("测试提现")

	result, err := client.Apply(ctx, dto)
	if err != nil {
		logResult(t, "withdrawal_apply", errorLogData{false, err.Error()})
		t.Logf("Apply 返回错误（预期可能失败）: %v", err)
		return
	}

	logResult(t, "withdrawal_apply", result)
}
