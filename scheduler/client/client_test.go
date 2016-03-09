package client

import (
	"fmt"
	"testing"
)

func Test_client_1(t *testing.T) {
	opts := ClientOpts{}
	client, err := NewClient(opts)
	if err != nil {

		t.Error("test not pass")
	}
	info, err := client.database.GetUserInfo("nobody")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(info)
}
