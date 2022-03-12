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

const DAPR_HTTP_PORT = "3500"
const DAPR_GRPC_PORT = "35000"

func DaprClient() daprc.Client {
	once.Do(func() {
		client, err := daprc.NewClientWithPort(DAPR_GRPC_PORT)

		if err != nil {
			fmt.Printf("DaprClient error %s", err.Error())
			return
		}
		instance = client
	})
	return instance
}
