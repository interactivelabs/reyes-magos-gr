package admin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getOrderId_err(t *testing.T) {
	mockOrder := NewMockOrder(t)
	mockOrder.EXPECT().Param("order_id").Return("not number")
	_, err := getOrderId(mockOrder)
	assert.Error(t, err)
}
