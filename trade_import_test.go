package hst_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/kainonly/hst"
)

func TestTradeImport(t *testing.T) {
	ctx := context.Background()

	// 从 logs/get_upload_token.log 读取上一次的 uploadToken
	var getTokenResult struct {
		UploadToken   string `json:"uploadToken"`
		ExpireSeconds int64  `json:"expireSeconds"`
	}
	readLastLogBizData(t, "get_upload_token", &getTokenResult)
	t.Logf("从 get_upload_token 日志读取到 uploadToken: %s", getTokenResult.UploadToken)

	// 上传 files/trade.xlsx
	filePath := filepath.Join("files", "trade.xlsx")
	dto := hst.NewTradeImportDto(
		cfg.ChannelId,              // channelId（渠道商 ID）
		getTokenResult.UploadToken, // uploadToken
		filePath,                   // XLSX 文件本地路径
	)

	busId, err := client.TradeImport(ctx, dto)
	if err != nil {
		logResult(t, "trade_import", errorLogData{false, err.Error()})
		t.Fatalf("TradeImport 失败: %v", err)
	}

	// busId 是裸字符串，用结构体包装写入日志（保持字段顺序一致）
	logResult(t, "trade_import", struct {
		BizSuccess bool   `json:"bizSuccess"`
		BizData    string `json:"bizData"`
	}{
		BizSuccess: true,
		BizData:    busId,
	})
}
