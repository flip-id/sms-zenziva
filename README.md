# zenziva
<a title="Doc for package zenziva." target="_blank" href="https://pkg.go.dev/github.com/flip-id/sms-zenziva?tab=doc"><img src="https://img.shields.io/badge/go.dev-doc-007d9c?style=flat-square&logo=read-the-docs"></a>
[![Go Report Card](https://goreportcard.com/badge/github.com/flip-id/sms-zenziva)](https://goreportcard.com/report/github.com/flip-id/sms-zenziva)

Package zenziva is a library to use the Zenziva service.
This library uses [Hystrix client](https://github.com/gojek/heimdall/v7) as its underlying HTTP client.

## Documentation

To show the documentation of the package, we can check the code directly or by running this command:
```bash
make doc
```

This will open the package documentation in local.
We can access it in `http://localhost:6060/pkg/github.com/flip-id/sms-zenziva`.

# Example

This library can be used based on the example shown in the URL below:

`http://localhost:6060/pkg/github.com/flip-id/sms-zenziva/#example_NewV1`

Script:
```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	c, err := NewV1(&Config{
		UserKey:     "test-user",
		PasswordKey: "test-password",
		Client:      http.DefaultClient,
	})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := c.SendSMSV1(ReqMessage{
		PhoneNumber: "+6281001002003",
		Text:        "Hello Zenziva!",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Resp: %+v\n", resp)
}
```