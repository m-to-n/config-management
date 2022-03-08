package dapr

import (
	"fmt"
	daprc "github.com/dapr/go-sdk/client"
	"sync"
)

var (
	once     sync.Once
	instance daprc.Client
)

func DaprClient() daprc.Client {
	once.Do(func() {
		client, err := daprc.NewClient()

		if err != nil {
			fmt.Printf("DaprClient error %s", err.Error())
			return
		}
		instance = client
	})
	return instance
}
