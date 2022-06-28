package zenziva

import (
	"encoding/xml"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestNewV1(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name          string
		args          func() args
		wantNilClient bool
		wantErr       bool
	}{
		{
			name: "nil config",
			args: func() args {
				return args{nil}
			},
			wantNilClient: true,
			wantErr:       true,
		},
		{
			name: "empty user key",
			args: func() args {
				return args{config: &Config{}}
			},
			wantNilClient: true,
			wantErr:       true,
		},
		{
			name: "empty password key",
			args: func() args {
				return args{config: &Config{
					UserKey: "test-user",
				}}
			},
			wantNilClient: true,
			wantErr:       true,
		},
		{
			name: "success getting the client",
			args: func() args {
				c := &Config{
					UserKey:     "test",
					PasswordKey: "test",
				}
				return args{c}
			},
			wantNilClient: false,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.args()
			gotClient, err := NewV1(args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewV1() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantNilClient {
				assert.Nil(t, gotClient)
			}
		})
	}
}

func getSender(c *Config) *Sender {
	s, _ := NewV1(c)
	return s
}

func TestSender_SendSMSV1(t *testing.T) {
	type args struct {
		request ReqMessage
	}
	tests := []struct {
		name         string
		s            func() *Sender
		args         func() args
		wantRespBody func() ResponseBody
		wantErr      bool
		responder    func()
	}{
		{
			name: "invalid base URL",
			s: func() *Sender {
				return getSender(&Config{
					BaseURL:     "://test",
					UserKey:     "test-user",
					PasswordKey: "test-password",
				})
			},
			args: func() args {
				return args{
					request: ReqMessage{},
				}
			},
			wantRespBody: func() ResponseBody {
				return ResponseBody{}
			},
			wantErr: true,
		},
		{
			name: "server is unreachable",
			s: func() *Sender {
				return getSender(&Config{
					BaseURL:     "flip://test.local",
					UserKey:     "test-user",
					PasswordKey: "test-password",
				})
			},
			args: func() args {
				return args{
					request: ReqMessage{},
				}
			},
			wantRespBody: func() ResponseBody {
				return ResponseBody{}
			},
			wantErr: true,
		},
		{
			name: "bad request from the server",
			s: func() *Sender {
				return getSender(&Config{
					BaseURL:     "flip://test.local",
					UserKey:     "test-user",
					PasswordKey: "test-password",
					Client:      http.DefaultClient,
				})
			},
			args: func() args {
				return args{
					request: ReqMessage{},
				}
			},
			wantRespBody: func() ResponseBody {
				return ResponseBody{}
			},
			wantErr: true,
			responder: func() {
				httpmock.Activate()
				httpmock.RegisterResponder(
					http.MethodGet,
					"flip://test.local",
					httpmock.NewStringResponder(http.StatusBadRequest, "error request"),
				)
			},
		},
		{
			name: "xml decode error",
			s: func() *Sender {
				return getSender(&Config{
					BaseURL:     "flip://test.local",
					UserKey:     "test-user",
					PasswordKey: "test-password",
					Client:      http.DefaultClient,
				})
			},
			args: func() args {
				return args{
					request: ReqMessage{},
				}
			},
			wantRespBody: func() ResponseBody {
				return ResponseBody{}
			},
			wantErr: true,
			responder: func() {
				httpmock.Activate()
				httpmock.RegisterResponder(
					http.MethodGet,
					"flip://test.local",
					httpmock.NewXmlResponderOrPanic(http.StatusOK, ""),
				)
			},
		},
		{
			name: "no responder",
			s: func() *Sender {
				return getSender(&Config{
					BaseURL:     "flip://test.local",
					UserKey:     "test-user",
					PasswordKey: "test-password",
					Client:      http.DefaultClient,
				})
			},
			args: func() args {
				return args{
					request: ReqMessage{},
				}
			},
			wantRespBody: func() ResponseBody {
				return ResponseBody{}
			},
			wantErr: true,
			responder: func() {
				httpmock.Activate()
				httpmock.RegisterNoResponder(
					httpmock.NewXmlResponderOrPanic(http.StatusOK, ""),
				)
			},
		},
		{
			name: "nil body",
			s: func() *Sender {
				return getSender(&Config{
					BaseURL:     "flip://test.local",
					UserKey:     "test-user",
					PasswordKey: "test-password",
					Client:      http.DefaultClient,
				})
			},
			args: func() args {
				return args{
					request: ReqMessage{},
				}
			},
			wantRespBody: func() ResponseBody {
				return ResponseBody{}
			},
			wantErr: true,
			responder: func() {
				httpmock.Activate()
				httpmock.RegisterResponder(http.MethodGet, "flip://test.local", func(r *http.Request) (*http.Response, error) {
					return &http.Response{
						Close:      false,
						Body:       nil,
						StatusCode: http.StatusOK,
						Status:     http.StatusText(http.StatusOK),
					}, nil
				})
			},
		},
		{
			name: "no user key",
			s: func() *Sender {
				s := getSender(&Config{
					BaseURL:     "flip://test.local",
					UserKey:     "test-user",
					PasswordKey: "test-password",
					Client:      http.DefaultClient,
				})
				s.config.UserKey = ""
				return s
			},
			args: func() args {
				return args{
					request: ReqMessage{},
				}
			},
			wantRespBody: func() ResponseBody {
				return ResponseBody{
					XMLName: xml.Name{Local: "response"},
					Message: Message{
						XMLName: xml.Name{Local: "message"},
						Status:  5,
						Text:    "Userkey atau Passkey Salah",
					},
				}
			},
			wantErr: false,
			responder: func() {
				httpmock.Activate()
				httpmock.RegisterResponder(http.MethodGet, "flip://test.local", httpmock.NewXmlResponderOrPanic(http.StatusOK, httpmock.File("example/error.xml")))
			},
		},
		{
			name: "success send sms",
			s: func() *Sender {
				s := getSender(&Config{
					BaseURL:     "flip://test.local",
					UserKey:     "test-user",
					PasswordKey: "test-password",
					Client:      http.DefaultClient,
				})
				s.config.UserKey = ""
				return s
			},
			args: func() args {
				return args{
					request: ReqMessage{},
				}
			},
			wantRespBody: func() ResponseBody {
				return ResponseBody{
					XMLName: xml.Name{Local: "response"},
					Message: Message{
						XMLName:   xml.Name{Local: "message"},
						MessageID: "59167697",
						To:        "+6289662233555",
						Status:    0,
						Text:      "Success",
						Balance:   4695,
					},
				}
			},
			wantErr: false,
			responder: func() {
				httpmock.Activate()
				httpmock.RegisterResponder(http.MethodGet, "flip://test.local", httpmock.NewXmlResponderOrPanic(http.StatusOK, httpmock.File("example/success.xml")))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer httpmock.DeactivateAndReset()
			if tt.responder != nil {
				tt.responder()
			}

			args := tt.args()
			s := tt.s()
			gotRespBody, err := s.SendSMSV1(args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sender.SendSMSV1() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.EqualValues(t, tt.wantRespBody(), gotRespBody)
		})
	}
}
