package hst_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/kainonly/hst"
)

func TestGetUploadToken(t *testing.T) {
	ctx := context.Background()

	// 读取真实 XLSX 文件并计算 SM3 哈希
	filePath := filepath.Join("files", "trade.xlsx")
	fileSM3Hash := sm3HashFile(t, filePath)

	dto := hst.NewGetUploadTokenDto(
		cfg.ChannelId,  // partnerId（渠道商 ID）
		cfg.MerchantNo, // merchantId（商户 ID）
		"trade.xlsx",   // fileName
		fileSM3Hash,    // fileSM3Hash（64 位十六进制）
	)
	result, err := client.GetUploadToken(ctx, dto)
	if err != nil {
		logResult(t, "get_upload_token", map[string]any{
			"bizSuccess": false,
			"error":      err.Error(),
		})
		t.Fatalf("GetUploadToken 失败: %v", err)
	}

	logResult(t, "get_upload_token", result)
}
