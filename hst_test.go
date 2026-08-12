package hst_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bytedance/sonic"
	"github.com/kainonly/hst"
	"gopkg.in/yaml.v3"
)

var client *hst.Hst
var v *Values
var cfg *ConfigMap

type Values struct {
	Context string                `yaml:"context"`
	Configs map[string]*ConfigMap `yaml:"configs"`
}

func (v *Values) Config() *ConfigMap {
	return v.Configs[v.Context]
}

type ConfigMap struct {
	BaseURL    string `yaml:"base_url"`
	PriKey     string `yaml:"pri_key"`      // 客户端私钥
	WicoPubKey string `yaml:"wico_pub_key"` // 支付中心公钥
	EncryptKey string `yaml:"encrypt_key"`  // 加密密钥
	MerchantNo string `yaml:"merchant_no"`  // 商户号
	ChannelId  string `yaml:"channel_id"`   // 客户端身份标识
	// 自然人商户测试挡板信息
	PersonMobile        string `yaml:"person_mobile"`         // 手机号
	PersonCertNo        string `yaml:"person_cert_no"`        // 证件号码
	PersonUsername      string `yaml:"person_username"`       // 用户名
	PersonUid           string `yaml:"person_uid"`            // 支付宝 uid
	PersonMybankAccount string `yaml:"person_mybank_account"` // 网商二类户
	PersonBankCardNo    string `yaml:"person_bank_card_no"`   // 外部银行卡号
}

func TestMain(m *testing.M) {
	var err error
	var b []byte
	if b, err = os.ReadFile(`./config/values.yml`); err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(b, &v); err != nil {
		panic(err)
	}
	cfg = v.Config()
	if client, err = hst.NewHst(&hst.Option{
		BaseURL:    cfg.BaseURL,
		PriKey:     cfg.PriKey,
		WicoPubKey: cfg.WicoPubKey,
		EncryptKey: cfg.EncryptKey,
		MerchantNo: cfg.MerchantNo,
		ChannelId:  cfg.ChannelId,
	}); err != nil {
		panic(err)
	}
	// 确保 logs 目录存在
	if err = os.MkdirAll("logs", 0755); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

// logResult 将业务结果写入 logs/<name>.log（追加模式）并在测试日志中展示。
// result 为任意业务响应（含 BizSuccess/BizCode/BizMsg/BizData）。
func logResult(t *testing.T, name string, result any) {
	t.Helper()

	// 序列化完整 result 为 JSON
	jsonBytes, err := sonic.Marshal(result)
	if err != nil {
		t.Logf("序列化结果失败: %v", err)
		return
	}
	jsonStr := string(jsonBytes)

	// 追加写入 logs/<name>.log（带时间戳）
	logPath := filepath.Join("logs", name+".log")
	line := fmt.Sprintf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), jsonStr)
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		t.Logf("打开日志文件失败: %v", err)
		return
	}
	defer f.Close()
	if _, err := f.WriteString(line); err != nil {
		t.Logf("写入日志文件失败: %v", err)
	}

	// 测试日志展示完整 JSON
	t.Logf("=== %s ===", name)
	t.Logf("result: %s", jsonStr)
}

// readLastLogBizData 从 logs/<name>.log 读取最后一行，解析 bizData 字段到传入的指针。
// 用于测试间数据传递（如 upload_files 测试从 create_prepare 日志读取 uploadToken）。
// 支持两种日志格式：
//  1. 有外层包装 {"bizSuccess":true,"bizData":{...}} → 取 bizData
//  2. 直接是 bizData {...} → 整行直接解析
func readLastLogBizData(t *testing.T, name string, bizData any) {
	t.Helper()
	logPath := filepath.Join("logs", name+".log")
	b, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("读取日志文件 %s 失败: %v", logPath, err)
	}
	lines := splitNonEmptyLines(string(b))
	if len(lines) == 0 {
		t.Fatalf("日志文件 %s 无有效记录", logPath)
	}
	lastLine := lines[len(lines)-1]
	// 去掉行首的时间戳 [yyyy-MM-dd HH:mm:ss]
	jsonPart := lastLine
	if idx := indexOfByte(lastLine, ']'); idx >= 0 && idx+1 < len(lastLine) {
		jsonPart = lastLine[idx+1:]
	}
	// 去掉可能的行首空格
	for len(jsonPart) > 0 && (jsonPart[0] == ' ' || jsonPart[0] == '\t') {
		jsonPart = jsonPart[1:]
	}
	// 尝试解析为外层包装（含 bizSuccess/bizData）
	var wrapper struct {
		BizSuccess bool            `json:"bizSuccess"`
		BizCode    string          `json:"bizCode"`
		BizMsg     string          `json:"bizMsg"`
		BizData    json.RawMessage `json:"bizData"`
	}
	if err := sonic.UnmarshalString(jsonPart, &wrapper); err == nil && len(wrapper.BizData) > 0 {
		// 有外层包装
		if !wrapper.BizSuccess {
			t.Fatalf("日志记录的业务失败: bizCode=%s bizMsg=%s", wrapper.BizCode, wrapper.BizMsg)
		}
		if err := sonic.Unmarshal(wrapper.BizData, bizData); err != nil {
			t.Fatalf("解析 bizData 失败: %v", err)
		}
		return
	}
	// 无外层包装，整行直接解析为 bizData
	if err := sonic.UnmarshalString(jsonPart, bizData); err != nil {
		t.Fatalf("解析日志 JSON 失败: %v (line=%s)", err, lastLine)
	}
}

// splitNonEmptyLines 按换行符分割，过滤空行。
func splitNonEmptyLines(s string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			line := s[start:i]
			// 去掉 \r
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			if line != "" {
				result = append(result, line)
			}
			start = i + 1
		}
	}
	// 处理最后一行（不以 \n 结尾的情况）
	if start < len(s) {
		line := s[start:]
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
		if line != "" {
			result = append(result, line)
		}
	}
	return result
}

// indexOfByte 返回字节 b 在 s 中首次出现的位置，不存在返回 -1。
func indexOfByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}
