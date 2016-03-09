package client

import (
	"fmt"
	"testing"
	"time"
)

func Test_client_1(t *testing.T) {
	client, err := NewClient()
	if err != nil {

		t.Error("test not pass")
	}
	/*	info, err := client.database.GetUserInfo("nobody")
		if err != nil {
			t.Error(err)
		}*/
	for i := 0; i < 1000; i++ {
		go func() {
			info, err := client.registry.ListImages()
			if err == nil {

				fmt.Println(string(info.([]byte)))
			} else {
				fmt.Println(err)
			}
		}()

		go func() {
			info1, err := client.registry.GetImageTags("hh/redis")
			if err == nil {

				fmt.Println(string(info1.([]byte)))
			} else {
				fmt.Println(err)
			}
		}()

	}
	time.Sleep(50 * time.Second)

}
