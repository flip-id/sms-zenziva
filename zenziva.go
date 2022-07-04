package zenziva

import (
	"context"
	"encoding/xml"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"net/http"
)

// Sender is a client of Zenziva.
type Sender interface {
	// SendSMSV1 function to send message using V1 Zenziva API.
	// This function is based on this documentation: https://reguler.zenziva.net/apps/download/Zenziva-SMSReguler-HttpApi.pdf.
	// This function will return error if the response from Zenziva is not successful.
	// This function will return respBody with the error message even if the response from Zenziva is not successful.
	SendSMSV1(ctx context.Context, request RequestSendSMSV1) (respBody ResponseXML, err error)
}

type sender struct {
	opt *Option
}

// Assign assigns the given options to the sender.
func (s *sender) Assign(opt *Option) *sender {
	if opt == nil {
		return s
	}

	newOpt := *opt
	s.opt = &newOpt
	return s
}

// NewV1 initializes a new Sender for the version 1 of Zenziva API.
func NewV1(opts ...FnOption) (client Sender, err error) {
	opt := (new(Option)).Assign(opts...).DefaultV1()

	err = opt.Validate()
	if err != nil {
		return
	}

	client = new(sender).Assign(opt)
	return
}

// ResponseXML is an XML template response from Zenziva.
type ResponseXML struct {
	XMLName xml.Name           `xml:"response"`
	Message ResponseXMLMessage `xml:"message"`
}

// GetError returns an error if the response from Zenziva is not successful.
func (r *ResponseXML) GetError() (err error) {
	err = (new(Error)).Assign(r.Message)
	return
}

// ResponseXMLMessage is a response message from Zenziva.
type ResponseXMLMessage struct {
	XMLName   xml.Name        `xml:"message"`
	MessageID string          `xml:"messageId"`
	To        string          `xml:"to"`
	Status    int             `xml:"status"`
	Text      string          `xml:"text"`
	Balance   decimal.Decimal `xml:"balance"`
}

// RequestSendSMSV1 is a request for send SMS version 1 to Zenziva.
type RequestSendSMSV1 struct {
	PhoneNumber string
	Text        string
}

// SendSMSV1 function to send message using V1 Zenziva API.
// This function is based on this documentation: https://reguler.zenziva.net/apps/download/Zenziva-SMSReguler-HttpApi.pdf.
func (s *sender) SendSMSV1(ctx context.Context, request RequestSendSMSV1) (respBody ResponseXML, err error) {
	req, err := http.NewRequest(http.MethodGet, s.opt.BaseURL, nil)
	if err != nil {
		return
	}

	param := req.URL.Query()
	param.Add("userkey", s.opt.UserKey)
	param.Add("passkey", s.opt.PasswordKey)
	param.Add("nohp", request.PhoneNumber)
	param.Add("pesan", request.Text)
	req.URL.RawQuery = param.Encode()
	req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationForm)
	req = req.WithContext(ctx)
	res, err := s.opt.client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		if res.Body == nil {
			return
		}

		_ = res.Body.Close()
	}()

	err = s.formatUnknown(res)
	if err != nil {
		return
	}

	err = xml.NewDecoder(res.Body).Decode(&respBody)
	if err != nil {
		return
	}

	err = respBody.GetError()
	return
}

func (s *sender) isError(resp *http.Response) bool {
	return resp.StatusCode >= http.StatusBadRequest
}

func (s *sender) formatUnknown(resp *http.Response) (err error) {
	if !s.isError(resp) {
		return
	}

	err = formatUnknown(resp)
	return
}
