package model

import (
	"fmt"
	"time"
)

// SignObjectReq 请求对象模型，用于签名和加密处理
type SignObjectReq struct {
	// 时间戳
	Timestamp string `json:"timestamp"`
	// 请求流水号
	ReqTxn string `json:"reqTxn"`
	// 初始化向量(Hex格式)
	IvHex string `json:"ivHex"`
	// 客户端身份标识
	ChannelId string `json:"channelId"`
	// 版本号固定
	Version string `json:"version"`
	// 签名字段
	Signature string `json:"signature"`
	// 实际业务数据对象的JSON字符串,加签时会进行加密处理
	Body string `json:"body"`
}

// NewSignObjectReq 创建一个新的 SignObjectReq 并设置默认值
func NewSignObjectReq() *SignObjectReq {
	return &SignObjectReq{
		Timestamp: fmt.Sprintf("%d", time.Now().UnixMilli()),
		Version:   "1.0",
	}
}
