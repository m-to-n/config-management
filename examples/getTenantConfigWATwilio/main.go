package main

import (
	"context"
	"encoding/json"
	"fmt"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/m-to-n/common/logging"
	"github.com/m-to-n/common/tenants"
	"log"
	"os"
)

func main() {
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()
	ctx := context.Background()

	content := &dapr.DataContent{
		ContentType: "application/json",
		Data:        []byte(fmt.Sprintf(`{ "accountSid": "%s", "receiverPhoneNumber": "%s" }`, os.Args[1], os.Args[2])),
	}

	result, err := client.InvokeMethodWithContent(ctx, "config-management", "getTenantConfigForTwilioWAReq", "get", content)
	log.Println("calling getTenantConfigForTwilioWAReq...")
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
