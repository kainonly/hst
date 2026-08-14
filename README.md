# hst

WiCoin 渠道商 Go SDK，封装了进件、分账、余额查询、提现等全部渠道接口，内置 SM2 签名 + SM4 加密的标准安全流程。

## 特性

- **国密安全**：内置 SM2 签名/验签、SM4-CBC 加密/解密、SM3 文件哈希
- **JSON 高性能**：使用 bytedance/sonic 进行 JSON 编解码
- **链式调用**：非必填字段通过 `Set` 方法链式设置
- **泛型响应**：利用 Go 1.26 泛型，统一 `SignObjectRespResult[T]` 包装业务响应
- **接口全覆盖**：进件 5 个 + 分账 5 个 + 余额查询 2 个 + 提现 2 个，共 14 个接口

## 安装

```bash
go get github.com/kainonly/hst
```

## 快速开始

```go
package main

import (
    "context"
    "fmt"
    "github.com/kainonly/hst"
)

func main() {
    client, err := hst.NewHst(&hst.Option{
        BaseURL:    "https://test-cqpay.oldbird.tech/api/v1",
        PriKey:     "<your_pri_key_hex>",
        WicoPubKey: "<your_wico_pub_key_hex>",
        EncryptKey: "<your_encrypt_key_hex>",
        MerchantNo: "<your_merchant_no>",
        ChannelId:  "<your_channel_id>",
    })
    if err != nil {
        panic(err)
    }

    // 进件：创建草稿并申请上传凭证
    result, err := client.CreatePrepare(context.Background(), hst.NewCreatePrepareDto(
        []string{"WICOIN_PAY"},
        "上海某某科技有限公司",
        "某某科技",
        // ... 其余必填参数
        &hst.FileManifest{
            CertPhotoAFiles: []string{"<sm3_hash>"},
        },
    ).SetMerchantType("OTHER"))
    if err != nil {
        panic(err)
    }
    fmt.Println(result.BizData.UploadToken)
}
```

## 初始化

```go
client, err := hst.NewHst(&hst.Option{
    BaseURL:    "https://test-cqpay.oldbird.tech/api/v1",
    PriKey:     "<SM2 私钥 Hex>",
    WicoPubKey: "<支付中心 SM2 公钥 Hex>",
    EncryptKey: "<SM4 加密密钥 Hex>",
    MerchantNo: "<商户号>",
    ChannelId:  "<渠道商 ID>",
})
```

| 字段 | 说明 |
|---|---|
| `BaseURL` | 网关地址，到 `/api/v1` 为止 |
| `PriKey` | 客户端 SM2 私钥（Hex 或 Base64） |
| `WicoPubKey` | 支付中心 SM2 公钥（Hex 或 Base64） |
| `EncryptKey` | SM4 加密密钥（Hex 或 Base64） |
| `MerchantNo` | 平台商户号 |
| `ChannelId` | 渠道商 ID（即 `partnerId` / `channelNo`） |

> SDK 内部会自动填充 `ReqTimestamp`、`PartnerId`、`MerchantId`、`ChannelNo` 等信封字段，无需手动设置。

## 进件接口

商户入驻相关接口，采用两步上传流程（先提交业务字段 + 文件哈希清单换取凭证，再上传实际文件）。

| 方法 | 路径 | 说明 |
|---|---|---|
| `CreatePrepare` | `/channel/merchant_info_draft/create/prepare` | 创建草稿并申请上传凭证（Step 1） |
| `UploadFiles` | `/channel-multi-file/merchant_info_draft/upload_files` | 上传资质文件（Step 2，multipart） |
| `UpdatePrepare` | `/channel/merchant_info_draft/update/prepare` | 更新草稿并申请上传凭证 |
| `Confirm` | `/channel/merchant_info_draft/confirm` | 确认提交草稿 |
| `SettlementStatus` | `/channel/merchant_info_draft/settlement_status` | 查询入驻状态 |

### 创建草稿并申请上传凭证

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

### 上传资质文件

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

### 更新草稿并申请上传凭证

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

### 确认提交草稿

将草稿数据写入正式业务表，触发商户入驻申请。

```go
dto := hst.NewConfirmDto("<draft_id>")
result, err := client.Confirm(ctx, dto)
// result.BizData.DraftStatus — CONFIRMED 或 FAILED
// result.BizData.MerchantId  — 确认成功后回填
// result.BizData.OrgId       — 企业唯一号
// result.BizData.AccountId   — 结算账户 ID
```

### 查询入驻状态

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

## 分账接口

文档交易订单文件导入相关接口，采用两步上传流程（先申请凭证，再上传文件）。

| 方法 | 路径 | 说明 |
|---|---|---|
| `GetUploadToken` | `/channel/doc-trade-file/getUploadToken` | 申请文件上传凭证（Step 1） |
| `TradeImport` | `/channel-file/doc-trade-file/import` | 上传交易订单文件（Step 2，multipart） |
| `TradeConfirm` | `/channel/doc-trade-file/confirm` | 确认导入 |
| `TradeStatus` | `/channel/doc-trade-file/status` | 查询主记录状态 |
| `TradeCancel` | `/channel/doc-trade-file/cancel` | 取消导入 |

### 申请文件上传凭证

两步上传 **Step 1**。预先计算文件 SM3 哈希并纳入签名体。

```go
// 计算文件 SM3 哈希
fileBytes, _ := os.ReadFile("trade.xlsx")
fileSM3Hash := hex.EncodeToString(sm3.Sm3Sum(fileBytes))

dto := hst.NewGetUploadTokenDto(
    "trade.xlsx",      // fileName 文件名
    fileSM3Hash,       // fileSM3Hash 64位十六进制
)
result, err := client.GetUploadToken(ctx, dto)
// result.BizData.UploadToken   — 上传凭证
// result.BizData.ExpireSeconds — 有效期（秒，默认 300）
```

### 上传交易订单文件

两步上传 **Step 2**。`multipart/form-data` 上传 XLSX 文件。

```go
dto := hst.NewTradeImportDto(
    "PC2025xxxx",      // channelId
    "<upload_token>",  // uploadToken
    "files/trade.xlsx", // filePath XLSX 文件路径
)
busId, err := client.TradeImport(ctx, dto)
// busId — 业务主记录唯一 ID，用于后续确认/查询/取消
```

> `uploadToken` 一次性消费，上传失败需从 Step 1 重新申请。

### 确认导入

确认已上传的文档交易批次并触发补单分账。

```go
dto := hst.NewTradeConfirmDto("<bus_id>")
result, err := client.TradeConfirm(ctx, dto)
// result.BizData — bool，true 表示确认成功
```

### 查询主记录状态

```go
dto := hst.NewTradeStatusDto("<bus_id>")
result, err := client.TradeStatus(ctx, dto)
// result.BizData.DocStatus        — IMPORTING/PENDING/SUCCESS/FAILED
// result.BizData.TotalDetailCount — 明细总数
// result.BizData.SuccessCount     — 成功数
// result.BizData.FailCount        — 失败数
// result.BizData.SuccessAmount    — 成功金额（元）
// result.BizData.FailAmount       — 失败金额（元）
// result.BizData.ProcessingAmount — 处理中金额（元）
```

### 取消导入

取消尚未确认的文档交易导入批次。

```go
dto := hst.NewTradeCancelDto("<bus_id>")
result, err := client.TradeCancel(ctx, dto)
// result.BizData — bool，true 表示取消成功
```

## 余额查询接口

查询商户资金账户余额，用于结算对账与提现前的额度判断。

| 方法 | 路径 | 说明 |
|---|---|---|
| `AvailableBalance` | `/channel/merchant_account/available_balance` | 查询商户可用余额 |
| `BrandBalance` | `/channel/merchant_account/brand-balance` | 查询品牌商户专户余额 |

### 查询商户可用余额

```go
dto := hst.NewAvailableBalanceDto("<merchant_no>")
result, err := client.AvailableBalance(ctx, dto)
// result.BizData.BalanceInfos — 余额明细列表
// 按accountType取值，不要按下标取：
//   AVAILABLE_BALANCE — 可用余额（可提现）
//   PENDING_BALANCE  — 待结算金额（不可提现）
```

> `PENDING_BALANCE` 是尚未解冻的待结算金额，不可提现。把两者相加当作可提现额度会导致提现失败。

### 查询品牌商户专户余额

```go
dto := hst.NewBrandBalanceDto("<merchant_no>")
result, err := client.BrandBalance(ctx, dto)
// result.BizData — 字符串，品牌专户余额（单位元）
```

> 此余额是平台侧备付金，同一平台配置下不同商户号查到的是同一个余额。不能作为商户可用/可提现额度。

## 提现接口

渠道商代商户发起提现，把可用余额提现到商户绑定的结算银行卡。

| 方法 | 路径 | 说明 |
|---|---|---|
| `Apply` | `/channel/merchant/withdrawal/apply` | 商户提现申请 |
| `TradeQuery` | `/channel/merchant/withdrawal/query` | 查询提现订单 |

### 商户提现申请

```go
dto := hst.NewApplyDto(
    "<merchant_no>",     // merchantNo
    "W20260806-0001",   // outWithdrawNo 幂等键
    "100.00",           // totalAmount 提现金额（元）
).SetRemark("日常结算提现")

result, err := client.Apply(ctx, dto)
// result.BizData.WithdrawNo — 平台提现单号
// result.BizData.Status    — 提现状态
```

> `outWithdrawNo` 是幂等键，重复调用返回首次结果。超时/中断时须用原单号查询，不可换单号重发。

### 查询提现订单

```go
dto := hst.NewTradeQueryDto(
    "<merchant_no>",
    "W20260806-0001",  // outWithdrawNo 申请时的单号
)
result, err := client.TradeQuery(ctx, dto)
// result.BizData.Status             — 提现状态
// result.BizData.WithdrawFinishDate — 完成时间
```

| status | 说明 |
|---|---|
| `DEALING` | 处理中 |
| `WAIT_CONFIRM` | 待确认 |
| `SUCCESS` | 成功（资金已到卡） |
| `FAIL` | 失败 |
| `UNKNOWN` | 未知（联系平台） |

## 错误处理

SDK 统一两层错误格式：

| 层级 | 格式 | 说明 |
|---|---|---|
| 网关层 | `第三方接口响应失败! code=%s msg=%s` | HTTP 非 200 或网关 Code 非 SUCCESS |
| 业务层 | `第三方接口业务失败! bizCode=%s bizMsg=%s` | `bizSuccess` 为 false |

所有错误通过 `help.E(code, msg)` 封装，`code` 为 `int64`（当前固定 `0`），`msg` 包含完整的错误码和描述。

```go
result, err := client.CreatePrepare(ctx, dto)
if err != nil {
    // err 包含 bizCode 和 bizMsg
    log.Printf("进件失败: %v", err)
    return
}
```

## 数据类型

### FileManifest

文件哈希清单，用于进件接口。每个字段对应一种资质，值为按上传顺序排列的 SM3 哈希列表（64 位 Hex）。

```go
type FileManifest struct {
    CertPhotoAFiles             []string `json:"certPhotoAFiles,omitempty"`             // 身份证人像面
    CertPhotoBFiles             []string `json:"certPhotoBFiles,omitempty"`             // 身份证国徽面
    LicensePhotoFiles           []string `json:"licensePhotoFiles,omitempty"`           // 营业执照
    PrgPhotoFiles               []string `json:"prgPhotoFiles,omitempty"`               // 组织机构代码证
    IndustryLicensePhotoFiles   []string `json:"industryLicensePhotoFiles,omitempty"`   // 开户许可证
    ShopPhotoFiles              []string `json:"shopPhotoFiles,omitempty"`              // 门头照
    OtherPhotoFiles             []string `json:"otherPhotoFiles,omitempty"`             // 其他资料
    CertPhotoCFiles             []string `json:"certPhotoCFiles,omitempty"`             // 手持身份证
    RegisterProtocolPhotoFiles  []string `json:"registerProtocolPhotoFiles,omitempty"`  // 商户入驻协议
    ContractPhotoFiles          []string `json:"contractPhotoFiles,omitempty"`          // 租赁协议
    ShopEntrancePhotoFiles      []string `json:"shopEntrancePhotoFiles,omitempty"`      // 门店内景
    CheckstandPhotoFiles        []string `json:"checkstandPhotoFiles,omitempty"`        // 收银台
    MerchantAgreementPhotoFiles []string `json:"merchantAgreementPhotoFiles,omitempty"` // 商户协议
}
```

### SignObjectRespResult

泛型业务响应包装。

```go
type SignObjectRespResult[T any] struct {
    BizSuccess bool   `json:"bizSuccess"`
    BizCode    string `json:"bizCode"`
    BizMsg     string `json:"bizMsg"`
    BizData    T      `json:"bizData"`
}
```

### SM3 文件哈希计算

```go
import (
    "encoding/hex"
    "github.com/tjfoc/gmsm/sm3"
)

func sm3HashFile(path string) string {
    b, _ := os.ReadFile(path)
    return hex.EncodeToString(sm3.Sm3Sum(b))
}
```

## 项目结构

```
hst/
├── hst.go                         核心类型 + 加解密函数
├── draft_create_prepare.go        进件：创建草稿
├── draft_update_prepare.go        进件：更新草稿
├── draft_upload_files.go          进件：上传资质文件
├── draft_confirm.go              进件：确认提交
├── draft_settlement_status.go     进件：查询入驻状态
├── trade_get_upload_token.go      分账：申请上传凭证
├── trade_import.go               分账：上传交易订单
├── trade_confirm.go              分账：确认导入
├── trade_status.go               分账：查询状态
├── trade_cancel.go               分账：取消导入
├── account_available_balance.go   余额：查询可用余额
├── account_brand_balance.go       余额：查询品牌专户余额
├── withdrawal_apply.go            提现：提现申请
├── withdrawal_query.go            提现：查询提现订单
├── config/                        配置文件
│   ├── values.yml                 实际配置（gitignore）
│   └── values.example.yml         配置模板
├── files/                         测试用资质文件（gitignore）
├── logs/                          测试日志（gitignore）
└── *_test.go                      单元测试
```

## 配置

参考 `config/values.example.yml`：

```yaml
context: dev
configs:
  dev:
    base_url: https://test-cqpay.oldbird.tech/api/v1
    pri_key: <your_pri_key_hex>
    wico_pub_key: <your_wico_pub_key_hex>
    encrypt_key: <your_encrypt_key_hex>
    merchant_no: <your_merchant_no>
    channel_id: <your_channel_id>
```

## License

MIT
