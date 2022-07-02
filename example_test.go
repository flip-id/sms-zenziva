package zenziva

import (
	"fmt"
	"log"
	"net/http"
)

func ExampleNewV1() {
	c, err := NewV1(
		WithUserKey("userkey"),
		WithPasswordKey("passwordkey"),
		WithClient(http.DefaultClient),
	)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := c.SendSMSV1(RequestSendSMSV1{
		PhoneNumber: "+6281001002003",
		Text:        "Hello Zenziva!",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Resp: %+v\n", resp)
}
