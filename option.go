package zenziva

import (
	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/hystrix"
	"net/http"
	"time"
)

const (
	// DefaultTimeout sets the default timeout of the HTTP client.
	DefaultTimeout = 30 * time.Second
	// BaseURLV1 sets the base URL of the API version 1.
	BaseURLV1 = "https://reguler.zenziva.net/apps/smsapi.php"
)

// FnOption is a function that sets the option.
type FnOption func(o *Option)

// WithBaseURL sets the base URL of the API.
func WithBaseURL(s string) FnOption {
	return func(o *Option) {
		o.BaseURL = s
	}
}

// WithUserKey sets the user key.
func WithUserKey(s string) FnOption {
	return func(o *Option) {
		o.UserKey = s
	}
}

// WithPasswordKey sets the password key.
func WithPasswordKey(s string) FnOption {
	return func(o *Option) {
		o.PasswordKey = s
	}
}

// WithTimeout sets the timeout of the HTTP client.
func WithTimeout(t time.Duration) FnOption {
	return func(o *Option) {
		o.ConnectTimeout = t
	}
}

// WithClient sets the HTTP client.
func WithClient(c heimdall.Doer) FnOption {
	return func(o *Option) {
		o.Client = c
	}
}

// WithHystrixOptions sets the hystrix options.
func WithHystrixOptions(opts ...hystrix.Option) FnOption {
	return func(o *Option) {
		o.HystrixOptions = append(o.HystrixOptions, opts...)
	}
}

// Option is a config for Zenziva.
type Option struct {
	BaseURL        string
	UserKey        string
	PasswordKey    string
	ConnectTimeout time.Duration
	Client         heimdall.Doer
	HystrixOptions []hystrix.Option
	client         *hystrix.Client
}

func (o *Option) Assign(opts ...FnOption) *Option {
	for _, opt := range opts {
		opt(o)
	}

	return o
}

// Validate validates the config variables to ensure smooth integration.
func (o *Option) Validate() (err error) {
	if o.UserKey == "" {
		err = ErrEmptyUserKey
		return
	}

	if o.PasswordKey == "" {
		err = ErrEmptyPasswordKey
	}
	return
}

// DefaultV1 sets the config default value for version 1.
func (o *Option) DefaultV1() *Option {
	if o.BaseURL == "" {
		o.BaseURL = BaseURLV1
	}
	return o.defaultVal()
}

func (o *Option) defaultVal() *Option {
	if o.ConnectTimeout < DefaultTimeout {
		o.ConnectTimeout = DefaultTimeout
	}

	if o.Client == nil {
		o.Client = http.DefaultClient
	}

	opts := append([]hystrix.Option{
		hystrix.WithHTTPTimeout(o.ConnectTimeout),
		hystrix.WithHystrixTimeout(o.ConnectTimeout),
		hystrix.WithHTTPClient(o.Client),
	}, o.HystrixOptions...)
	o.client = hystrix.NewClient(opts...)
	return o
}
