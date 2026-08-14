# 进件接口

商户入驻相关接口，采用两步上传流程（先提交业务字段 + 文件哈希清单换取凭证，再上传实际文件）。

## 接口列表

| 方法 | 路径 | 说明 |
|---|---|---|
| `CreatePrepare` | `/channel/merchant_info_draft/create/prepare` | 创建草稿并申请上传凭证（Step 1） |
| `UploadFiles` | `/channel-multi-file/merchant_info_draft/upload_files` | 上传资质文件（Step 2，multipart） |
| `UpdatePrepare` | `/channel/merchant_info_draft/update/prepare` | 更新草稿并申请上传凭证 |
| `Confirm` | `/channel/merchant_info_draft/confirm` | 确认提交草稿 |
| `SettlementStatus` | `/channel/merchant_info_draft/settlement_status` | 查询入驻状态 |

---

## 创建草稿并申请上传凭证

两步上传 **Step 1**。提交商户入驻业务字段与资质文件 SM3 哈希清单，换取一次性上传凭证。

```go
dto := hst.NewCreatePrepareDto(
    []string{"WICOIN_PAY"},           // productCode 产品编码列表
    "上海某某科技有限公司",               // legalName 法定名称
    "某某科技",                        // shortName 商户简称
    "03",                             // merchantBaseType 01自然人/02个体/03企业
    "brand_other",                    // subRoleType 商户角色
    "01",                             // dealType 经营类型
    "7299",                           // mcc 经营类目
    "13800000000",                    // contactMobile 联系人手机号
    "张三",                            // contactName 联系人姓名
    "zhangsan@example.com",           // email 邮箱
    "13800000001",                    // principalMobile 负责人手机号
    "100",                            // principalCertType 100身份证
    "310***********1234",             // principalCertNo 负责人证件号
    "李四",                            // principalPerson 负责人姓名
    "2035-01-01 00:00:00",            // principalCertVld 证件有效期
    "310000",                         // province 省
    "310100",                         // city 市
    "310104",                         // district 区
    "上海市徐汇区XX路XX号",              // address 详细地址
    "021-00000000",                   // servicePhoneNo 客服电话
    "M",                              // personSex 性别 M/F
    "企业法人",                         // personProfession 职业
    "01",                             // settlementAccountType 01银行卡
    "62220000000000000",              // bankCardNo 银行卡号
    "上海某某科技有限公司",               // bankCertName 银行账户户名
    "02",                             // accountType 01对私/02对公
    "310100000000",                   // contactLine 联行号
    "某银行上海徐汇支行",                 // branchName 开户支行
    "310000",                         // branchProvince 开户省
    "310100",                         // branchCity 开户市
    "01",                             // certType 01身份证
    "310***********1234",              // certNo 持卡人证件号
    "上海市徐汇区XX路XX号",              // cardHolderAddress 持卡人地址
    "zhangsan@example.com",           // logonId 支付宝登录号
    "2088***********",                // userId 支付宝用户ID
    &hst.FileManifest{
        CertPhotoAFiles:   []string{"<sm3_hash>"},  // 身份证人像面
        CertPhotoBFiles:   []string{"<sm3_hash>"},  // 身份证国徽面
        LicensePhotoFiles: []string{"<sm3_hash>"},  // 营业执照
    },
).SetMerchantType("OTHER")             // 可选：商户业务类型
  .SetAlipayAccount("2088xxx")          // 可选：支付宝收款账号
  .SetRemark("备注")                    // 可选：备注

result, err := client.CreatePrepare(ctx, dto)
// result.BizData.UploadToken   — 上传凭证（UUID）
// result.BizData.ExpireSeconds — 凭证有效期（秒，固定 900）
```

---

## 上传资质文件

两步上传 **Step 2**。以 `multipart/form-data` 携带凭证与文件上传，服务端按字段名 + 顺序重算 SM3 比对。

```go
dto := hst.NewUploadFilesDto(
    "PC2025xxxx",     // channelId 渠道商 ID
    "<upload_token>",  // uploadToken 来自 CreatePrepare/UpdatePrepare
).SetCertPhotoAFiles("files/sfz-a.jpg").  // 身份证人像面
  SetCertPhotoBFiles("files/sfz-b.jpg").  // 身份证国徽面
  SetLicensePhotoFiles("files/yyzz.jpg")  // 营业执照

bizData, err := client.UploadFiles(ctx, dto)
// bizData.DraftId     — 草稿 ID
// bizData.DraftStatus — 草稿状态 EDITING/SUBMITTING/CONFIRMED/FAILED
```

> `uploadToken` 为一次性凭证，无论成功与否都会被消费。上传字段须与 Step 1 的 `FileManifest` 一致。

---

## 更新草稿并申请上传凭证

仅 `EDITING` 或 `FAILED` 状态的草稿可更新。业务字段须整体重新提交，`fileManifest` 只需包含新文件字段。

```go
dto := hst.NewUpdatePrepareDto(
    "<draft_id>",      // draftId 草稿 ID
    []string{"WICOIN_PAY"},
    // ... 其余必填字段与 CreatePrepare 相同
    &hst.FileManifest{
        ShopPhotoFiles: []string{"<sm3_hash>"},  // 仅更新门头照
    },
).SetMerchantType("OTHER")

result, err := client.UpdatePrepare(ctx, dto)
// result.BizData.UploadToken — 新的上传凭证
```

---

## 确认提交草稿

将草稿数据写入正式业务表，触发商户入驻申请。

```go
dto := hst.NewConfirmDto("<draft_id>")
result, err := client.Confirm(ctx, dto)
// result.BizData.DraftStatus — CONFIRMED 或 FAILED
// result.BizData.MerchantId  — 确认成功后回填
// result.BizData.OrgId       — 企业唯一号
// result.BizData.AccountId   — 结算账户 ID
```

---

## 查询入驻状态

```go
dto := hst.NewSettlementStatusDto("<draft_id>")
result, err := client.SettlementStatus(ctx, dto)
// result.BizData.SettlementStatus — 入驻状态码
// result.BizData.MerchantCreated  — 是否已生成商户号
// result.BizData.ActivateUrl      — 待激活时的激活链接
```

| settlementStatus | 说明 |
|---|---|
| `0` | 审核中 |
| `1` | 成功 |
| `2` | 失败 |
| `3` | 待激活 |
| `4` | 激活中 |
