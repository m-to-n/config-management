package main

import (
	"context"
	dapr "github.com/dapr/go-sdk/client"
	"io/ioutil"
	"log"
	"os"
)

/**
https://github.com/dapr/go-sdk/tree/main/examples
*/

func main() {
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()
	ctx := context.Background()

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal("File reading error", err)
	}

	/* var tenantConfig common_tenants.TenantConfig
	err = json.Unmarshal(data, &tenantConfig)
	if err != nil {
		log.Fatal("error when parsing file data(must be TenantConfig): %s", err.Error())
	} */

	content := &dapr.DataContent{
		ContentType: "application/json",
		Data:        data,
	}

	result, err := client.InvokeMethodWithContent(ctx, "config-management", "createTenantConfig", "post", content)
	log.Println("Tenant config creation requested.")
	log.Println("Result: ")
	log.Println(result)
}
