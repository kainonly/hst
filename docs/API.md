# 接口文档

## 文档导航

- [初始化](#初始化)
- [进件接口](merchant.md)
- [分账接口](trade.md)
- [余额查询接口](account.md)
- [提现接口](withdrawal.md)
- [错误处理](#错误处理)
- [数据类型](#数据类型)

---

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

---

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

---

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
