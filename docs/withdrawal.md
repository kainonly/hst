# 提现接口

渠道商代商户发起提现，把可用余额提现到商户绑定的结算银行卡。

## 接口列表

| 方法 | 路径 | 说明 |
|---|---|---|
| `Apply` | `/channel/merchant/withdrawal/apply` | 商户提现申请 |
| `TradeQuery` | `/channel/merchant/withdrawal/query` | 查询提现订单 |

---

## 商户提现申请

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

---

## 查询提现订单

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
