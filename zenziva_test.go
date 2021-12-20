package sms_zenziva_local

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// integration test
func TestSendSmsSuccess(t *testing.T) {
	//client := resty.New().SetTimeout(time.Second * time.Duration(15))
	//httpmock.ActivateNonDefault(client.GetClient())
	//defer httpmock.DeactivateAndReset()
	//
	//fixture := `<?xml version="1.0" encoding="UTF-8"?>
	//			<response>
	//				<message>
	//					<messageId>59167697</messageId>
	//					<to>+6289662233555</to>
	//					<status>0</status>
	//					<text>Success</text>
	//					<balance>4695</balance>
	//				</message>
	//			</response>`
	//
	//responder, _ := httpmock.NewXmlResponder(200, fixture)
	//fakeUrl :=  "https://testzenziva.com"
	//httpmock.RegisterResponder("GET", fakeUrl, responder)

	config := Config{
		BaseUrl:        "https://testzenziva.com",
		UserKey:        "test",
		PasswordKey:    "test",
		ConnectTimeout: 15,
	}
	sender := New(config)
	reqMsg := ReqMessage{
		PhoneNumber: "089662233555",
		Text:        "tes",
	}

	res, err := sender.SendSMS(reqMsg)

	assert.Nil(t, err)
	assert.NotNil(t, res)
}
