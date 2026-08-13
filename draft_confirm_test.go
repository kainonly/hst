package hst_test

import (
	"context"
	"testing"

	"github.com/kainonly/hst"
)

func TestConfirm(t *testing.T) {
	ctx := context.Background()

	// 从 logs/upload_files.log 读取上一次的 draftId
	var uploadFilesResult struct {
		DraftId     string `json:"draftId"`
		DraftStatus string `json:"draftStatus"`
	}
	readLastLogBizData(t, "upload_files", &uploadFilesResult)
	t.Logf("从 upload_files 日志读取到 draftId: %s", uploadFilesResult.DraftId)

	dto := hst.NewConfirmDto(uploadFilesResult.DraftId)
	result, err := client.Confirm(ctx, dto)
	if err != nil {
		logResult(t, "confirm", errorLogData{false, err.Error()})
		t.Fatalf("Confirm 失败: %v", err)
	}

	logResult(t, "confirm", result)
}
