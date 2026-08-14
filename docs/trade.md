# 分账接口

文档交易订单文件导入相关接口，采用两步上传流程（先申请凭证，再上传文件）。

## 接口列表

| 方法 | 路径 | 说明 |
|---|---|---|
| `GetUploadToken` | `/channel/doc-trade-file/getUploadToken` | 申请文件上传凭证（Step 1） |
| `TradeImport` | `/channel-file/doc-trade-file/import` | 上传交易订单文件（Step 2，multipart） |
| `TradeConfirm` | `/channel/doc-trade-file/confirm` | 确认导入 |
| `TradeStatus` | `/channel/doc-trade-file/status` | 查询主记录状态 |
| `TradeCancel` | `/channel/doc-trade-file/cancel` | 取消导入 |

---

## 申请文件上传凭证

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

---

## 上传交易订单文件

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

---

## 确认导入

确认已上传的文档交易批次并触发补单分账。

```go
dto := hst.NewTradeConfirmDto("<bus_id>")
result, err := client.TradeConfirm(ctx, dto)
// result.BizData — bool，true 表示确认成功
```

---

## 查询主记录状态

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

---

## 取消导入

取消尚未确认的文档交易导入批次。

```go
dto := hst.NewTradeCancelDto("<bus_id>")
result, err := client.TradeCancel(ctx, dto)
// result.BizData — bool，true 表示取消成功
```
