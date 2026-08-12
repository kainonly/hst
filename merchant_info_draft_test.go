package hst_test

import (
	"context"
	"testing"

	"github.com/kainonly/hst"
	"github.com/stretchr/testify/assert"
)

func TestMerchantInfoDraft(t *testing.T) {
	ctx := context.Background()

	body := hst.NewSettlementStatusBody(`123`)
	bizData, err := client.SettlementStatus(ctx, body)
	assert.NoError(t, err)

	t.Log(bizData)
}
