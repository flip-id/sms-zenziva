package zenziva

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
)

// List of errors in the Zenziva package.
var (
	ErrFailedToSendSMS  = errors.New("failed to send SMS")
	ErrEmptyUserKey     = errors.New("empty user key")
	ErrEmptyPasswordKey = errors.New("empty password key")
)

// Error is a wrapper for the error returned by the Zenziva package.
type Error struct {
	MessageID string          `json:"message_id"`
	To        string          `json:"to"`
	Status    int             `json:"status"`
	Text      string          `json:"text"`
	Balance   decimal.Decimal `json:"balance"`
}

// Error returns the error message.
func (e *Error) Error() (res string) {
	res = fmt.Sprintf(
		"failed to send SMS to: %s, message ID: %s, status: %d, text: %s, balance: %s",
		e.To,
		e.MessageID,
		e.Status,
		e.Text,
		e.Balance,
	)
	return
}

// Assign assigns the response to the Error.
func (e *Error) Assign(resp ResponseXMLMessage) error {
	if resp.Status == 0 {
		return nil
	}

	e.MessageID = resp.MessageID
	e.To = resp.To
	e.Status = resp.Status
	e.Text = resp.Text
	e.Balance = resp.Balance
	return e
}
