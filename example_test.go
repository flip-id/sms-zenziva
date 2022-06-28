package smszenziva

import (
	"fmt"
	"log"
	"net/http"
)

func ExampleNewV1() {
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
