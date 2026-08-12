package hst_test

import (
	"context"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/kainonly/hst"
	"github.com/stretchr/testify/assert"
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
	certPhotoAHash := sm3HashFile(t, filepath.Join(filesDir, "sfz-a.jpg"))  // 身份证人像面
	certPhotoBHash := sm3HashFile(t, filepath.Join(filesDir, "sfz-b.jpg"))  // 身份证国徽面
	licensePhotoHash := sm3HashFile(t, filepath.Join(filesDir, "yyzz.jpg")) // 营业执照

	dto := hst.NewCreatePrepareDto(
		[]string{"WICOIN_PAY"},  // productCode
		cfg.PersonUsername,      // legalName（自然人即姓名）
		cfg.PersonUsername,      // shortName
		"01",                    // merchantBaseType：01 自然人
		"brand_other",           // subRoleType：其他服务公司
		"01",                    // dealType：实体特约商户
		"7299",                  // mcc：其他生活服务
		cfg.PersonMobile,        // contactMobile
		cfg.PersonUsername,      // contactName
		"test@example.com",      // email
		cfg.PersonMobile,        // principalMobile（法人手机号）
		"100",                   // principalCertType：100 身份证
		cfg.PersonCertNo,        // principalCertNo
		cfg.PersonUsername,      // principalPerson
		"2031-05-20 00:00:00",   // principalCertVld（证件有效期）
		"460000",                // province：海南省
		"460100",                // city：海口市
		"460105",                // district：秀英区
		"海南省海口市秀英区XX路XX号",       // address
		cfg.PersonMobile,        // servicePhoneNo
		"M",                     // personSex：M 男性（证件号第17位 1）
		"自由职业",                  // personProfession
		"01",                    // settlementAccountType：01 银行卡
		cfg.PersonBankCardNo,    // bankCardNo
		cfg.PersonUsername,      // bankCertName
		"01",                    // accountType：01 对私账户（自然人）
		"310100000000",          // contactLine：联行号
		"中国农业银行海口秀英支行",          // branchName
		"460000",                // branchProvince
		"460100",                // branchCity
		"01",                    // certType：01 身份证
		cfg.PersonCertNo,        // certNo（持卡人证件号码）
		"海南省海口市秀英区XX路XX号",       // cardHolderAddress
		cfg.PersonMybankAccount, // logonId（支付宝登陆账号，用网商二类户）
		cfg.PersonUid,           // userId
		&hst.FileManifest{
			CertPhotoAFiles:   []string{certPhotoAHash},   // 身份证人像面
			CertPhotoBFiles:   []string{certPhotoBHash},   // 身份证国徽面
			LicensePhotoFiles: []string{licensePhotoHash}, // 营业执照
		},
	).SetAlipayAccount(cfg.PersonMybankAccount) // 支付宝收款账号（网商二类户）

	result, err := client.CreatePrepare(ctx, dto)
	assert.NoError(t, err)

	logResult(t, "create_prepare", result)
}
