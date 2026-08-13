package hst_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/kainonly/hst"
)

func TestUpdatePrepare(t *testing.T) {
	ctx := context.Background()

	// 从 logs/upload_files.log 读取上一次的 draftId
	var uploadFilesResult struct {
		DraftId string `json:"draftId"`
	}
	readLastLogBizData(t, "upload_files", &uploadFilesResult)
	t.Logf("从 upload_files 日志读取到 draftId: %s", uploadFilesResult.DraftId)

	// fileManifest：本次只更新门头照（演示增量上传）
	// 对 files 目录下真实文件计算 SM3
	shopPhotoPath := filepath.Join("files", "shop.jpg")
	shopPhotoHash := sm3HashFile(t, shopPhotoPath)

	dto := hst.NewUpdatePrepareDto(
		uploadFilesResult.DraftId,               // draftId（定位已有草稿）
		[]string{"WICOIN_MYBANK_SPLIT_ACCOUNT"}, // productCode
		cfg.PersonUsername,                      // legalName（自然人即姓名）
		cfg.PersonUsername+"_更新",                // shortName（修改：加"_更新"后缀）
		"01",                                    // merchantBaseType：01 自然人
		"brand_other",                           // subRoleType：其他服务公司
		"01",                                    // dealType：实体特约商户
		"7299",                                  // mcc：其他生活服务
		cfg.PersonMobile,                        // contactMobile
		cfg.PersonUsername,                      // contactName
		"test@example.com",                      // email
		cfg.PersonMobile,                        // principalMobile（法人手机号）
		"100",                                   // principalCertType：100 身份证
		cfg.PersonCertNo,                        // principalCertNo
		cfg.PersonUsername,                      // principalPerson
		"2031-05-20 00:00:00",                   // principalCertVld（证件有效期）
		"460000",                                // province：海南省
		"460100",                                // city：海口市
		"460105",                                // district：秀英区
		"海南省海口市秀英区更新路XX号", // address（修改：地址变更）
		cfg.PersonMobile,     // servicePhoneNo
		"M",                  // personSex：M 男性
		"自由职业",               // personProfession
		"01",                 // settlementAccountType：01 银行卡
		cfg.PersonBankCardNo, // bankCardNo
		cfg.PersonUsername,   // bankCertName
		"01",                 // accountType：01 对私账户（自然人）
		"310100000000",       // contactLine：联行号
		"中国农业银行海口秀英支行",       // branchName
		"460000",             // branchProvince
		"460100",             // branchCity
		"01",                 // certType：01 身份证
		cfg.PersonCertNo,     // certNo（持卡人证件号码）
		"海南省海口市秀英区更新路XX号",       // cardHolderAddress
		cfg.PersonMybankAccount, // logonId
		cfg.PersonUid,           // userId
		&hst.FileManifest{
			ShopPhotoFiles: []string{shopPhotoHash}, // 门头照（增量更新）
		},
	).SetAlipayAccount(cfg.PersonMybankAccount).
		SetMerchantType("OTHER")

	result, err := client.UpdatePrepare(ctx, dto)
	if err != nil {
		logResult(t, "update_prepare", errorLogData{false, err.Error()})
		t.Fatalf("UpdatePrepare 失败: %v", err)
	}

	// 写入日志（包含 fileManifest 字段名，供 upload_files 测试读取）
	logResult(t, "update_prepare", prepareLogData{
		UploadToken:   result.BizData.UploadToken,
		ExpireSeconds: result.BizData.ExpireSeconds,
		FileFields:    []string{"shopPhotoFiles"},
	})
}
