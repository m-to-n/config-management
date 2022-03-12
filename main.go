package main

import (
	"context"
	"fmt"
	"github.com/dapr/go-sdk/service/common"
	common_tenants "github.com/m-to-n/common/tenants"
	"github.com/m-to-n/config-management/dapr"
	"github.com/m-to-n/config-management/tenant"
	"github.com/m-to-n/config-management/utils"
	"log"
)

func testingStuff() {
	/* if err := tenant.CreateDummyTenant(); err != nil {
		panic(err)
	} */

	tcp, err := tenant.GetTenantConfig("tenant-123")
	if err != nil {
		panic(err)
	}

	tenantConfig := *tcp

	fmt.Println("I am done here %s", tenantConfig)

	fmt.Println("I am done here!!!")
}

// curl http://localhost:3500/v1.0/invoke/config-management/method/getTenantConfig/tenant-123
func getTenantConfigHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	log.Printf("getTenantConfigHandler - ContentType:%s, Verb:%s, QueryString:%s, %+v", in.ContentType, in.Verb, in.QueryString, string(in.Data))
	// do something with the invocation here

	dummyTenantData, _ := common_tenants.TenantConfigToJson(utils.GetDummyTenant())

	out = &common.Content{
		// Data:        in.Data,
		Data:        dummyTenantData,
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}

	return
}

func main() {
	s := dapr.DaprService()

	if err := s.AddServiceInvocationHandler("getTenantConfig", getTenantConfigHandler); err != nil {
		log.Fatalf("error adding invocation handler: %v", err)
	}

	if err := s.Start(); err != nil {
		log.Fatalf("dapr server error: %v", err)
	}
}
