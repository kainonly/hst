package model

// BizResp 业务响应对象
type BizResp struct {
	BizData    interface{} `json:"bizData"`
	BizMsg     string      `json:"bizMsg"`
	BizCode    string      `json:"bizCode"`
	BizSuccess bool        `json:"bizSuccess"`
}
