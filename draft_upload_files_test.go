package hst_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/kainonly/hst"
)

func TestUploadFiles(t *testing.T) {
	ctx := context.Background()

	// 从 logs/create_prepare.log 读取上一次的 uploadToken
	var createPrepareResult struct {
		UploadToken   string `json:"uploadToken"`
		ExpireSeconds int64  `json:"expireSeconds"`
	}
	readLastLogBizData(t, "create_prepare", &createPrepareResult)
	t.Logf("从 create_prepare 日志读取到 uploadToken: %s", createPrepareResult.UploadToken)

	// 上传 files/ 目录下的真实资质图片
	// 字段须与 create_prepare 的 fileManifest 一致
	dto := hst.NewUploadFilesDto(
		cfg.ChannelId,
		createPrepareResult.UploadToken,
	).SetCertPhotoAFiles(filepath.Join("files", "sfz-a.jpg")). // 身份证人像面
									SetCertPhotoBFiles(filepath.Join("files", "sfz-b.jpg")). // 身份证国徽面
									SetLicensePhotoFiles(filepath.Join("files", "yyzz.jpg")) // 营业执照

	result, err := client.UploadFiles(ctx, dto)
	if err != nil {
		// 把错误也写入日志，便于排查
		logResult(t, "upload_files", map[string]any{
			"bizSuccess": false,
			"error":      err.Error(),
		})
		t.Fatalf("UploadFiles 失败: %v", err)
	}

	logResult(t, "upload_files", result)
}
