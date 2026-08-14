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

## 文档

详细接口文档见 [docs/API.md](docs/API.md)：

| 模块 | 文档 | 接口数 |
|---|---|---|
| 进件 | [docs/merchant.md](docs/merchant.md) | 5 |
| 分账 | [docs/trade.md](docs/trade.md) | 5 |
| 余额查询 | [docs/account.md](docs/account.md) | 2 |
| 提现 | [docs/withdrawal.md](docs/withdrawal.md) | 2 |

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
