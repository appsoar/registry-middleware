package debug

import (
	"fmt"
)

var (
	Debug = false
)

func Print(dat ...interface{}) {
	if Debug {
		for _, v := range dat {
			fmt.Println(v)
		}
	}

}

func init() {
	Debug = true
}
