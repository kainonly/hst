package hst_test

import (
	"context"
	"encoding/hex"
	"os"
	"testing"

	"github.com/kainonly/hst"
	"github.com/tjfoc/gmsm/sm3"
)

// sm3HashFile 对文件原始字节计算 SM3，返回 64 位十六进制字符串。
func sm3HashFile(t *testing.T, path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("读取文件 %s 失败: %v", path, err)
	}
	return hex.EncodeToString(sm3.Sm3Sum(b))
}

func TestCreatePrepare(t *testing.T) {
	ctx := context.Background()

	// fileManifest：对 files 目录下的真实资质图片计算 SM3 哈希
	const filesDir = "files"
	certPhotoAHash := sm3HashFile(t, filepath.Join(filesDir, "sfz-a.jpg"))
	//certPhotoBHash := sm3HashFile(t, filepath.Join(filesDir, "sfz-b.jpg"))
	//licensePhotoHash := sm3HashFile(t, filepath.Join(filesDir, "yyzz.jpg"))

	dto := hst.NewCreatePrepareDto(
		cfg.PersonUsername,
		cfg.PersonUsername,
		[]string{"WICOIN_MYBANK_SPLIT_ACCOUNT"},
		"01", "brand_other", "01", "7299",
		cfg.PersonMobile,
		cfg.PersonUsername,
		"test@example.com",
	).
		SetPrincipal(cfg.PersonMobile, "100", cfg.PersonCertNo, cfg.PersonUsername, "2031-05-20 00:00:00").
		SetLocation("460000", "460100", "460105", "海南省海口市秀英区XX路XX号", cfg.PersonMobile).
		SetPerson("M", "自由职业").
		SetSettlementAccountType("01"). // 结算类型：01 银行卡
		SetBank(cfg.PersonBankCardNo, "01",
			"中国农业银行海口秀英支行", "460000", "460100",
		).
		SetContactLine("310100000000"). // 联行号（非必填）
		SetCardHolder("01", cfg.PersonCertNo, "海南省海口市秀英区XX路XX号")

	//SetLogonId(cfg.PersonMybankAccount). // logonId（网商二类户）
	//SetUserId(cfg.PersonUid).
	//SetFileManifest(&hst.FileManifest{
	//	CertPhotoAFiles:   []string{certPhotoAHash},   // 身份证人像面
	//	CertPhotoBFiles:   []string{certPhotoBHash},   // 身份证国徽面
	//	LicensePhotoFiles: []string{licensePhotoHash}, // 营业执照
	//}).
	//SetAlipayAccount(cfg.PersonMybankAccount). // 支付宝收款账号（网商二类户）
	//SetMerchantType("OTHER") // 商户业务类型：其他

	result, err := client.CreatePrepare(ctx, dto)
	if err != nil {
		logResult(t, "create_prepare", errorLogData{false, err.Error()})
		t.Fatalf("CreatePrepare 失败: %v", err)
	}

	// 写入日志（包含 fileManifest 字段名，供 upload_files 测试读取）
	logResult(t, "create_prepare", prepareLogData{
		UploadToken:   result.BizData.UploadToken,
		ExpireSeconds: result.BizData.ExpireSeconds,
		FileFields:    []string{"certPhotoAFiles", "certPhotoBFiles", "licensePhotoFiles"},
	})
}
