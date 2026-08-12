package hst_test

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/kainonly/hst"
	"github.com/stretchr/testify/assert"
	"github.com/tjfoc/gmsm/sm3"
)

func TestGetUploadToken(t *testing.T) {
	ctx := context.Background()

	// 构造测试用的文件内容并计算 SM3 哈希
	// 注意：真实 XLSX 是二进制格式，此处仅用 CSV 文本演示 SM3 计算；
	// 实际使用时需先生成合规的 XLSX 文件再计算哈希。
	fileContent := []byte("交易金额(元),订单摘要,收款商户ID\n100.00,6月直播佣金,MC2025xxxx\n")
	fileSM3Hash := hex.EncodeToString(sm3.Sm3Sum(fileContent))

	dto := hst.NewGetUploadTokenDto(
		cfg.ChannelId,  // partnerId（此处复用 channelId 作为渠道商 ID，按实际配置调整）
		cfg.MerchantNo, // merchantId
		"trade-test.xlsx",
		fileSM3Hash,
	)
	result, err := client.GetUploadToken(ctx, dto)
	assert.NoError(t, err)

	logResult(t, "get_upload_token", result)
}
