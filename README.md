# hst

WiCoin 渠道商 Go SDK，封装了进件、分账、余额查询、提现等全部渠道接口，内置 SM2 签名 + SM4 加密的标准安全流程。

> **AI 阅读指引**：本文档按「机制 → 流程 → API 参考 → 数据类型」组织。
> 理解本 SDK 需先读「核心机制」与「标识符生命周期」两节；调用任何接口前须核对「业务流程」中的顺序约束与「API 参考」中的函数签名（文档中所有代码示例的参数顺序均与真实签名一致）。
> 所有接口分两类：**加密信封类**（JSON 请求，SM2 签名 + SM4 加密，见核心机制）与 **multipart 类**（`UploadFiles` / `TradeImport`，普通 HTTP 上传，无信封）。

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
    result, _, err := client.CreatePrepare(context.Background(), hst.NewCreatePrepareDto(
        "上海某某科技有限公司",  // legalName 法定名称
        "某某科技",             // shortName 商户简称
        []string{"WICOIN_PAY"}, // productCode 产品编码列表
        "03",                   // merchantBaseType 01自然人/02个体/03企业
        "brand_other",          // subRoleType 商户角色
        "01",                   // dealType 经营类型 01实体/02网络/03兼有
        "7299",                 // mcc 经营类目
        "13800000000",          // contactMobile 联系人手机号
        "张三",                  // contactName 联系人姓名
        "zhangsan@example.com", // email 邮箱
    ).SetPrincipal( // 可选组：负责人信息
        "13800000001", "100", "310***********1234", "李四", "2035-01-01 00:00:00",
    ).SetFileManifest(&hst.FileManifest{ // 文件哈希清单（两步上传 Step 1）
        CertPhotoAFiles:   []string{"<sm3_hash>"}, // 身份证人像面
        LicensePhotoFiles: []string{"<sm3_hash>"}, // 营业执照
    }))
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

## 核心机制

### 加密信封（除 multipart 外的全部接口）

SDK 与网关之间的 JSON 请求/响应均包裹在加密信封中，流程固定：

```
请求（SDK → 网关）：
  业务 DTO ──sonic.Marshal──> 明文 JSON body
    ├─ SM2 私钥签名（userId = ChannelId）        → signature（Base64）
    └─ SM4-CBC 加密（key=EncryptKey，随机 IV）    → body（Base64）
  组装 SignObjectReq {timestamp, reqTxn, ivHex, channelId, version:"1.0", signature, body}
  POST <BaseURL><path>（application/json）

响应（网关 → SDK）：
  SignObjectResp {code, msg, timestamp, clientReqTxn, serverAppId, ivHex, signature, body}
    ├─ code != "SUCCESS"            → 网关层错误（HTTP 非 200 同理）
    ├─ SM4-CBC 解密 body（iv=ivHex）→ 业务明文 JSON
    ├─ SM2 公钥验签（userId = serverAppId）→ 失败即错误
    └─ 解析为 SignObjectRespResult[T] {bizSuccess, bizCode, bizMsg, bizData}
         └─ bizSuccess == false     → 业务层错误
```

- 信封对调用方完全透明：业务方法返回的已是解密验签后的 `SignObjectRespResult[T]` 或 `BizData`
- 如需原始信封（审计），见「获取网关响应信封」
- `(*Hst).NewSignObjectReq`（组信封）与 `(*Hst).Request`（发请求 + 解密验签）为公开方法，仅特殊场景直接使用

### multipart 上传（`UploadFiles` / `TradeImport`）

不走加密信封。`multipart/form-data` 直接携带 `channelId`、`uploadToken` 与文件流；
服务端按 Step 1 声明的 SM3 哈希逐字段 + 按顺序重算比对，不一致即拒绝。

### 两步上传模式

进件与分账的文件传输均为两步：

1. **Step 1（加密信封）**：提交业务字段 + 文件 SM3 哈希清单（进件为 `FileManifest`，分账为 `fileName + fileSM3Hash`），换取一次性 `uploadToken`
2. **Step 2（multipart）**：携带 `uploadToken` 上传真实文件，服务端重算 SM3 与 Step 1 比对

`uploadToken` 一次性消费（无论成败），失败须从 Step 1 重新开始。

## 标识符生命周期

| 标识符 | 产生位置 | 消费/使用位置 | 说明 |
|---|---|---|---|
| `channelId` / `partnerId` / `channelNo` | 初始化配置 `Option.ChannelId` | SDK 自动填充到请求 | 渠道商身份，三者同义 |
| `uploadToken` | `CreatePrepare` / `UpdatePrepare` / `GetUploadToken` 响应 | `UploadFiles` / `TradeImport` | 一次性，进件 900s / 分账 300s 有效 |
| `draftId` | `UploadFiles` 响应 `BizData.DraftId` | `UpdatePrepare` / `Confirm` / `SettlementStatus` | 草稿唯一 ID，确认前有效 |
| `merchantId` / `orgId` / `accountId` | `Confirm` 响应（确认成功后回填） | 渠道侧留存 | 商户/企业/结算账户标识 |
| `merchantNo` | 进件成功后平台分配（渠道侧自备） | 余额查询、提现接口入参 | 商户号 |
| `busId` | `TradeImport` 返回 | `TradeConfirm` / `TradeStatus` / `TradeCancel` | 分账批次主记录 ID |
| `outWithdrawNo` | 调用方自行生成 | `Apply` / `TradeQuery` | 提现幂等键，超时重查不可换单号 |
| `withdrawNo` | `Apply` 响应 | 渠道侧留存 | 平台提现单号 |

## 业务流程

### 进件（商户入驻）

```
CreatePrepare ──> UploadFiles ──> [EDITING] ──> Confirm ──> SettlementStatus（轮询）
 (Step1 业务字段   (Step2 上传                      (确认提交,  (0审核中/1成功/2失败/
  + FileManifest)   资质文件)                        回填       3待激活/4激活中)
                                                      merchantId)
失败路径：SettlementStatus=2（FAILED）或草稿仍是 EDITING/FAILED 时
        └─> UpdatePrepare（业务字段整体重报，fileManifest 仅含新文件）
              ──> UploadFiles ──> Confirm ──> 轮询（循环直至成功）
```

顺序约束：`CreatePrepare` → `UploadFiles` 严格两步；`Confirm` 前须完成 `UploadFiles`；
仅 `EDITING` / `FAILED` 状态可 `UpdatePrepare`。

### 分账（文档交易订单导入）

```
GetUploadToken ──> TradeImport ──> TradeConfirm ──> TradeStatus（轮询）
 (Step1 fileName    (Step2 上传      (确认导入,       DocStatus:
  + fileSM3Hash)     XLSX → busId)    触发补单分账)   IMPORTING/PENDING/SUCCESS/FAILED)

取消路径：TradeCancel（仅取消尚未 Confirm 的批次）
```

### 提现

```
Apply（outWithdrawNo 幂等键 + totalAmount）──> TradeQuery（轮询）
                                              status: DEALING/WAIT_CONFIRM/SUCCESS/FAIL/UNKNOWN
危险规则：Apply 超时/中断时资金可能已出账，必须用原 outWithdrawNo 调 TradeQuery 查明结果，
        换单号重发等于再提现一笔；UNKNOWN 须联系平台，不能自行判为失败。
```

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

构造函数仅含 10 个必填参数，其余字段全部通过链式 `Set` 方法设置：

```go
dto := hst.NewCreatePrepareDto(
    "上海某某科技有限公司",  // legalName 法定名称
    "某某科技",             // shortName 商户简称
    []string{"WICOIN_PAY"}, // productCode 产品编码列表
    "03",                   // merchantBaseType 01自然人/02个体/03企业
    "brand_other",          // subRoleType 商户角色
    "01",                   // dealType 经营类型 01实体/02网络/03兼有
    "7299",                 // mcc 经营类目
    "13800000000",          // contactMobile 联系人手机号
    "张三",                  // contactName 联系人姓名
    "zhangsan@example.com", // email 邮箱
).SetPrincipal( // 负责人信息
    "13800000001",           // mobile
    "100",                   // certType 100身份证
    "310***********1234",    // certNo
    "李四",                   // person
    "2035-01-01 00:00:00",   // certVld
).SetLocation( // 地址信息
    "310000",                // province
    "310100",                // city
    "310104",                // district
    "上海市徐汇区XX路XX号",     // address
    "021-00000000",          // servicePhoneNo
).SetPerson("M", "企业法人"). // 自然人商户：性别、职业
    SetSettlementAccountType("01"). // 结算类型 01银行卡/02支付宝/03支付宝虚拟账户
    SetBank( // 结算账户
        "62220000000000000",     // bankCardNo
        "02",                    // accountType 01对私/02对公
        "某银行上海徐汇支行",       // branchName
        "310000", "310100",      // branchProvince / branchCity
    ).SetBankCertName("上海某某科技有限公司"). // 对公账户需设户名；对私无需
    SetContactLine("310100000000"). // 可选：联行号
    SetCardHolder( // 持卡人
        "01",                    // certType 01身份证
        "310***********1234",    // certNo
        "上海市徐汇区XX路XX号",    // cardHolderAddress
    ).SetFileManifest(&hst.FileManifest{ // 文件哈希清单
        CertPhotoAFiles:   []string{"<sm3_hash>"}, // 身份证人像面
        CertPhotoBFiles:   []string{"<sm3_hash>"}, // 身份证国徽面
        LicensePhotoFiles: []string{"<sm3_hash>"}, // 营业执照
    }).SetMerchantType("OTHER").  // 可选：商户业务类型
    SetTaxNum("31000000000000").  // 可选：税务登记证号码
    SetShareholder( // 可选：控股股东
        "王五", "100", "310***********1234", "2035-01-01 00:00:00",
    ).SetRemark("备注") // 可选：备注

result, _, err := client.CreatePrepare(ctx, dto)
// result.BizData.UploadToken   — 上传凭证（UUID）
// result.BizData.ExpireSeconds — 凭证有效期（秒，固定 900）
```

其余可选链式方法：`SetMangerLogonId`（管理员支付宝登录号）、`SetBussAuthVld`（执照有效期）、
`SetBussAuthType` / `SetBussAuthNo`（执照证件类型/号码）、`SetPersonCertVld`（自然人证件期限）、
`SetPartnerId`、`SetAgreementNo`、`SetAlipayPid`、`SetAlipayAccount`、`SetLogicGroupId`、
`SetWxSubMchId`、`SetWxSubMchAccount`、`SetLogonId`、`SetUserId`。

### 上传资质文件

两步上传 **Step 2**。以 `multipart/form-data` 携带凭证与文件上传，服务端按字段名 + 顺序重算 SM3 比对。

文件通过 `hst.NewUploadFile(name, data)` 构造（文件名 + 内容字节），再按 `FileManifest` 字段名经 `SetFiles` 挂载。典型场景是业务方在自己的接口中已收到用户上传的文件（`[]byte`），Step 1 用同一份字节计算 SM3 哈希提交 `FileManifest`，Step 2 直接透传该字节，全程无需落盘：

```go
dataA, _ := os.ReadFile("files/sfz-a.jpg") // 实际场景为接口收到的上传内容
dto := hst.NewUploadFilesDto("<upload_token>"). // uploadToken 来自 CreatePrepare/UpdatePrepare
    SetFiles("certPhotoAFiles", hst.NewUploadFile("sfz-a.jpg", dataA)). // 身份证人像面
    SetFiles("certPhotoBFiles", hst.NewUploadFile("sfz-b.jpg", dataB)). // 身份证国徽面
    SetFiles("licensePhotoFiles", hst.NewUploadFile("yyzz.jpg", dataC)) // 营业执照

bizData, err := client.UploadFiles(ctx, dto)
// bizData.DraftId     — 草稿 ID
// bizData.DraftStatus — 草稿状态 EDITING/SUBMITTING/CONFIRMED/FAILED
```

`SetFiles` 的字段名即 `FileManifest` 的 JSON 字段名（`certPhotoAFiles`、`licensePhotoFiles` 等，见上文 fileManifest 字段名清单）；非法字段名 `UploadFiles` 会直接返回错误，不会静默丢文件。

Hertz 接口中直接转发（不落盘）：

```go
func uploadHandler(ctx app.RequestContext, client *hst.Hst) {
    fh, _ := ctx.FormFile("certPhotoA")
    f, _ := fh.Open()
    defer f.Close()
    data, _ := io.ReadAll(f)
    // data 在 Step 1 已计算 SM3 哈希，此处直接透传
    dto := hst.NewUploadFilesDto(token).
        SetFiles("certPhotoAFiles", hst.NewUploadFile(fh.Filename, data))
    bizData, err := client.UploadFiles(ctx, dto)
    _ = bizData
}
```

> `Data` 字节须与 Step 1 计算 SM3 哈希的字节完全一致，否则服务端重算比对不一致将拒绝。

> `uploadToken` 为一次性凭证，无论成功与否都会被消费。上传字段须与 Step 1 的 `FileManifest` 一致。

### 更新草稿并申请上传凭证

仅 `EDITING` 或 `FAILED` 状态的草稿可更新。业务字段须整体重新提交，`fileManifest` 只需包含新文件字段。

```go
dto := hst.NewUpdatePrepareDto(
    "<draft_id>",       // draftId 草稿 ID
    "上海某某科技有限公司", // legalName
    "某某科技",          // shortName
    []string{"WICOIN_PAY"}, // productCode
    // ... 其余 7 个必填参数与 CreatePrepare 相同（merchantBaseType 至 email）
    // 链式 Set 方法与 CreatePrepare 完全一致
).SetFileManifest(&hst.FileManifest{
    ShopPhotoFiles: []string{"<sm3_hash>"},  // 仅更新门头照
}).SetMerchantType("OTHER")

result, _, err := client.UpdatePrepare(ctx, dto)
// result.BizData.UploadToken — 新的上传凭证
```

### 确认提交草稿

将草稿数据写入正式业务表，触发商户入驻申请。

```go
dto := hst.NewConfirmDto("<draft_id>")
result, _, err := client.Confirm(ctx, dto)
// result.BizData.DraftStatus — CONFIRMED 或 FAILED
// result.BizData.MerchantId  — 确认成功后回填
// result.BizData.OrgId       — 企业唯一号
// result.BizData.AccountId   — 结算账户 ID
```

### 查询入驻状态

```go
dto := hst.NewSettlementStatusDto("<draft_id>")
result, _, err := client.SettlementStatus(ctx, dto)
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
    "trade.xlsx", // fileName 文件名（须与 Step 2 上传的文件名一致）
    fileSM3Hash,  // fileSM3Hash 64位十六进制
)
result, _, err := client.GetUploadToken(ctx, dto)
// result.BizData.UploadToken   — 上传凭证
// result.BizData.ExpireSeconds — 有效期（秒，默认 300）
```

### 上传交易订单文件

两步上传 **Step 2**。`multipart/form-data` 上传 XLSX 文件。

```go
dto := hst.NewTradeImportDto(
    "<upload_token>",   // uploadToken 来自 GetUploadToken
    "files/trade.xlsx", // filePath XLSX 文件本地路径
)
busId, err := client.TradeImport(ctx, dto)
// busId — 业务主记录唯一 ID，用于后续确认/查询/取消
```

> `channelId` 由 SDK 自动填充，无需传入。multipart 字段名固定为 `file`，文件名取自路径基名。

> `uploadToken` 一次性消费，上传失败需从 Step 1 重新申请。

### 确认导入

确认已上传的文档交易批次并触发补单分账。

```go
dto := hst.NewTradeConfirmDto("<bus_id>")
result, _, err := client.TradeConfirm(ctx, dto)
// result.BizData — bool，true 表示确认成功
```

### 查询主记录状态

```go
dto := hst.NewTradeStatusDto("<bus_id>")
result, _, err := client.TradeStatus(ctx, dto)
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
result, _, err := client.TradeCancel(ctx, dto)
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
result, _, err := client.AvailableBalance(ctx, dto)
// result.BizData.BalanceInfos — 余额明细列表
// 按accountType取值，不要按下标取：
//   AVAILABLE_BALANCE — 可用余额（可提现）
//   PENDING_BALANCE  — 待结算金额（不可提现）
```

> `PENDING_BALANCE` 是尚未解冻的待结算金额，不可提现。把两者相加当作可提现额度会导致提现失败。

### 查询品牌商户专户余额

```go
dto := hst.NewBrandBalanceDto("<merchant_no>")
result, _, err := client.BrandBalance(ctx, dto)
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

result, _, err := client.Apply(ctx, dto)
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
result, _, err := client.TradeQuery(ctx, dto)
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
result, _, err := client.CreatePrepare(ctx, dto)
if err != nil {
    // err 包含 bizCode 和 bizMsg
    log.Printf("进件失败: %v", err)
    return
}
```

## API 参考

所有业务方法的接收者均为 `*Hst`，第一个参数均为 `ctx context.Context`。
加密信封类方法统一返回 `(result *SignObjectRespResult[T], err error)`；
multipart 类方法返回业务数据本身（无信封包装）。

### 方法签名（14 个）

```go
// 进件（加密信封，除 UploadFiles）
func (x *Hst) CreatePrepare(ctx context.Context, dto *CreatePrepareDto) (*SignObjectRespResult[*CreatePrepareBizData], *SignObjectResp, error)
func (x *Hst) UpdatePrepare(ctx context.Context, dto *UpdatePrepareDto) (*SignObjectRespResult[*UpdatePrepareBizData], *SignObjectResp, error)
func (x *Hst) UploadFiles(ctx context.Context, dto *UploadFilesDto) (*UploadFilesBizData, error) // multipart
func (x *Hst) Confirm(ctx context.Context, dto *ConfirmDto) (*SignObjectRespResult[*ConfirmBizData], *SignObjectResp, error)
func (x *Hst) SettlementStatus(ctx context.Context, dto *SettlementStatusDto) (*SignObjectRespResult[*SettlementStatusBizData], *SignObjectResp, error)

// 分账
func (x *Hst) GetUploadToken(ctx context.Context, dto *GetUploadTokenDto) (*SignObjectRespResult[*GetUploadTokenBizData], *SignObjectResp, error)
func (x *Hst) TradeImport(ctx context.Context, dto *TradeImportDto) (busId string, err error) // multipart
func (x *Hst) TradeConfirm(ctx context.Context, dto *TradeConfirmDto) (*SignObjectRespResult[bool], *SignObjectResp, error)
func (x *Hst) TradeStatus(ctx context.Context, dto *TradeStatusDto) (*SignObjectRespResult[*TradeStatusBizData], *SignObjectResp, error)
func (x *Hst) TradeCancel(ctx context.Context, dto *TradeCancelDto) (*SignObjectRespResult[bool], *SignObjectResp, error)

// 余额查询
func (x *Hst) AvailableBalance(ctx context.Context, dto *AvailableBalanceDto) (*SignObjectRespResult[*AvailableBalanceBizData], *SignObjectResp, error)
func (x *Hst) BrandBalance(ctx context.Context, dto *BrandBalanceDto) (*SignObjectRespResult[string], *SignObjectResp, error)

// 提现
func (x *Hst) Apply(ctx context.Context, dto *ApplyDto) (*SignObjectRespResult[*ApplyBizData], *SignObjectResp, error)
func (x *Hst) TradeQuery(ctx context.Context, dto *TradeQueryDto) (*SignObjectRespResult[*TradeQueryBizData], *SignObjectResp, error)

// 特殊用途
func NewHst(option *Option) (*Hst, error)
func (x *Hst) NewSignObjectReq(body Body) (*SignObjectReq, error) // 组装加密信封（特殊场景）
func (x *Hst) Request(ctx context.Context, path string, signObjectReq *SignObjectReq) (*SignObjectResp, error) // 发送信封请求，返回解密后的完整信封
```

### DTO 构造函数签名

```go
// 进件
func NewCreatePrepareDto(legalName, shortName string, productCode []string, merchantBaseType, subRoleType, dealType, mcc, contactMobile, contactName, email string) *CreatePrepareDto
func NewUpdatePrepareDto(draftId, legalName, shortName string, productCode []string, merchantBaseType, subRoleType, dealType, mcc, contactMobile, contactName, email string) *UpdatePrepareDto
func NewUploadFilesDto(uploadToken string) *UploadFilesDto
func NewConfirmDto(draftId string) *ConfirmDto
func NewSettlementStatusDto(draftId string) *SettlementStatusDto

// 分账
func NewGetUploadTokenDto(fileName, fileSM3Hash string) *GetUploadTokenDto
func NewTradeImportDto(uploadToken, filePath string) *TradeImportDto
func NewTradeConfirmDto(busId string) *TradeConfirmDto
func NewTradeStatusDto(busId string) *TradeStatusDto
func NewTradeCancelDto(busId string) *TradeCancelDto

// 余额查询
func NewAvailableBalanceDto(merchantNo string) *AvailableBalanceDto
func NewBrandBalanceDto(merchantNo string) *BrandBalanceDto

// 提现
func NewApplyDto(merchantNo, outWithdrawNo, totalAmount string) *ApplyDto
func NewTradeQueryDto(merchantNo, outWithdrawNo string) *TradeQueryDto

// 上传文件源（仅 UploadFiles 使用）
func NewUploadFile(name string, data []byte) *UploadFile // 文件名 + 内容字节（与 Step 1 SM3 哈希的字节一致）
func (x *UploadFilesDto) SetFiles(field string, files ...*UploadFile) *UploadFilesDto // 按 FileManifest 字段名挂载文件
```

### 响应 BizData 关键字段

| 方法 | BizData 类型 | 关键字段 |
|---|---|---|
| `CreatePrepare` / `UpdatePrepare` | `*CreatePrepareBizData` | `UploadToken`、`ExpireSeconds` |
| `UploadFiles` | `*UploadFilesBizData` | `DraftId`、`DraftStatus`（EDITING/SUBMITTING/CONFIRMED/FAILED）及草稿全量字段 |
| `Confirm` | `*ConfirmBizData` | `DraftStatus`、`MerchantId`、`OrgId`、`AccountId` |
| `SettlementStatus` | `*SettlementStatusBizData` | `SettlementStatus`（0-4）、`MerchantCreated`、`ActivateUrl` |
| `GetUploadToken` | `*GetUploadTokenBizData` | `UploadToken`、`ExpireSeconds` |
| `TradeImport` | `string`（busId） | — |
| `TradeConfirm` / `TradeCancel` | `bool` | — |
| `TradeStatus` | `*TradeStatusBizData` | `DocStatus`、`TotalDetailCount`、`SuccessCount`、`FailCount`、金额字段 |
| `AvailableBalance` | `*AvailableBalanceBizData` | `BalanceInfos`（按 `accountType` 取：AVAILABLE_BALANCE 可提现 / PENDING_BALANCE 待结算） |
| `BrandBalance` | `string` | 品牌专户余额（元，平台备付金，非商户额度） |
| `Apply` | `*ApplyBizData` | `WithdrawNo`、`Status`、`ErrorDesc`（幂等命中时另有金额/时间快照） |
| `TradeQuery` | `*TradeQueryBizData` | 订单完整快照（`Status`、`WithdrawFinishDate`、`ErrorDesc`） |

## 获取网关响应信封

加密信封类接口会同时返回业务结果和完整的 `SignObjectResp`：

```go
result, resp, err := client.CreatePrepare(ctx, dto)
if resp != nil {
    log.Printf("txn=%s timestamp=%s", resp.ClientReqTxn, resp.Timestamp)
}
```

要点：

- 网关响应能解析时，即使网关报错或验签失败，`resp` 也会返回
- 不需要外层响应的应用可用 `_` 忽略：`result, _, err := client.CreatePrepare(ctx, dto)`
- `UploadFiles` / `TradeImport` 为 multipart 普通响应，不产生加密信封

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
