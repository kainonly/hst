package model

import "strings"

// SignObjectResp 响应对象模型，用于签名和加密处理
type SignObjectResp struct {
	// 响应码 成功: SUCCESS
	Code string `json:"code"`
	// 响应消息
	Msg string `json:"msg"`
	// 时间戳
	Timestamp string `json:"timestamp"`
	// 客户端请求流水号
	ClientReqTxn string `json:"clientReqTxn"`
	// 服务端应用ID
	ServerAppId string `json:"serverAppId"`
	// 初始化向量(Hex格式)
	IvHex string `json:"ivHex"`
	// 版本号固定
	Version string `json:"version"`
	// 签名字段
	Signature string `json:"signature"`
	// 实际业务数据对象的JSON字符串,加签时会进行加密处理
	Body string `json:"body"`
}

// IsSuccess 判断响应是否成功
func (r *SignObjectResp) IsSuccess() bool {
	return strings.EqualFold(r.Code, "SUCCESS")
}
