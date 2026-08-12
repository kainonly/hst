package hst_test

import (
	"context"
	"testing"

	"github.com/kainonly/hst"
	"github.com/stretchr/testify/assert"
)

func TestSettlementStatus(t *testing.T) {
	ctx := context.Background()

	dto := hst.NewSettlementStatusDto(`123`)
	result, err := client.SettlementStatus(ctx, dto)
	assert.NoError(t, err)

	t.Log(result)
}
