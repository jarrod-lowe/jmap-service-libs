package jmaperror_test

import (
	"errors"
	"fmt"

	"github.com/jarrod-lowe/jmap-service-libs/jmaperror"
)

func ExampleInvalidArguments() {
	err := jmaperror.InvalidArguments("mailboxId must be provided")
	fmt.Println(err.Type())
	fmt.Println(err.Error())
	// Output:
	// invalidArguments
	// invalidArguments: mailboxId must be provided
}

func ExampleServerFail() {
	underlying := errors.New("connection refused")
	err := jmaperror.ServerFail("database unavailable", underlying)
	fmt.Println(err.Type())
	fmt.Println(errors.Is(err, underlying))
	// Output:
	// serverFail
	// true
}

func ExampleInvalidProperties() {
	err := jmaperror.InvalidProperties("invalid values", []string{"name", "email"})
	m := err.ToMap()
	fmt.Println(m["type"])
	fmt.Println(m["properties"])
	// Output:
	// invalidProperties
	// [name email]
}

func ExampleNotJSON() {
	err := jmaperror.NotJSON("request body is not valid JSON")
	m := err.ToMap()
	fmt.Println(m["type"])
	fmt.Println(m["status"])
	// Output:
	// urn:ietf:params:jmap:error:notJSON
	// 400
}

func ExampleLimit() {
	err := jmaperror.Limit("maxSizeRequest", "request exceeds 10MB")
	m := err.ToMap()
	fmt.Println(m["type"])
	fmt.Println(m["limit"])
	// Output:
	// urn:ietf:params:jmap:error:limit
	// maxSizeRequest
}
