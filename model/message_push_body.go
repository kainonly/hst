package model

// MessagePushBody 消息推送体
type MessagePushBody struct {
	// 投递类型 1-单推 2-标签推，3-批量推
	SubmitType string `json:"submitType"`
	// 设备SN编号
	Sn string `json:"sn"`
	// 商户编号
	Mid string `json:"mid"`
	// 播报数字
	BroadcastNumber string `json:"broadcastNumber"`
	// 语种(中文：CN；英文：EN；马来:MS）
	Lang string `json:"lang"`
	// 支付方式
	PayType string `json:"payType"`
	// 金额
	Amount string `json:"amount"`
}

// NewMessagePushBody 创建一个新的 MessagePushBody 并设置默认值
func NewMessagePushBody() *MessagePushBody {
	return &MessagePushBody{
		Lang:    "CN",
		PayType: "TNG",
	}
}
