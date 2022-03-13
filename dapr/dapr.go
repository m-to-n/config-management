package dapr

import (
	daprc "github.com/dapr/go-sdk/client"
	daprcommon "github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"
	"log"
	"sync"
)

var (
	once       sync.Once
	instance_c daprc.Client
	instance_d daprcommon.Service
)

// dapr sidecar http port
const DAPR_HTTP_PORT = "3500"

// dapr sidecar grp port
const DAPR_GRPC_PORT = "35000"

// dapr app grpc address
const DAPR_APP_GRPC_ADDR = ":35001"

func DaprClient() daprc.Client {
	once.Do(func() {
		client, err := daprc.NewClientWithPort(DAPR_GRPC_PORT)

		if err != nil {
			// fmt.Printf("DaprClient error %s", err.Error())
			// return
			log.Fatalf("failed to start dapr client: %v", err)
		}
		instance_c = client
	})
	return instance_c
}

func DaprService() daprcommon.Service {
	once.Do(func() {

		s, err := daprd.NewService(DAPR_APP_GRPC_ADDR)
		if err != nil {
			log.Fatalf("failed to start dapr server: %v", err)
		}

		instance_d = s
	})
	return instance_d
}
