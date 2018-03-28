package broker

import "fmt"

func i2s(any interface{}) string {
	return fmt.Sprintf("%v", any)
}