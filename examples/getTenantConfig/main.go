package main

import (
	"context"
	dapr "github.com/dapr/go-sdk/client"
	"log"
)

/**
https://github.com/dapr/go-sdk/tree/main/examples
*/

func main() {
	tenantId := "tenant-123"
	client, err := dapr.NewClientWithPort("35001")
	if err != nil {
		panic(err)
	}
	defer client.Close()
	ctx := context.Background()
	//Using Dapr SDK to invoke a method
	result, err := client.InvokeMethod(ctx, "config-management", "getTenantConfig", "get")
	log.Println("Tenant config requested: " + tenantId)
	log.Println("Result: ")
	log.Println(result)
}
