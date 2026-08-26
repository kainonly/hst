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
		cfg.PersonUsername,                      // legalName（自然人即姓名）
		cfg.PersonUsername+"_更新",                // shortName（修改：加"_更新"后缀）
		[]string{"WICOIN_MYBANK_SPLIT_ACCOUNT"}, // productCode
		"01",                                    // merchantBaseType：01 自然人
		"brand_other",                           // subRoleType：其他服务公司
		"01",                                    // dealType：实体特约商户
		"7299",                                  // mcc：其他生活服务
		cfg.PersonMobile,
		cfg.PersonUsername,
		"test@example.com",
	).
		SetPrincipal(cfg.PersonMobile, "100", cfg.PersonCertNo, cfg.PersonUsername, "2031-05-20 00:00:00").
		SetLocation("460000", "460100", "460105", "海南省海口市秀英区更新路XX号", cfg.PersonMobile). // 修改：地址变更
		SetPerson("M", "自由职业").
		SetSettlementAccountType("01"). // 结算类型：01 银行卡
		SetBank(cfg.PersonBankCardNo, "01",
			"中国农业银行海口秀英支行", "460000", "460100",
		).
		SetContactLine("310100000000"). // 联行号（非必填）
		SetCardHolder("01", cfg.PersonCertNo, "海南省海口市秀英区更新路XX号").
		SetLogonId(cfg.PersonMybankAccount). // logonId（网商二类户）
		SetUserId(cfg.PersonUid).
		SetFileManifest(&hst.FileManifest{
			ShopPhotoFiles: []string{shopPhotoHash}, // 门头照（增量更新）
		}).
		SetAlipayAccount(cfg.PersonMybankAccount). // 支付宝收款账号（网商二类户）
		SetMerchantType("OTHER")                   // 商户业务类型：其他

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
