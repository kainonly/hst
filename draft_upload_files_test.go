package hst_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/kainonly/hst"
)

// fileFieldMap 资质字段名 -> 本地文件路径
var fileFieldMap = map[string]string{
	"certPhotoAFiles":   filepath.Join("files", "sfz-a.jpg"), // 身份证人像面
	"certPhotoBFiles":   filepath.Join("files", "sfz-b.jpg"), // 身份证国徽面
	"licensePhotoFiles": filepath.Join("files", "yyzz.jpg"),  // 营业执照
	"shopPhotoFiles":    filepath.Join("files", "shop.jpg"),  // 门头照（如有）
}

// prepareLogData prepare 接口（create/update）记录到日志的结构。
// 包含 uploadToken 和本次声明的 fileManifest 字段名列表，
// 供 upload_files 测试读取，决定上传哪些文件。
type prepareLogData struct {
	UploadToken   string   `json:"uploadToken"`
	ExpireSeconds int64    `json:"expireSeconds"`
	FileFields    []string `json:"fileFields"` // 本次 fileManifest 中声明的字段名
}

func TestUploadFiles(t *testing.T) {
	ctx := context.Background()

	// 优先从 logs/update_prepare.log 读取（更新流程），不存在则从 logs/create_prepare.log 读取（创建流程）
	var sourceLog = "create_prepare"
	if _, err := os.Stat(filepath.Join("logs", "update_prepare.log")); err == nil {
		sourceLog = "update_prepare"
	}
	var prepareResult prepareLogData
	readLastLogBizData(t, sourceLog, &prepareResult)
	t.Logf("从 %s 日志读取到 uploadToken: %s, fileFields: %v",
		sourceLog, prepareResult.UploadToken, prepareResult.FileFields)

	// 构造上传 DTO，只上传 prepare 的 fileManifest 中声明且本地有对应文件的字段
	dto := hst.NewUploadFilesDto(prepareResult.UploadToken)
	for _, fieldName := range prepareResult.FileFields {
		filePath, ok := fileFieldMap[fieldName]
		if !ok {
			t.Logf("字段 %s 无本地文件映射，跳过", fieldName)
			continue
		}
		if _, err := os.Stat(filePath); err != nil {
			t.Logf("本地文件 %s 不存在，跳过字段 %s", filePath, fieldName)
			continue
		}
		switch fieldName {
		case "certPhotoAFiles":
			dto.SetCertPhotoAFiles(filePath)
		case "certPhotoBFiles":
			dto.SetCertPhotoBFiles(filePath)
		case "licensePhotoFiles":
			dto.SetLicensePhotoFiles(filePath)
		case "shopPhotoFiles":
			dto.SetShopPhotoFiles(filePath)
		}
	}

	result, err := client.UploadFiles(ctx, dto)
	if err != nil {
		logResult(t, "upload_files", errorLogData{false, err.Error()})
		t.Fatalf("UploadFiles 失败: %v", err)
	}

	logResult(t, "upload_files", result)
}
