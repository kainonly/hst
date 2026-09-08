package hst_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/kainonly/hst"
)

// TestTradeImport 流程 A：上传 trade.xlsx
func TestTradeImport(t *testing.T) {
	ctx := context.Background()

	var getTokenResult struct {
		UploadToken   string `json:"uploadToken"`
		ExpireSeconds int64  `json:"expireSeconds"`
	}
	readLastLogBizData(t, "get_upload_token", &getTokenResult)
	t.Logf("从 get_upload_token 日志读取到 uploadToken: %s", getTokenResult.UploadToken)

	filePath := filepath.Join("files", "trade.xlsx")
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("读取交易文件失败: %v", err)
	}
	dto := hst.NewTradeImportDto(getTokenResult.UploadToken, hst.NewUploadFile(filepath.Base(filePath), fileData))
	busId, err := client.TradeImport(ctx, dto)
	if err != nil {
		logResult(t, "trade_import", errorLogData{false, err.Error()})
		t.Fatalf("TradeImport 失败: %v", err)
	}
	logResult(t, "trade_import", busIdLogData{true, busId})
}

// TestTradeImport2 流程 B：上传 trade-2.xlsx
func TestTradeImport2(t *testing.T) {
	ctx := context.Background()

	var getTokenResult struct {
		UploadToken   string `json:"uploadToken"`
		ExpireSeconds int64  `json:"expireSeconds"`
	}
	readLastLogBizData(t, "get_upload_token_2", &getTokenResult)
	t.Logf("从 get_upload_token_2 日志读取到 uploadToken: %s", getTokenResult.UploadToken)

	filePath := filepath.Join("files", "trade-2.xlsx")
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("读取交易文件失败: %v", err)
	}
	dto := hst.NewTradeImportDto(getTokenResult.UploadToken, hst.NewUploadFile(filepath.Base(filePath), fileData))
	busId, err := client.TradeImport(ctx, dto)
	if err != nil {
		logResult(t, "trade_import_2", errorLogData{false, err.Error()})
		t.Fatalf("TradeImport2 失败: %v", err)
	}
	logResult(t, "trade_import_2", busIdLogData{true, busId})
}
