# 余额查询接口

查询商户资金账户余额，用于结算对账与提现前的额度判断。

## 接口列表

| 方法 | 路径 | 说明 |
|---|---|---|
| `AvailableBalance` | `/channel/merchant_account/available_balance` | 查询商户可用余额 |
| `BrandBalance` | `/channel/merchant_account/brand-balance` | 查询品牌商户专户余额 |

---

## 查询商户可用余额

```go
dto := hst.NewAvailableBalanceDto("<merchant_no>")
result, err := client.AvailableBalance(ctx, dto)
// result.BizData.BalanceInfos — 余额明细列表
// 按accountType取值，不要按下标取：
//   AVAILABLE_BALANCE — 可用余额（可提现）
//   PENDING_BALANCE  — 待结算金额（不可提现）
```

> `PENDING_BALANCE` 是尚未解冻的待结算金额，不可提现。把两者相加当作可提现额度会导致提现失败。

---

## 查询品牌商户专户余额

```go
dto := hst.NewBrandBalanceDto("<merchant_no>")
result, err := client.BrandBalance(ctx, dto)
// result.BizData — 字符串，品牌专户余额（单位元）
```

> 此余额是平台侧备付金，同一平台配置下不同商户号查到的是同一个余额。不能作为商户可用/可提现额度。
