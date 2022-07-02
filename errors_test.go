package zenziva

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestError_Error(t *testing.T) {
	r := ResponseXMLMessage{
		MessageID: "300",
		To:        "081234567890",
		Status:    100,
		Text:      "ERROR",
		Balance:   decimal.NewFromInt(100),
	}
	err := (new(Error)).Assign(r)
	text := fmt.Sprintf(
		"failed to send SMS to: %s, message ID: %s, status: %d, text: %s, balance: %s",
		r.To,
		r.MessageID,
		r.Status,
		r.Text,
		r.Balance,
	)
	assert.Equal(t, text, err.Error())
}
