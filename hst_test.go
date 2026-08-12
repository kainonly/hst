package hst_test

import (
	"os"
	"testing"

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
	os.Exit(m.Run())
}
