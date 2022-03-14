package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dapr/go-sdk/service/common"
	"github.com/m-to-n/config-management/dapr"
	"github.com/m-to-n/config-management/tenant"
	"log"
)

func testingStuff() {
	/*if err := tenant.CreateDummyTenant(); err != nil {
		panic(err)
	}  */

	ctx := context.Background()
	tcp, err := tenant.GetTenantConfig(ctx, "tenant-123")
	if err != nil {
		panic(err)
	}

	tenantConfig := *tcp

	fmt.Println("I am done here %s", tenantConfig)

	fmt.Println("I am done here!!!") /**/
}

type getTenantConfigReqData struct {
	TenantId string `json:"tenantId"`
}

// curl http://localhost:3500/v1.0/invoke/config-management/method/getTenantConfig/tenant-123
func getTenantConfigHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, errO error) {
	log.Printf("getTenantConfigHandler - ContentType:%s, Verb:%s, QueryString:%s, %+v", in.ContentType, in.Verb, in.QueryString, string(in.Data))
	// do something with the invocation here

	// using dummyTenantData requires import: common_tenants "github.com/m-to-n/common/tenants"
	// dummyTenantData, _ := common_tenants.TenantConfigToJson(utils.GetDummyTenant())

	var reqData getTenantConfigReqData
	err := json.Unmarshal(in.Data, &reqData)
	if err != nil {
		log.Printf("error when parsing request: %s", err.Error())
		return nil, err
	}

	log.Println("Calling GetTenantConfig for: " + reqData.TenantId)
	dbTenantData, _ := tenant.GetTenantConfig(ctx, reqData.TenantId)
	dbTenantDataBytes, _ := json.Marshal(*dbTenantData)

	out = &common.Content{
		// Data:        in.Data,
		// Data: dummyTenantData,
		Data:        dbTenantDataBytes,
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}

	return
}

func main() {
	// testingStuff()
	s := dapr.DaprService()

	if err := s.AddServiceInvocationHandler("getTenantConfig", getTenantConfigHandler); err != nil {
		log.Fatalf("error adding invocation handler: %v", err)
	}

	if err := s.Start(); err != nil {
		log.Fatalf("dapr server error: %v", err)
	} /**/
}
