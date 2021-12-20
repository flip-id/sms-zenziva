package sms_zenziva_local

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// Config config app for zenziva
type Config struct {
	BaseUrl        string
	UserKey        string
	PasswordKey    string
	ConnectTimeout int
}

type Sender struct {
	Config Config
}
// end of config

// Message Response from zenziva
type Message struct {
	XMLName   xml.Name `xml:"message"`
	MessageId string   `xml:"messageId"`
	To        string   `xml:"to"`
	Status    int      `xml:"status"`
	Text      string   `xml:"text"`
	Balance   int      `xml:"balance"`
}

type ResponseBody struct {
	XMLName xml.Name `xml:"response"`
	Message Message  `xml:"message"`
}

// end of response

// ReqMessage request for zenziva
type ReqMessage struct {
	PhoneNumber string
	Text        string
}

// end of request

func New(config Config) *Sender {
	return &Sender{
		Config: config,
	}
}

// CallbackData callback data
type CallbackData struct {
	Error struct {
		Description string `json:"description"`
		GroupID     int    `json:"group_id"`
		GroupName   string `json:"group_name"`
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Permanent   bool   `json:"permanent"`
	} `json:"error"`
	MessageID string `json:"message_id"`
	To        string `json:"to"`
	Status    int    `json:"status"`
	Text      string `json:"text"`
	Balance   int    `json:"balance"`
	SentAt    string `json:"sent_at"`
}
// end of callback data

// SendSMS function to send message
func (s *Sender) SendSMS(request ReqMessage) (ResponseBody, error) {
	path := fmt.Sprintf("%s", s.Config.BaseUrl)

	req, err := http.NewRequest(http.MethodGet, path, nil)
	param := req.URL.Query()
	param.Add("userkey", s.Config.UserKey)
	param.Add("passkey", s.Config.PasswordKey)
	param.Add("nohp", request.PhoneNumber)
	param.Add("pesan", request.Text)
	req.URL.RawQuery = param.Encode()
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		log.Error(err)
		return ResponseBody{}, err
	}

	timeout := s.Config.ConnectTimeout
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	res, err := client.Do(req)

	if err != nil {
		log.Error(err)
		return ResponseBody{}, err
	}

	if res.StatusCode >= 400 {
		log.Error(res)
		return ResponseBody{}, errors.New("failed to send SMS")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(res)
			return
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	resBody := ResponseBody{}

	if err != nil {
		log.Error(res)
		return ResponseBody{}, err
	}

	err = xml.Unmarshal(body, &resBody)

	if err != nil {
		log.Error(res)
		return ResponseBody{}, err
	}

	return resBody, nil
}
