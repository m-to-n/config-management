package main

import (
	"context"
	"encoding/json"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/m-to-n/common/logging"
	"github.com/m-to-n/common/tenants"
	"log"
)

/**
https://github.com/dapr/go-sdk/tree/main/examples
*/

func main() {
	tenantId := "tenant-123"
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()
	ctx := context.Background()

	content := &dapr.DataContent{
		ContentType: "application/json",
		Data:        []byte(`{ "tenantId": "tenant-123" }`),
	}

	result, err := client.InvokeMethodWithContent(ctx, "config-management", "getTenantConfig", "get", content)
	log.Println("Tenant config requested: " + tenantId)
	log.Println("Result: ")
	log.Println(result)

	var tenant tenants.TenantConfig
	err = json.Unmarshal(result, &tenant)
	if err != nil {
		log.Fatalf("tenant unamrshaling error: %s ", err.Error())
	}
	structStrPtr, _ := logging.StructToPrettyString(tenant)
	log.Println(*structStrPtr)
}
