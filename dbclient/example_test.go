package dbclient_test

import (
	"fmt"

	"github.com/jarrod-lowe/jmap-service-libs/dbclient"
)

func ExampleAccountPK() {
	pk := dbclient.AccountPK("abc123")
	fmt.Println(pk)
	// Output: ACCOUNT#abc123
}

func ExampleUserPK() {
	pk := dbclient.UserPK("user456")
	fmt.Println(pk)
	// Output: USER#user456
}
