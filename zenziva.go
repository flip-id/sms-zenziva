package smszenziva

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/hystrix"
	log "github.com/sirupsen/logrus"
)

const (
	// DefaultTimeout sets the default timeout of the HTTP client.
	DefaultTimeout = 30 * time.Second
)

// Config is a config app for Zenziva.
type Config struct {
	BaseURL        string
	UserKey        string
	PasswordKey    string
	ConnectTimeout time.Duration
	Client         heimdall.Doer
}

// Validate validates the config variables to ensure smooth integration.
func (c *Config) Validate() (err error) {
	if c.UserKey == "" {
		err = ErrEmptyUserKey
		return
	}

	if c.PasswordKey == "" {
		err = ErrEmptyPasswordKey
		return
	}

	return
}

// DefaultV1 sets the config default value for version 1.
func (c *Config) DefaultV1() *Config {
	if c.BaseURL == "" {
		c.BaseURL = "https://reguler.zenziva.net/apps/smsapi.php"
	}

	return c.defaultVal()
}

func (c *Config) defaultVal() *Config {
	if c.ConnectTimeout < DefaultTimeout {
		c.ConnectTimeout = DefaultTimeout
	}

	if c.Client == nil {
		c.Client = http.DefaultClient
	}

	return c
}

// Sender is a client of Zenziva.
type Sender struct {
	client *hystrix.Client
	config *Config
}

// end of config

// Message is a response from Zenziva.
type Message struct {
	XMLName   xml.Name `xml:"message"`
	MessageID string   `xml:"messageId"`
	To        string   `xml:"to"`
	Status    int      `xml:"status"`
	Text      string   `xml:"text"`
	Balance   int      `xml:"balance"`
}

// ResponseBody is an XML response body from the Zenziva.
type ResponseBody struct {
	XMLName xml.Name `xml:"response"`
	Message Message  `xml:"message"`
}

// end of response

// ReqMessage is a request for Zenziva.
type ReqMessage struct {
	PhoneNumber string
	Text        string
}

// end of request

// NewV1 initializes a new Sender for the version 1 of Zenziva API.
func NewV1(config *Config) (client *Sender, err error) {
	if config == nil {
		err = ErrNilArgs
		return
	}

	err = config.Validate()
	if err != nil {
		return
	}

	config = config.DefaultV1()
	client = &Sender{
		client: hystrix.NewClient(
			hystrix.WithHTTPTimeout(config.ConnectTimeout),
			hystrix.WithHystrixTimeout(config.ConnectTimeout),
			hystrix.WithHTTPClient(config.Client),
		),
		config: config,
	}
	return
}

// CallbackData is a callback data.
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

// SendSMSV1 function to send message using V1 Zenziva API.
func (s *Sender) SendSMSV1(request ReqMessage) (respBody ResponseBody, err error) {
	req, err := http.NewRequest(http.MethodGet, s.config.BaseURL, nil)
	if err != nil {
		log.Error(err)
		return
	}

	param := req.URL.Query()
	param.Add("userkey", s.config.UserKey)
	param.Add("passkey", s.config.PasswordKey)
	param.Add("nohp", request.PhoneNumber)
	param.Add("pesan", request.Text)
	req.URL.RawQuery = param.Encode()
	req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationForm)
	res, err := s.client.Do(req)
	if err != nil {
		log.Error(res, err)
		return
	}

	defer func(Body io.ReadCloser) {
		if Body == nil {
			return
		}

		err := Body.Close()
		if err != nil {
			log.Error(res, err)
		}
	}(res.Body)

	if res.StatusCode >= http.StatusBadRequest {
		log.Error(res)
		err = ErrFailedToSendSMS
		return
	}

	err = xml.NewDecoder(res.Body).Decode(&respBody)
	if err != nil {
		log.Error(res, err)
	}

	return
}
