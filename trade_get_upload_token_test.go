package hst_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/kainonly/hst"
)

// busIdLogData 用于把 trade_import 返回的裸字符串 busId 包装写入日志。
type busIdLogData struct {
	BizSuccess bool   `json:"bizSuccess"`
	BizData    string `json:"bizData"`
}

// TestGetUploadToken 流程 A：用 trade.xlsx 申请上传凭证
func TestGetUploadToken(t *testing.T) {
	ctx := context.Background()
	filePath := filepath.Join("files", "trade.xlsx")
	fileSM3Hash := sm3HashFile(t, filePath)

	dto := hst.NewGetUploadTokenDto(
		cfg.ChannelId,
		cfg.MerchantNo,
		"trade.xlsx",
		fileSM3Hash,
	)
	result, err := client.GetUploadToken(ctx, dto)
	if err != nil {
		logResult(t, "get_upload_token", errorLogData{false, err.Error()})
		t.Fatalf("GetUploadToken 失败: %v", err)
	}
	logResult(t, "get_upload_token", result)
}

// TestGetUploadToken2 流程 B：用 trade-2.xlsx 申请上传凭证
func TestGetUploadToken2(t *testing.T) {
	ctx := context.Background()
	filePath := filepath.Join("files", "trade-2.xlsx")
	fileSM3Hash := sm3HashFile(t, filePath)

	dto := hst.NewGetUploadTokenDto(
		cfg.ChannelId,
		cfg.MerchantNo,
		"trade-2.xlsx",
		fileSM3Hash,
	)
	result, err := client.GetUploadToken(ctx, dto)
	if err != nil {
		logResult(t, "get_upload_token_2", errorLogData{false, err.Error()})
		t.Fatalf("GetUploadToken2 失败: %v", err)
	}
	logResult(t, "get_upload_token_2", result)
}
